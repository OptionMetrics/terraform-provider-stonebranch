package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &WorkflowEdgeResource{}
)

func NewWorkflowEdgeResource() resource.Resource {
	return &WorkflowEdgeResource{}
}

// WorkflowEdgeResource defines the resource implementation.
type WorkflowEdgeResource struct {
	client *client.Client
}

// WorkflowEdgeResourceModel describes the resource data model.
type WorkflowEdgeResourceModel struct {
	// Identity - composite key
	WorkflowName types.String `tfsdk:"workflow_name"`
	SourceId     types.String `tfsdk:"source_id"` // vertex ID of source task
	TargetId     types.String `tfsdk:"target_id"` // vertex ID of target task

	// Optional
	StraightEdge types.Bool `tfsdk:"straight_edge"`
}

// WorkflowEdgeAPIModel represents the API request structure for creating an edge.
type WorkflowEdgeAPIModel struct {
	SourceId     *EdgeVertexRef `json:"sourceId,omitempty"`
	TargetId     *EdgeVertexRef `json:"targetId,omitempty"`
	StraightEdge bool           `json:"straightEdge,omitempty"`
}

// EdgeVertexRef represents a vertex reference for edge endpoints.
type EdgeVertexRef struct {
	Value string `json:"value,omitempty"` // vertexId
}

// WorkflowEdgeResponseModel represents the API response structure.
type WorkflowEdgeResponseModel struct {
	SysId        string                `json:"sysId,omitempty"`
	SourceId     *EdgeVertexRefResp    `json:"sourceId,omitempty"`
	TargetId     *EdgeVertexRefResp    `json:"targetId,omitempty"`
	StraightEdge bool                  `json:"straightEdge,omitempty"`
}

type EdgeVertexRefResp struct {
	TaskName  string `json:"taskName,omitempty"`
	TaskAlias string `json:"taskAlias,omitempty"`
	Value     string `json:"value,omitempty"`
}

func (r *WorkflowEdgeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_edge"
}

func (r *WorkflowEdgeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a dependency (edge) between tasks within a StoneBranch Workflow. Use this resource to define execution order between workflow vertices.",

		Attributes: map[string]schema.Attribute{
			"workflow_name": schema.StringAttribute{
				MarkdownDescription: "Name of the workflow containing the dependency.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_id": schema.StringAttribute{
				MarkdownDescription: "Vertex ID of the source task (the predecessor). This task must complete before the target task runs.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_id": schema.StringAttribute{
				MarkdownDescription: "Vertex ID of the target task (the successor). This task waits for the source task to complete.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"straight_edge": schema.BoolAttribute{
				MarkdownDescription: "Whether to draw the edge as a straight line in the workflow diagram.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

func (r *WorkflowEdgeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *WorkflowEdgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkflowEdgeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating workflow edge", map[string]any{
		"workflow":  data.WorkflowName.ValueString(),
		"source_id": data.SourceId.ValueString(),
		"target_id": data.TargetId.ValueString(),
	})

	// Build API model
	apiModel := &WorkflowEdgeAPIModel{
		SourceId: &EdgeVertexRef{
			Value: data.SourceId.ValueString(),
		},
		TargetId: &EdgeVertexRef{
			Value: data.TargetId.ValueString(),
		},
		StraightEdge: data.StraightEdge.ValueBool(),
	}

	// Add the edge to the workflow
	query := url.Values{}
	query.Set("workflowname", data.WorkflowName.ValueString())

	_, err := r.client.Post(ctx, "/resources/workflow/edges?"+query.Encode(), apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Workflow Edge",
			fmt.Sprintf("Could not create edge from %s to %s in workflow %s: %s",
				data.SourceId.ValueString(), data.TargetId.ValueString(), data.WorkflowName.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Created workflow edge", map[string]any{
		"source_id": data.SourceId.ValueString(),
		"target_id": data.TargetId.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowEdgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkflowEdgeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query for all edges in the workflow
	query := url.Values{}
	query.Set("workflowname", data.WorkflowName.ValueString())

	respBody, err := r.client.Get(ctx, "/resources/workflow/edges", query)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			tflog.Debug(ctx, "Workflow not found, removing edge from state", map[string]any{
				"workflow": data.WorkflowName.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Workflow Edges",
			fmt.Sprintf("Could not read edges for workflow %s: %s", data.WorkflowName.ValueString(), err),
		)
		return
	}

	// Parse response - API returns an array of edges
	var edges []WorkflowEdgeResponseModel
	if err := json.Unmarshal(respBody, &edges); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Response",
			fmt.Sprintf("Could not parse edges response: %s", err),
		)
		return
	}

	// Find the specific edge matching source and target
	found := false
	for _, edge := range edges {
		if edge.SourceId != nil && edge.TargetId != nil &&
			edge.SourceId.Value == data.SourceId.ValueString() &&
			edge.TargetId.Value == data.TargetId.ValueString() {
			found = true
			data.StraightEdge = types.BoolValue(edge.StraightEdge)
			break
		}
	}

	if !found {
		tflog.Debug(ctx, "Workflow edge not found, removing from state", map[string]any{
			"source_id": data.SourceId.ValueString(),
			"target_id": data.TargetId.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowEdgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkflowEdgeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating workflow edge", map[string]any{
		"source_id": data.SourceId.ValueString(),
		"target_id": data.TargetId.ValueString(),
	})

	// Build API model for update
	apiModel := &WorkflowEdgeAPIModel{
		SourceId: &EdgeVertexRef{
			Value: data.SourceId.ValueString(),
		},
		TargetId: &EdgeVertexRef{
			Value: data.TargetId.ValueString(),
		},
		StraightEdge: data.StraightEdge.ValueBool(),
	}

	query := url.Values{}
	query.Set("workflowname", data.WorkflowName.ValueString())

	_, err := r.client.Put(ctx, "/resources/workflow/edges?"+query.Encode(), apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Workflow Edge",
			fmt.Sprintf("Could not update edge from %s to %s: %s",
				data.SourceId.ValueString(), data.TargetId.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowEdgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkflowEdgeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Removing workflow edge", map[string]any{
		"workflow":  data.WorkflowName.ValueString(),
		"source_id": data.SourceId.ValueString(),
		"target_id": data.TargetId.ValueString(),
	})

	query := url.Values{}
	query.Set("workflowname", data.WorkflowName.ValueString())
	query.Set("sourceid", data.SourceId.ValueString())
	query.Set("targetid", data.TargetId.ValueString())

	_, err := r.client.Delete(ctx, "/resources/workflow/edges", query)
	if err != nil {
		// Ignore 404 errors (already deleted)
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Removing Workflow Edge",
			fmt.Sprintf("Could not remove edge from %s to %s in workflow %s: %s",
				data.SourceId.ValueString(), data.TargetId.ValueString(), data.WorkflowName.ValueString(), err),
		)
		return
	}
}
