package data_sources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TaskDataSource{}

func NewTaskDataSource() datasource.DataSource {
	return &TaskDataSource{}
}

// TaskDataSource defines the data source implementation.
type TaskDataSource struct {
	client *client.Client
}

// TaskDataSourceModel describes the data source data model.
type TaskDataSourceModel struct {
	// Input (required)
	Name types.String `tfsdk:"name"`

	// Output - common fields across all task types
	SysId       types.String `tfsdk:"sys_id"`
	Type        types.String `tfsdk:"type"`
	Version     types.Int64  `tfsdk:"version"`
	Summary     types.String `tfsdk:"summary"`
	Agent       types.String `tfsdk:"agent"`
	AgentCluster types.String `tfsdk:"agent_cluster"`
	Credentials types.String `tfsdk:"credentials"`
	OpswiseGroups types.List `tfsdk:"opswise_groups"`
}

// SingleTaskAPIModel represents the API response structure (common fields).
type SingleTaskAPIModel struct {
	SysId        string   `json:"sysId"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Version      int64    `json:"version"`
	Summary      string   `json:"summary,omitempty"`
	Agent        string   `json:"agent,omitempty"`
	AgentCluster string   `json:"agentCluster,omitempty"`
	Credentials  string   `json:"credentials,omitempty"`
	OpswiseGroups []string `json:"opswiseGroups,omitempty"`
}

func (d *TaskDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

func (d *TaskDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Looks up a task by name from StoneBranch Universal Controller. Returns common fields available across all task types.",

		Attributes: map[string]schema.Attribute{
			// Input
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the task to look up.",
				Required:            true,
			},

			// Output - common fields
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the task.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the task (e.g., 'taskUnix', 'taskWindows', 'taskSql', 'taskWorkflow').",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the task.",
				Computed:            true,
			},
			"summary": schema.StringAttribute{
				MarkdownDescription: "Summary/description of the task.",
				Computed:            true,
			},
			"agent": schema.StringAttribute{
				MarkdownDescription: "Name of the agent assigned to execute this task.",
				Computed:            true,
			},
			"agent_cluster": schema.StringAttribute{
				MarkdownDescription: "Name of the agent cluster assigned to execute this task.",
				Computed:            true,
			},
			"credentials": schema.StringAttribute{
				MarkdownDescription: "Name of the credentials used by this task.",
				Computed:            true,
			},
			"opswise_groups": schema.ListAttribute{
				MarkdownDescription: "List of business service names this task belongs to.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *TaskDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TaskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TaskDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Looking up task", map[string]any{"name": data.Name.ValueString()})

	// Build query
	query := url.Values{}
	query.Set("taskname", data.Name.ValueString())

	// Fetch task
	respBody, err := d.client.Get(ctx, "/resources/task", query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Task",
			fmt.Sprintf("Could not read task %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	var apiModel SingleTaskAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Task Response",
			fmt.Sprintf("Could not parse task response: %s", err),
		)
		return
	}

	// Map API response to model
	data.SysId = types.StringValue(apiModel.SysId)
	data.Type = types.StringValue(apiModel.Type)
	data.Version = types.Int64Value(apiModel.Version)

	if apiModel.Summary != "" {
		data.Summary = types.StringValue(apiModel.Summary)
	} else {
		data.Summary = types.StringNull()
	}

	if apiModel.Agent != "" {
		data.Agent = types.StringValue(apiModel.Agent)
	} else {
		data.Agent = types.StringNull()
	}

	if apiModel.AgentCluster != "" {
		data.AgentCluster = types.StringValue(apiModel.AgentCluster)
	} else {
		data.AgentCluster = types.StringNull()
	}

	if apiModel.Credentials != "" {
		data.Credentials = types.StringValue(apiModel.Credentials)
	} else {
		data.Credentials = types.StringNull()
	}

	// Handle opswise_groups
	if len(apiModel.OpswiseGroups) > 0 {
		groups, diags := types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		resp.Diagnostics.Append(diags...)
		data.OpswiseGroups = groups
	} else {
		data.OpswiseGroups = types.ListNull(types.StringType)
	}

	tflog.Debug(ctx, "Found task", map[string]any{
		"sys_id": apiModel.SysId,
		"type":   apiModel.Type,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
