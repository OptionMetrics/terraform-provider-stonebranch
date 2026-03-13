package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/OptionMetrics/terraform-provider-stonebranch/cmd/sb2tf/generator"
)

var (
	listFilter string

	listCmd = &cobra.Command{
		Use:   "list [resource-type]",
		Short: "List available resources",
		Long: `List resources from StoneBranch Universal Controller.

Without arguments, shows available resource types.
With a resource type, lists resources of that type.

Examples:
  sb2tf list                    # Show all resource types
  sb2tf list tasks              # List all tasks
  sb2tf list task_unix          # List Unix tasks only
  sb2tf list triggers           # List all triggers
  sb2tf list --filter "prod-*"  # Filter by name pattern`,
		Args: cobra.MaximumNArgs(1),
		RunE: runList,
	}
)

func init() {
	listCmd.Flags().StringVar(&listFilter, "filter", "", "Filter resources by name pattern (supports wildcards)")
}

func runList(cmd *cobra.Command, args []string) error {
	// No arguments: show available resource types
	if len(args) == 0 {
		return listResourceTypes()
	}

	resourceType := args[0]

	// Handle category shortcuts (tasks, triggers, connections, etc.)
	switch resourceType {
	case "tasks":
		return listAllTasks()
	case "triggers":
		return listAllTriggers()
	case "connections":
		return listAllConnections()
	default:
		return listSpecificType(resourceType)
	}
}

func listResourceTypes() error {
	fmt.Println("Available resource types:")
	fmt.Println()

	categories := generator.GetResourceCategories()
	for _, cat := range categories {
		fmt.Printf("  %s:\n", cat.Name)
		for _, rt := range cat.Types {
			fmt.Printf("    %-25s %s\n", rt.CLIName, rt.TerraformResource)
		}
		fmt.Println()
	}

	fmt.Println("Shortcuts:")
	fmt.Println("  tasks       - List all task types")
	fmt.Println("  triggers    - List all trigger types")
	fmt.Println("  connections - List all connection types")

	return nil
}

func listAllTasks() error {
	ctx := context.Background()
	client := GetClient()

	query := url.Values{}
	if listFilter != "" {
		query.Set("taskname", listFilter)
	}

	respBody, err := client.Get(ctx, "/resources/task/listadv", query)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	var tasks []ResourceItem
	if err := json.Unmarshal(respBody, &tasks); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	printResourceTable("Tasks", tasks, true)
	return nil
}

func listAllTriggers() error {
	ctx := context.Background()
	client := GetClient()

	query := url.Values{}
	if listFilter != "" {
		query.Set("triggername", listFilter)
	}

	respBody, err := client.Get(ctx, "/resources/trigger/listadv", query)
	if err != nil {
		return fmt.Errorf("failed to list triggers: %w", err)
	}

	var triggers []ResourceItem
	if err := json.Unmarshal(respBody, &triggers); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	printResourceTable("Triggers", triggers, true)
	return nil
}

func listAllConnections() error {
	ctx := context.Background()
	client := GetClient()

	var allConnections []ResourceItem

	// Database connections
	respBody, err := client.Get(ctx, "/resources/databaseconnection/list", nil)
	if err == nil {
		var dbConns []ResourceItem
		if json.Unmarshal(respBody, &dbConns) == nil {
			for i := range dbConns {
				dbConns[i].Type = "database_connection"
			}
			allConnections = append(allConnections, dbConns...)
		}
	}

	// Email connections
	respBody, err = client.Get(ctx, "/resources/emailconnection/list", nil)
	if err == nil {
		var emailConns []ResourceItem
		if json.Unmarshal(respBody, &emailConns) == nil {
			for i := range emailConns {
				emailConns[i].Type = "email_connection"
			}
			allConnections = append(allConnections, emailConns...)
		}
	}

	// Filter if needed
	if listFilter != "" {
		allConnections = filterByName(allConnections, listFilter)
	}

	printResourceTable("Connections", allConnections, true)
	return nil
}

func listSpecificType(resourceType string) error {
	rt := generator.GetResourceType(resourceType)
	if rt == nil {
		return fmt.Errorf("unknown resource type: %s\nRun 'sb2tf list' to see available types", resourceType)
	}

	ctx := context.Background()
	client := GetClient()

	query := url.Values{}
	// Don't filter by type in API - filter locally instead
	if listFilter != "" {
		query.Set(rt.NameQueryParam, listFilter)
	}

	respBody, err := client.Get(ctx, rt.ListEndpoint, query)
	if err != nil {
		return fmt.Errorf("failed to list %s: %w", resourceType, err)
	}

	// Parse as raw JSON to handle different field names
	var rawItems []map[string]interface{}
	if err := json.Unmarshal(respBody, &rawItems); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Determine which field contains the name
	nameField := rt.NameField
	if nameField == "" {
		nameField = "name"
	}

	// Convert to ResourceItem
	var items []ResourceItem
	for _, raw := range rawItems {
		item := ResourceItem{}
		if name, ok := raw[nameField].(string); ok {
			item.Name = name
		}
		if t, ok := raw["type"].(string); ok {
			item.Type = t
		}
		if summary, ok := raw["summary"].(string); ok {
			item.Summary = summary
		}
		if sysId, ok := raw["sysId"].(string); ok {
			item.SysId = sysId
		}
		items = append(items, item)
	}

	// Filter by type locally if needed
	if rt.APITypeValue != "" {
		var filtered []ResourceItem
		for _, item := range items {
			if item.Type == rt.APITypeValue {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	// Set the type for display if not returned by API
	for i := range items {
		if items[i].Type == "" {
			items[i].Type = resourceType
		}
	}

	printResourceTable(strings.Title(resourceType), items, rt.HasTypeField)
	return nil
}

// ResourceItem represents a generic resource from the API.
type ResourceItem struct {
	SysId   string `json:"sysId"`
	Name    string `json:"name"`
	Type    string `json:"type,omitempty"`
	Summary string `json:"summary,omitempty"`
}

func printResourceTable(title string, items []ResourceItem, showType bool) {
	if len(items) == 0 {
		fmt.Printf("No %s found.\n", strings.ToLower(title))
		return
	}

	// Sort by name
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	fmt.Printf("%s (%d):\n", title, len(items))
	fmt.Println()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if showType {
		fmt.Fprintln(w, "NAME\tTYPE\tSUMMARY")
		fmt.Fprintln(w, "----\t----\t-------")
		for _, item := range items {
			summary := truncate(item.Summary, 50)
			fmt.Fprintf(w, "%s\t%s\t%s\n", item.Name, item.Type, summary)
		}
	} else {
		fmt.Fprintln(w, "NAME\tSUMMARY")
		fmt.Fprintln(w, "----\t-------")
		for _, item := range items {
			summary := truncate(item.Summary, 60)
			fmt.Fprintf(w, "%s\t%s\n", item.Name, summary)
		}
	}
	w.Flush()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func filterByName(items []ResourceItem, pattern string) []ResourceItem {
	// Simple wildcard matching (only supports * at beginning and/or end)
	pattern = strings.ToLower(pattern)
	prefix := strings.HasPrefix(pattern, "*")
	suffix := strings.HasSuffix(pattern, "*")
	pattern = strings.Trim(pattern, "*")

	var result []ResourceItem
	for _, item := range items {
		name := strings.ToLower(item.Name)
		match := false
		if prefix && suffix {
			match = strings.Contains(name, pattern)
		} else if prefix {
			match = strings.HasSuffix(name, pattern)
		} else if suffix {
			match = strings.HasPrefix(name, pattern)
		} else {
			match = name == pattern
		}
		if match {
			result = append(result, item)
		}
	}
	return result
}
