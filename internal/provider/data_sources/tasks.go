package data_sources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TasksDataSource{}

func NewTasksDataSource() datasource.DataSource {
	return &TasksDataSource{}
}

// TasksDataSource defines the data source implementation.
type TasksDataSource struct {
	client *client.Client
}

// TasksDataSourceModel describes the data source data model.
type TasksDataSourceModel struct {
	// Filter inputs
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	AgentName        types.String `tfsdk:"agent_name"`
	BusinessServices types.String `tfsdk:"business_services"`
	WorkflowName     types.String `tfsdk:"workflow_name"`

	// Output
	Tasks types.List `tfsdk:"tasks"`
}

// TaskModel describes a single task in the results.
type TaskModel struct {
	SysId         types.String `tfsdk:"sys_id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	Summary       types.String `tfsdk:"summary"`
	Version       types.Int64  `tfsdk:"version"`
	Agent         types.String `tfsdk:"agent"`
	AgentCluster  types.String `tfsdk:"agent_cluster"`
	Credentials   types.String `tfsdk:"credentials"`
	OpswiseGroups types.List   `tfsdk:"opswise_groups"`
}

// TaskAPIModel represents the API response structure.
type TaskAPIModel struct {
	SysId         string   `json:"sysId"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Summary       string   `json:"summary,omitempty"`
	Version       int64    `json:"version,omitempty"`
	Agent         string   `json:"agent,omitempty"`
	AgentCluster  string   `json:"agentCluster,omitempty"`
	Credentials   string   `json:"credentials,omitempty"`
	OpswiseGroups []string `json:"opswiseGroups,omitempty"`
}

func (d *TasksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tasks"
}

func (d *TasksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of tasks from StoneBranch Universal Controller.",

		Attributes: map[string]schema.Attribute{
			// Filter inputs
			"name": schema.StringAttribute{
				MarkdownDescription: "Filter tasks by name (supports wildcards).",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Filter tasks by type. Values: 'taskUnix', 'taskWindows', 'taskSql', 'taskEmail', 'taskWorkflow', 'taskFileMonitor', 'taskFileTransfer', etc.",
				Optional:            true,
			},
			"agent_name": schema.StringAttribute{
				MarkdownDescription: "Filter tasks by agent name.",
				Optional:            true,
			},
			"business_services": schema.StringAttribute{
				MarkdownDescription: "Filter tasks by business service name.",
				Optional:            true,
			},
			"workflow_name": schema.StringAttribute{
				MarkdownDescription: "Filter tasks by workflow membership (returns tasks within the workflow).",
				Optional:            true,
			},

			// Output
			"tasks": schema.ListNestedAttribute{
				MarkdownDescription: "List of tasks matching the filter criteria.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sys_id": schema.StringAttribute{
							MarkdownDescription: "System ID of the task.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the task.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the task.",
							Computed:            true,
						},
						"summary": schema.StringAttribute{
							MarkdownDescription: "Summary/description of the task.",
							Computed:            true,
						},
						"version": schema.Int64Attribute{
							MarkdownDescription: "Version number of the task.",
							Computed:            true,
						},
						"agent": schema.StringAttribute{
							MarkdownDescription: "Agent assigned to the task.",
							Computed:            true,
						},
						"agent_cluster": schema.StringAttribute{
							MarkdownDescription: "Agent cluster assigned to the task.",
							Computed:            true,
						},
						"credentials": schema.StringAttribute{
							MarkdownDescription: "Credentials used by the task.",
							Computed:            true,
						},
						"opswise_groups": schema.ListAttribute{
							MarkdownDescription: "List of business service names this task belongs to.",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *TasksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *TasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TasksDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading tasks list")

	// Build query parameters from filters
	query := url.Values{}
	if !data.Name.IsNull() {
		query.Set("taskname", data.Name.ValueString())
	}
	if !data.Type.IsNull() {
		query.Set("type", data.Type.ValueString())
	}
	if !data.AgentName.IsNull() {
		query.Set("agentname", data.AgentName.ValueString())
	}
	if !data.BusinessServices.IsNull() {
		query.Set("businessServices", data.BusinessServices.ValueString())
	}
	if !data.WorkflowName.IsNull() {
		query.Set("workflowname", data.WorkflowName.ValueString())
	}

	// Make API call
	respBody, err := d.client.Get(ctx, "/resources/task/listadv", query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tasks",
			fmt.Sprintf("Could not read tasks: %s", err),
		)
		return
	}

	// Parse response (array of tasks)
	var apiModels []TaskAPIModel
	if err := json.Unmarshal(respBody, &apiModels); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Response",
			fmt.Sprintf("Could not parse tasks response: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Read tasks", map[string]any{"count": len(apiModels)})

	// Convert to Terraform model
	tasks, diags := d.fromAPIModels(ctx, apiModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Tasks = tasks

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// fromAPIModels converts a list of API models to a Terraform list.
func (d *TasksDataSource) fromAPIModels(ctx context.Context, apiModels []TaskAPIModel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	taskType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"sys_id":         types.StringType,
			"name":           types.StringType,
			"type":           types.StringType,
			"summary":        types.StringType,
			"version":        types.Int64Type,
			"agent":          types.StringType,
			"agent_cluster":  types.StringType,
			"credentials":    types.StringType,
			"opswise_groups": types.ListType{ElemType: types.StringType},
		},
	}

	if len(apiModels) == 0 {
		return types.ListValueMust(taskType, []attr.Value{}), diags
	}

	tasks := make([]attr.Value, len(apiModels))
	for i, apiModel := range apiModels {
		// Handle opswise_groups
		var opswiseGroups types.List
		if len(apiModel.OpswiseGroups) > 0 {
			opswiseGroups, _ = types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		} else {
			opswiseGroups = types.ListValueMust(types.StringType, []attr.Value{})
		}

		taskObj, objDiags := types.ObjectValue(taskType.AttrTypes, map[string]attr.Value{
			"sys_id":         types.StringValue(apiModel.SysId),
			"name":           types.StringValue(apiModel.Name),
			"type":           types.StringValue(apiModel.Type),
			"summary":        stringValueOrNull(apiModel.Summary),
			"version":        types.Int64Value(apiModel.Version),
			"agent":          stringValueOrNull(apiModel.Agent),
			"agent_cluster":  stringValueOrNull(apiModel.AgentCluster),
			"credentials":    stringValueOrNull(apiModel.Credentials),
			"opswise_groups": opswiseGroups,
		})
		diags.Append(objDiags...)
		if diags.HasError() {
			return types.ListNull(taskType), diags
		}
		tasks[i] = taskObj
	}

	return types.ListValueMust(taskType, tasks), diags
}
