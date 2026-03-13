package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"terraform-provider-stonebranch/internal/client"
)

// Generator handles exporting StoneBranch resources to Terraform HCL.
type Generator struct {
	client       *client.Client
	output       string
	noDeps       bool
	exported     map[string]bool   // Track exported resources to avoid duplicates
	nameCounters map[string]int    // Counters for generating sequential resource names
	hclBuffer    *bytes.Buffer     // Single buffer for all output
	nameMap      map[string]string // Maps resource key to generated terraform name
}

// NewGenerator creates a new Generator.
func NewGenerator(client *client.Client, output string, noDeps bool) *Generator {
	return &Generator{
		client:       client,
		output:       output,
		noDeps:       noDeps,
		exported:     make(map[string]bool),
		nameCounters: make(map[string]int),
		hclBuffer:    &bytes.Buffer{},
		nameMap:      make(map[string]string),
	}
}

// generateResourceName creates a sequential identifier like "task_unix_001"
func (g *Generator) generateResourceName(resourceType string) string {
	g.nameCounters[resourceType]++
	return fmt.Sprintf("%s_%03d", resourceType, g.nameCounters[resourceType])
}

// getOrCreateResourceName returns the terraform resource name for a given resource
func (g *Generator) getOrCreateResourceName(resourceType, name string) string {
	key := fmt.Sprintf("%s/%s", resourceType, name)
	if tfName, ok := g.nameMap[key]; ok {
		return tfName
	}
	tfName := g.generateResourceName(resourceType)
	g.nameMap[key] = tfName
	return tfName
}

// ExportResource exports a single resource by type and name, appending to the buffer.
func (g *Generator) ExportResource(ctx context.Context, resourceType, name string) error {
	if g.isExported(resourceType, name) {
		return nil
	}

	rt := GetResourceType(resourceType)
	if rt == nil {
		return fmt.Errorf("unknown resource type: %s", resourceType)
	}

	// Fetch the resource
	data, err := g.fetchResource(ctx, rt, name)
	if err != nil {
		return err
	}

	// Generate sequential resource name
	tfName := g.getOrCreateResourceName(resourceType, name)

	// Add template fields
	data["_resourceName"] = tfName
	data["_terraformResource"] = rt.TerraformResource
	data["_originalName"] = name

	// Generate HCL
	hcl, err := g.generateHCLFromData(rt, data)
	if err != nil {
		return err
	}

	// Mark as exported
	g.markExported(resourceType, name)

	// Append to buffer with comment showing original name
	g.hclBuffer.WriteString(fmt.Sprintf("# %s: %s\n", resourceType, name))
	g.hclBuffer.WriteString(hcl)
	g.hclBuffer.WriteString("\n")

	return nil
}

// ExportAll exports all resources of a given type matching the filter.
func (g *Generator) ExportAll(ctx context.Context, resourceType, filter string) error {
	rt := GetResourceType(resourceType)
	if rt == nil {
		return fmt.Errorf("unknown resource type: %s", resourceType)
	}

	items, err := g.listResources(ctx, rt, filter)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Fprintf(os.Stderr, "No %s resources found\n", resourceType)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d %s resources\n", len(items), resourceType)

	for _, item := range items {
		if item.Name == "" {
			continue
		}

		// For workflows, use the complete export
		if resourceType == "task_workflow" && !g.noDeps {
			if err := g.exportWorkflowComplete(ctx, item.Name); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to export workflow %s: %v\n", item.Name, err)
			}
		} else {
			if err := g.ExportResource(ctx, resourceType, item.Name); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to export %s/%s: %v\n", resourceType, item.Name, err)
			}
		}
	}

	return nil
}

// ExportTasks exports tasks matching a filter pattern.
// Workflows will include their contained tasks, vertices, and edges.
func (g *Generator) ExportTasks(ctx context.Context, filter string) error {
	// List all tasks matching the filter (API supports * and ? wildcards)
	query := url.Values{}
	if filter != "" {
		query.Set("taskname", filter)
	}

	respBody, err := g.client.Get(ctx, "/resources/task/listadv", query)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	var rawItems []map[string]interface{}
	if err := json.Unmarshal(respBody, &rawItems); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(rawItems) == 0 {
		fmt.Fprintf(os.Stderr, "No tasks found matching filter: %s\n", filter)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d tasks matching filter\n", len(rawItems))

	// Separate workflows from other tasks
	var workflows []map[string]interface{}
	var otherTasks []map[string]interface{}

	for _, item := range rawItems {
		taskType, _ := item["type"].(string)
		if taskType == "taskWorkflow" {
			workflows = append(workflows, item)
		} else {
			otherTasks = append(otherTasks, item)
		}
	}

	// First, export workflows (with their tasks, vertices, edges)
	for _, wf := range workflows {
		name, _ := wf["name"].(string)
		if name == "" {
			continue
		}
		if err := g.exportWorkflowComplete(ctx, name); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to export workflow %s: %v\n", name, err)
		}
	}

	// Then export other tasks (that weren't already exported as part of a workflow)
	for _, task := range otherTasks {
		name, _ := task["name"].(string)
		taskType, _ := task["type"].(string)
		if name == "" || taskType == "" {
			continue
		}

		cliType, ok := APITypeToResourceType[taskType]
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: unsupported task type %s for task %s\n", taskType, name)
			continue
		}

		if g.isExported(cliType, name) {
			continue
		}

		if err := g.ExportResource(ctx, cliType, name); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to export %s/%s: %v\n", cliType, name, err)
		}
	}

	return nil
}

// exportWorkflowComplete exports a workflow with its tasks, vertices, and edges.
// Order: workflow definition -> tasks -> vertices -> edges
func (g *Generator) exportWorkflowComplete(ctx context.Context, name string) error {
	if g.isExported("task_workflow", name) {
		return nil
	}

	rt := GetResourceType("task_workflow")
	if rt == nil {
		return fmt.Errorf("unknown resource type: task_workflow")
	}

	// Fetch the workflow task
	data, err := g.fetchResource(ctx, rt, name)
	if err != nil {
		return fmt.Errorf("failed to fetch workflow: %w", err)
	}

	// Generate workflow resource name
	wfTfName := g.getOrCreateResourceName("task_workflow", name)
	data["_resourceName"] = wfTfName
	data["_terraformResource"] = rt.TerraformResource
	data["_originalName"] = name

	// Generate workflow HCL
	workflowHCL, err := g.generateHCLFromData(rt, data)
	if err != nil {
		return fmt.Errorf("failed to generate workflow HCL: %w", err)
	}

	// Mark workflow as exported
	g.markExported("task_workflow", name)

	// Write workflow section header
	g.hclBuffer.WriteString(fmt.Sprintf("\n# ============================================================\n"))
	g.hclBuffer.WriteString(fmt.Sprintf("# Workflow: %s\n", name))
	g.hclBuffer.WriteString(fmt.Sprintf("# ============================================================\n\n"))
	g.hclBuffer.WriteString(fmt.Sprintf("# task_workflow: %s\n", name))
	g.hclBuffer.WriteString(workflowHCL)
	g.hclBuffer.WriteString("\n")

	// Fetch workflow vertices
	vertices, err := g.fetchWorkflowVertices(ctx, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to fetch workflow vertices for %s: %v\n", name, err)
		return nil
	}

	if len(vertices) == 0 {
		return nil
	}

	// Export all tasks in the workflow
	g.hclBuffer.WriteString("# --- Tasks in Workflow ---\n")
	for _, v := range vertices {
		taskName := v.Task.Value
		if taskName == "" {
			continue
		}

		// Fetch the task
		taskData, err := g.fetchResourceByName(ctx, "/resources/task", "taskname", taskName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to fetch task %s: %v\n", taskName, err)
			continue
		}

		taskType, ok := taskData["type"].(string)
		if !ok {
			continue
		}

		cliType, ok := APITypeToResourceType[taskType]
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: unsupported task type %s for task %s\n", taskType, taskName)
			continue
		}

		if g.isExported(cliType, taskName) {
			continue
		}

		taskRT := GetResourceType(cliType)
		if taskRT == nil {
			continue
		}

		// Generate task resource name
		taskTfName := g.getOrCreateResourceName(cliType, taskName)
		taskData["_resourceName"] = taskTfName
		taskData["_terraformResource"] = taskRT.TerraformResource
		taskData["_originalName"] = taskName

		taskHCL, err := g.generateHCLFromData(taskRT, taskData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to generate HCL for task %s: %v\n", taskName, err)
			continue
		}

		g.hclBuffer.WriteString(fmt.Sprintf("# %s: %s\n", cliType, taskName))
		g.hclBuffer.WriteString(taskHCL)
		g.hclBuffer.WriteString("\n")

		g.markExported(cliType, taskName)
	}

	// Generate workflow vertices and track their terraform names
	g.hclBuffer.WriteString("# --- Workflow Vertices ---\n")
	vertexTfNames := make(map[string]string) // Maps vertexId to terraform resource name
	for _, v := range vertices {
		taskName := v.Task.Value
		if taskName == "" {
			continue
		}
		vertexTfName, vertexHCL, err := g.generateWorkflowVertexHCLNew(name, wfTfName, v)
		if err == nil {
			vertexTfNames[v.VertexId] = vertexTfName
			g.hclBuffer.WriteString(vertexHCL)
			g.hclBuffer.WriteString("\n")
		}
	}

	// Generate workflow edges
	edges, err := g.fetchWorkflowEdges(ctx, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to fetch workflow edges for %s: %v\n", name, err)
	} else if len(edges) > 0 {
		g.hclBuffer.WriteString("# --- Workflow Edges ---\n")
		for _, e := range edges {
			edgeHCL, err := g.generateWorkflowEdgeHCLNew(name, wfTfName, e, vertices, vertexTfNames)
			if err == nil {
				g.hclBuffer.WriteString(edgeHCL)
				g.hclBuffer.WriteString("\n")
			}
		}
	}

	return nil
}

// ExportWorkflow exports a workflow and all its components including dependent tasks.
func (g *Generator) ExportWorkflow(ctx context.Context, name string) error {
	// Export the workflow task itself
	if err := g.ExportResource(ctx, "task_workflow", name); err != nil {
		return fmt.Errorf("failed to export workflow: %w", err)
	}

	if g.noDeps {
		return nil
	}

	// Fetch workflow vertices
	vertices, err := g.fetchWorkflowVertices(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to fetch workflow vertices: %w", err)
	}

	// Export each task in the workflow
	for _, v := range vertices {
		taskName := v.Task.Value
		if taskName == "" || g.isExported("task", taskName) {
			continue
		}

		// Determine task type and export
		taskData, err := g.fetchResourceByName(ctx, "/resources/task", "taskname", taskName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to fetch task %s: %v\n", taskName, err)
			continue
		}

		taskType, ok := taskData["type"].(string)
		if !ok {
			continue
		}

		cliType, ok := APITypeToResourceType[taskType]
		if !ok {
			fmt.Fprintf(os.Stderr, "Warning: unsupported task type %s for task %s\n", taskType, taskName)
			continue
		}

		if err := g.ExportResource(ctx, cliType, taskName); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to export task %s: %v\n", taskName, err)
		}

		// Generate workflow vertex
		vertexHCL, err := g.generateWorkflowVertexHCL(name, v)
		if err == nil {
			g.outputHCL("workflow_vertex", fmt.Sprintf("%s_%s", name, taskName), vertexHCL)
		}
	}

	// Fetch and generate workflow edges
	edges, err := g.fetchWorkflowEdges(ctx, name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to fetch workflow edges: %v\n", err)
	} else {
		for i, e := range edges {
			edgeHCL, err := g.generateWorkflowEdgeHCL(name, e, vertices)
			if err == nil {
				g.outputHCL("workflow_edge", fmt.Sprintf("%s_edge_%d", name, i), edgeHCL)
			}
		}
	}

	return nil
}

// Finalize writes the buffered output to file or stdout.
func (g *Generator) Finalize() error {
	if g.hclBuffer.Len() == 0 {
		return nil
	}

	// Add header
	header := "# Generated by sb2tf from StoneBranch Universal Controller\n\n"
	content := header + g.hclBuffer.String()

	if g.output == "" {
		// Write to stdout
		fmt.Print(content)
		return nil
	}

	// Write to file
	if err := os.MkdirAll(g.output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := filepath.Join(g.output, "main.tf")
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	fmt.Fprintf(os.Stderr, "Wrote %s\n", filename)

	return nil
}

// fetchResource fetches a single resource from the API.
func (g *Generator) fetchResource(ctx context.Context, rt *ResourceType, name string) (map[string]interface{}, error) {
	return g.fetchResourceByName(ctx, rt.APIEndpoint, rt.NameQueryParam, name)
}

func (g *Generator) fetchResourceByName(ctx context.Context, endpoint, paramName, name string) (map[string]interface{}, error) {
	query := url.Values{}
	query.Set(paramName, name)

	respBody, err := g.client.Get(ctx, endpoint, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch resource: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return data, nil
}

// ListResources lists all resources of a given type (exported for CLI use).
func (g *Generator) ListResources(ctx context.Context, rt *ResourceType, filter string) ([]ResourceItem, error) {
	return g.listResources(ctx, rt, filter)
}

// listResources lists all resources of a given type.
func (g *Generator) listResources(ctx context.Context, rt *ResourceType, filter string) ([]ResourceItem, error) {
	query := url.Values{}
	// Don't filter by type in API - filter locally instead
	if filter != "" {
		query.Set(rt.NameQueryParam, filter)
	}

	respBody, err := g.client.Get(ctx, rt.ListEndpoint, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	// Parse as raw JSON to handle different field names
	var rawItems []map[string]interface{}
	if err := json.Unmarshal(respBody, &rawItems); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to ResourceItem, handling custom name fields
	nameField := rt.NameField
	if nameField == "" {
		nameField = "name"
	}

	var items []ResourceItem
	for _, raw := range rawItems {
		item := ResourceItem{}
		if name, ok := raw[nameField].(string); ok {
			item.Name = name
		}
		if t, ok := raw["type"].(string); ok {
			item.Type = t
		}
		items = append(items, item)
	}

	// Filter by type locally if this resource type has a specific API type value
	if rt.APITypeValue != "" {
		var filtered []ResourceItem
		for _, item := range items {
			if item.Type == rt.APITypeValue {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	return items, nil
}

// ResourceItem represents a resource in list responses.
type ResourceItem struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
	Task struct {
		Value string `json:"value"`
	} `json:"task,omitempty"`
}

// exportDependencies exports resources that this resource depends on.
func (g *Generator) exportDependencies(ctx context.Context, rt *ResourceType, data map[string]interface{}) error {
	for _, dep := range rt.Dependencies {
		// Check condition if present
		if dep.Condition != "" {
			parts := strings.Split(dep.Condition, "=")
			if len(parts) == 2 {
				if val, ok := data[parts[0]].(string); !ok || val != parts[1] {
					continue
				}
			}
		}

		// Get the referenced resource name
		refName, ok := data[dep.Field].(string)
		if !ok || refName == "" {
			continue
		}

		// Skip if already exported
		if g.isExported(dep.ResourceType, refName) {
			continue
		}

		// Export the dependency
		if err := g.ExportResource(ctx, dep.ResourceType, refName); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to export dependency %s/%s: %v\n", dep.ResourceType, refName, err)
		}
	}

	return nil
}

// fetchWorkflowVertices fetches all vertices for a workflow.
func (g *Generator) fetchWorkflowVertices(ctx context.Context, workflowName string) ([]WorkflowVertex, error) {
	query := url.Values{}
	query.Set("workflowname", workflowName)

	respBody, err := g.client.Get(ctx, "/resources/workflow/vertices", query)
	if err != nil {
		return nil, err
	}

	var vertices []WorkflowVertex
	if err := json.Unmarshal(respBody, &vertices); err != nil {
		return nil, err
	}

	return vertices, nil
}

// fetchWorkflowEdges fetches all edges for a workflow.
func (g *Generator) fetchWorkflowEdges(ctx context.Context, workflowName string) ([]WorkflowEdge, error) {
	query := url.Values{}
	query.Set("workflowname", workflowName)

	respBody, err := g.client.Get(ctx, "/resources/workflow/edges", query)
	if err != nil {
		return nil, err
	}

	var edges []WorkflowEdge
	if err := json.Unmarshal(respBody, &edges); err != nil {
		return nil, err
	}

	return edges, nil
}

// WorkflowVertex represents a vertex in the workflow.
type WorkflowVertex struct {
	Task struct {
		Value string `json:"value"`
	} `json:"task"`
	VertexId string `json:"vertexId"`
	Alias    string `json:"alias,omitempty"`
	VertexX  string `json:"vertexX,omitempty"`
	VertexY  string `json:"vertexY,omitempty"`
}

// WorkflowEdge represents an edge in the workflow.
type WorkflowEdge struct {
	SourceId struct {
		Value string `json:"value"`
	} `json:"sourceId"`
	TargetId struct {
		Value string `json:"value"`
	} `json:"targetId"`
	StraightEdge bool `json:"straightEdge,omitempty"`
}

// generateHCL generates HCL for a resource (legacy - uses sanitized names).
func (g *Generator) generateHCL(rt *ResourceType, data map[string]interface{}) (string, error) {
	tmpl := GetTemplate(rt.CLIName)
	if tmpl == nil {
		return "", fmt.Errorf("no template for resource type: %s", rt.CLIName)
	}

	// Add computed fields to data
	name, _ := data["name"].(string)
	data["_resourceName"] = SanitizeName(name)
	data["_terraformResource"] = rt.TerraformResource

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// generateHCLFromData generates HCL for a resource using pre-set _resourceName.
func (g *Generator) generateHCLFromData(rt *ResourceType, data map[string]interface{}) (string, error) {
	tmpl := GetTemplate(rt.CLIName)
	if tmpl == nil {
		return "", fmt.Errorf("no template for resource type: %s", rt.CLIName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// generateWorkflowVertexHCL generates HCL for a workflow vertex.
func (g *Generator) generateWorkflowVertexHCL(workflowName string, v WorkflowVertex) (string, error) {
	data := map[string]interface{}{
		"workflowName":       workflowName,
		"taskName":           v.Task.Value,
		"alias":              v.Alias,
		"_resourceName":      SanitizeName(fmt.Sprintf("%s_%s", workflowName, v.Task.Value)),
		"_terraformResource": "stonebranch_workflow_vertex",
	}

	tmpl := GetTemplate("workflow_vertex")
	if tmpl == nil {
		return "", fmt.Errorf("no template for workflow_vertex")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateWorkflowEdgeHCL generates HCL for a workflow edge.
func (g *Generator) generateWorkflowEdgeHCL(workflowName string, e WorkflowEdge, vertices []WorkflowVertex) (string, error) {
	// Find task names for source and target vertex IDs
	var sourceTask, targetTask string
	for _, v := range vertices {
		if v.VertexId == e.SourceId.Value {
			sourceTask = v.Task.Value
		}
		if v.VertexId == e.TargetId.Value {
			targetTask = v.Task.Value
		}
	}

	data := map[string]interface{}{
		"workflowName":       workflowName,
		"sourceTask":         sourceTask,
		"targetTask":         targetTask,
		"straightEdge":       e.StraightEdge,
		"_resourceName":      SanitizeName(fmt.Sprintf("%s_%s_to_%s", workflowName, sourceTask, targetTask)),
		"_terraformResource": "stonebranch_workflow_edge",
	}

	tmpl := GetTemplate("workflow_edge")
	if tmpl == nil {
		return "", fmt.Errorf("no template for workflow_edge")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateWorkflowVertexHCLNew generates HCL for a workflow vertex using sequential naming.
// Returns the terraform resource name and the HCL string.
func (g *Generator) generateWorkflowVertexHCLNew(workflowName, workflowTfName string, v WorkflowVertex) (string, string, error) {
	taskName := v.Task.Value
	vertexTfName := g.generateResourceName("workflow_vertex")

	data := map[string]interface{}{
		"workflowName":       workflowName,
		"workflowTfName":     workflowTfName,
		"taskName":           taskName,
		"alias":              v.Alias,
		"_resourceName":      vertexTfName,
		"_terraformResource": "stonebranch_workflow_vertex",
	}

	tmpl := GetTemplate("workflow_vertex")
	if tmpl == nil {
		return "", "", fmt.Errorf("no template for workflow_vertex")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", err
	}

	return vertexTfName, buf.String(), nil
}

// generateWorkflowEdgeHCLNew generates HCL for a workflow edge using sequential naming.
func (g *Generator) generateWorkflowEdgeHCLNew(workflowName, workflowTfName string, e WorkflowEdge, vertices []WorkflowVertex, vertexTfNames map[string]string) (string, error) {
	// Find task names for source and target vertex IDs
	var sourceTask, targetTask string
	for _, v := range vertices {
		if v.VertexId == e.SourceId.Value {
			sourceTask = v.Task.Value
		}
		if v.VertexId == e.TargetId.Value {
			targetTask = v.Task.Value
		}
	}

	// Get the terraform resource names for source and target vertices
	sourceVertexTfName := vertexTfNames[e.SourceId.Value]
	targetVertexTfName := vertexTfNames[e.TargetId.Value]

	edgeTfName := g.generateResourceName("workflow_edge")

	data := map[string]interface{}{
		"workflowName":       workflowName,
		"workflowTfName":     workflowTfName,
		"sourceTask":         sourceTask,
		"targetTask":         targetTask,
		"sourceVertexTfName": sourceVertexTfName,
		"targetVertexTfName": targetVertexTfName,
		"straightEdge":       e.StraightEdge,
		"_resourceName":      edgeTfName,
		"_terraformResource": "stonebranch_workflow_edge",
	}

	tmpl := GetTemplate("workflow_edge")
	if tmpl == nil {
		return "", fmt.Errorf("no template for workflow_edge")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// outputHCL appends HCL to the buffer (legacy compatibility).
func (g *Generator) outputHCL(resourceType, name, hcl string) error {
	g.hclBuffer.WriteString(fmt.Sprintf("# %s: %s\n", resourceType, name))
	g.hclBuffer.WriteString(hcl)
	g.hclBuffer.WriteString("\n")
	return nil
}

// markExported marks a resource as exported.
func (g *Generator) markExported(resourceType, name string) {
	g.exported[fmt.Sprintf("%s/%s", resourceType, name)] = true
}

// isExported checks if a resource has already been exported.
func (g *Generator) isExported(resourceType, name string) bool {
	return g.exported[fmt.Sprintf("%s/%s", resourceType, name)]
}

// SanitizeName converts a StoneBranch resource name to a valid Terraform identifier.
func SanitizeName(name string) string {
	// Replace non-alphanumeric characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	sanitized := re.ReplaceAllString(name, "_")

	// Convert to lowercase
	sanitized = strings.ToLower(sanitized)

	// Remove consecutive underscores
	re = regexp.MustCompile(`_+`)
	sanitized = re.ReplaceAllString(sanitized, "_")

	// Trim leading/trailing underscores
	sanitized = strings.Trim(sanitized, "_")

	// Ensure it starts with a letter
	if len(sanitized) > 0 && unicode.IsDigit(rune(sanitized[0])) {
		sanitized = "r_" + sanitized
	}

	if sanitized == "" {
		sanitized = "resource"
	}

	return sanitized
}

// GetExportedResources returns the list of exported resources.
func (g *Generator) GetExportedResources() []string {
	var result []string
	for key := range g.exported {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

// GetTemplate returns the template for a resource type.
// This is a placeholder that will be implemented in templates.go
var GetTemplate func(resourceType string) *template.Template
