package data_sources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TaskInstancesDataSource{}

func NewTaskInstancesDataSource() datasource.DataSource {
	return &TaskInstancesDataSource{}
}

// TaskInstancesDataSource defines the data source implementation.
type TaskInstancesDataSource struct {
	client *client.Client
}

// TaskInstancesDataSourceModel describes the data source data model.
type TaskInstancesDataSourceModel struct {
	// Filter inputs
	TaskName             types.String `tfsdk:"task_name"`
	Status               types.String `tfsdk:"status"`
	Type                 types.String `tfsdk:"type"`
	AgentName            types.String `tfsdk:"agent_name"`
	UpdatedTimeType      types.String `tfsdk:"updated_time_type"`
	UpdatedTime          types.String `tfsdk:"updated_time"`
	WorkflowInstanceName types.String `tfsdk:"workflow_instance_name"`
	BusinessServices     types.String `tfsdk:"business_services"`

	// Output
	TaskInstances types.List `tfsdk:"task_instances"`
}

// TaskInstanceModel describes a single task instance in the results.
type TaskInstanceModel struct {
	SysId                  types.String `tfsdk:"sys_id"`
	Name                   types.String `tfsdk:"name"`
	Type                   types.String `tfsdk:"type"`
	Status                 types.String `tfsdk:"status"`
	StatusDescription      types.String `tfsdk:"status_description"`
	TriggerTime            types.String `tfsdk:"trigger_time"`
	StartTime              types.String `tfsdk:"start_time"`
	EndTime                types.String `tfsdk:"end_time"`
	ExitCode               types.String `tfsdk:"exit_code"`
	Agent                  types.String `tfsdk:"agent"`
	TaskName               types.String `tfsdk:"task_name"`
	TaskId                 types.String `tfsdk:"task_id"`
	InstanceNumber         types.Int64  `tfsdk:"instance_number"`
	TriggeredBy            types.String `tfsdk:"triggered_by"`
	WorkflowInstanceName   types.String `tfsdk:"workflow_instance_name"`
	WorkflowDefinitionName types.String `tfsdk:"workflow_definition_name"`
}

// TaskInstanceRequestAPIModel represents the API request structure.
type TaskInstanceRequestAPIModel struct {
	Name                 string `json:"name,omitempty"`
	Status               string `json:"status,omitempty"`
	Type                 string `json:"type,omitempty"`
	AgentName            string `json:"agentName,omitempty"`
	UpdatedTimeType      string `json:"updatedTimeType,omitempty"`
	UpdatedTime          string `json:"updatedTime,omitempty"`
	WorkflowInstanceName string `json:"workflowInstanceName,omitempty"`
	BusinessServices     string `json:"businessServices,omitempty"`
}

// TaskInstanceAPIModel represents the API response structure.
type TaskInstanceAPIModel struct {
	SysId                  string `json:"sysId"`
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	Status                 string `json:"status,omitempty"`
	StatusDescription      string `json:"statusDescription,omitempty"`
	TriggerTime            string `json:"triggerTime,omitempty"`
	StartTime              string `json:"startTime,omitempty"`
	EndTime                string `json:"endTime,omitempty"`
	ExitCode               string `json:"exitCode,omitempty"`
	Agent                  string `json:"agent,omitempty"`
	TaskName               string `json:"taskName,omitempty"`
	TaskId                 string `json:"taskId,omitempty"`
	InstanceNumber         int64  `json:"instanceNumber,omitempty"`
	TriggeredBy            string `json:"triggeredBy,omitempty"`
	WorkflowInstanceName   string `json:"workflowInstanceName,omitempty"`
	WorkflowDefinitionName string `json:"workflowDefinitionName,omitempty"`
}

func (d *TaskInstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task_instances"
}

func (d *TaskInstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of task instances from StoneBranch Universal Controller.",

		Attributes: map[string]schema.Attribute{
			// Filter inputs
			"task_name": schema.StringAttribute{
				MarkdownDescription: "Filter task instances by task name. Required by the API (use '*' for wildcard matching).",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Filter task instances by status (e.g., 'Running', 'Success', 'Failed', 'Waiting').",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Filter task instances by task type.",
				Optional:            true,
			},
			"agent_name": schema.StringAttribute{
				MarkdownDescription: "Filter task instances by agent name.",
				Optional:            true,
			},
			"updated_time_type": schema.StringAttribute{
				MarkdownDescription: "Type of time filter. Values: 'Today', 'Offset', 'Since', 'Older Than'.",
				Optional:            true,
			},
			"updated_time": schema.StringAttribute{
				MarkdownDescription: "Time value for filtering. Format depends on updated_time_type: 'mn' (minutes), 'h' (hours), 'd' (days). Example: '1h', '30mn', '2d'.",
				Optional:            true,
			},
			"workflow_instance_name": schema.StringAttribute{
				MarkdownDescription: "Filter task instances by parent workflow instance name.",
				Optional:            true,
			},
			"business_services": schema.StringAttribute{
				MarkdownDescription: "Filter task instances by business service name.",
				Optional:            true,
			},

			// Output
			"task_instances": schema.ListNestedAttribute{
				MarkdownDescription: "List of task instances matching the filter criteria.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sys_id": schema.StringAttribute{
							MarkdownDescription: "System ID of the task instance.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the task instance.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the task.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Current status of the task instance.",
							Computed:            true,
						},
						"status_description": schema.StringAttribute{
							MarkdownDescription: "Description of the current status.",
							Computed:            true,
						},
						"trigger_time": schema.StringAttribute{
							MarkdownDescription: "Time when the task was triggered.",
							Computed:            true,
						},
						"start_time": schema.StringAttribute{
							MarkdownDescription: "Time when the task started executing.",
							Computed:            true,
						},
						"end_time": schema.StringAttribute{
							MarkdownDescription: "Time when the task finished executing.",
							Computed:            true,
						},
						"exit_code": schema.StringAttribute{
							MarkdownDescription: "Exit code returned by the task.",
							Computed:            true,
						},
						"agent": schema.StringAttribute{
							MarkdownDescription: "Agent that executed the task.",
							Computed:            true,
						},
						"task_name": schema.StringAttribute{
							MarkdownDescription: "Name of the task definition.",
							Computed:            true,
						},
						"task_id": schema.StringAttribute{
							MarkdownDescription: "ID of the task definition.",
							Computed:            true,
						},
						"instance_number": schema.Int64Attribute{
							MarkdownDescription: "Instance number of the task execution.",
							Computed:            true,
						},
						"triggered_by": schema.StringAttribute{
							MarkdownDescription: "What triggered the task execution.",
							Computed:            true,
						},
						"workflow_instance_name": schema.StringAttribute{
							MarkdownDescription: "Name of the parent workflow instance.",
							Computed:            true,
						},
						"workflow_definition_name": schema.StringAttribute{
							MarkdownDescription: "Name of the parent workflow definition.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *TaskInstancesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TaskInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TaskInstancesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading task instances list")

	// Build request body from filters
	requestBody := TaskInstanceRequestAPIModel{}
	if !data.TaskName.IsNull() {
		requestBody.Name = data.TaskName.ValueString()
	}
	if !data.Status.IsNull() {
		requestBody.Status = data.Status.ValueString()
	}
	if !data.Type.IsNull() {
		requestBody.Type = data.Type.ValueString()
	}
	if !data.AgentName.IsNull() {
		requestBody.AgentName = data.AgentName.ValueString()
	}
	if !data.UpdatedTimeType.IsNull() {
		requestBody.UpdatedTimeType = data.UpdatedTimeType.ValueString()
	}
	if !data.UpdatedTime.IsNull() {
		requestBody.UpdatedTime = data.UpdatedTime.ValueString()
	}
	if !data.WorkflowInstanceName.IsNull() {
		requestBody.WorkflowInstanceName = data.WorkflowInstanceName.ValueString()
	}
	if !data.BusinessServices.IsNull() {
		requestBody.BusinessServices = data.BusinessServices.ValueString()
	}

	// Make POST API call
	respBody, err := d.client.Post(ctx, "/resources/taskinstance/listadv", requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Task Instances",
			fmt.Sprintf("Could not read task instances: %s", err),
		)
		return
	}

	// Parse response (array of task instances)
	var apiModels []TaskInstanceAPIModel
	if err := json.Unmarshal(respBody, &apiModels); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Response",
			fmt.Sprintf("Could not parse task instances response: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Read task instances", map[string]any{"count": len(apiModels)})

	// Convert to Terraform model
	taskInstances, diags := d.fromAPIModels(ctx, apiModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.TaskInstances = taskInstances

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// fromAPIModels converts a list of API models to a Terraform list.
func (d *TaskInstancesDataSource) fromAPIModels(ctx context.Context, apiModels []TaskInstanceAPIModel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	instanceType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"sys_id":                   types.StringType,
			"name":                     types.StringType,
			"type":                     types.StringType,
			"status":                   types.StringType,
			"status_description":       types.StringType,
			"trigger_time":             types.StringType,
			"start_time":               types.StringType,
			"end_time":                 types.StringType,
			"exit_code":                types.StringType,
			"agent":                    types.StringType,
			"task_name":                types.StringType,
			"task_id":                  types.StringType,
			"instance_number":          types.Int64Type,
			"triggered_by":             types.StringType,
			"workflow_instance_name":   types.StringType,
			"workflow_definition_name": types.StringType,
		},
	}

	if len(apiModels) == 0 {
		return types.ListValueMust(instanceType, []attr.Value{}), diags
	}

	instances := make([]attr.Value, len(apiModels))
	for i, apiModel := range apiModels {
		instanceObj, objDiags := types.ObjectValue(instanceType.AttrTypes, map[string]attr.Value{
			"sys_id":                   types.StringValue(apiModel.SysId),
			"name":                     types.StringValue(apiModel.Name),
			"type":                     types.StringValue(apiModel.Type),
			"status":                   stringValueOrNull(apiModel.Status),
			"status_description":       stringValueOrNull(apiModel.StatusDescription),
			"trigger_time":             stringValueOrNull(apiModel.TriggerTime),
			"start_time":               stringValueOrNull(apiModel.StartTime),
			"end_time":                 stringValueOrNull(apiModel.EndTime),
			"exit_code":                stringValueOrNull(apiModel.ExitCode),
			"agent":                    stringValueOrNull(apiModel.Agent),
			"task_name":                stringValueOrNull(apiModel.TaskName),
			"task_id":                  stringValueOrNull(apiModel.TaskId),
			"instance_number":          types.Int64Value(apiModel.InstanceNumber),
			"triggered_by":             stringValueOrNull(apiModel.TriggeredBy),
			"workflow_instance_name":   stringValueOrNull(apiModel.WorkflowInstanceName),
			"workflow_definition_name": stringValueOrNull(apiModel.WorkflowDefinitionName),
		})
		diags.Append(objDiags...)
		if diags.HasError() {
			return types.ListNull(instanceType), diags
		}
		instances[i] = instanceObj
	}

	return types.ListValueMust(instanceType, instances), diags
}
