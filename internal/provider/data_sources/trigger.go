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

	"github.com/OptionMetrics/terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TriggerDataSource{}

func NewTriggerDataSource() datasource.DataSource {
	return &TriggerDataSource{}
}

// TriggerDataSource defines the data source implementation.
type TriggerDataSource struct {
	client *client.Client
}

// TriggerDataSourceModel describes the data source data model.
type TriggerDataSourceModel struct {
	// Input (required)
	Name types.String `tfsdk:"name"`

	// Output - common fields across all trigger types
	SysId         types.String `tfsdk:"sys_id"`
	Type          types.String `tfsdk:"type"`
	Version       types.Int64  `tfsdk:"version"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Tasks         types.List   `tfsdk:"tasks"`
	TimeZone      types.String `tfsdk:"time_zone"`
	Calendar      types.String `tfsdk:"calendar"`
	OpswiseGroups types.List   `tfsdk:"opswise_groups"`
}

// TriggerAPIModel represents the API response structure (common fields).
type TriggerAPIModel struct {
	SysId         string   `json:"sysId"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Version       int64    `json:"version"`
	Description   string   `json:"description,omitempty"`
	Enabled       bool     `json:"enabled,omitempty"`
	Tasks         []string `json:"tasks,omitempty"`
	TimeZone      string   `json:"timeZone,omitempty"`
	Calendar      string   `json:"calendar,omitempty"`
	OpswiseGroups []string `json:"opswiseGroups,omitempty"`
}

func (d *TriggerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trigger"
}

func (d *TriggerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Looks up a trigger by name from StoneBranch Universal Controller. Returns common fields available across all trigger types.",

		Attributes: map[string]schema.Attribute{
			// Input
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the trigger to look up.",
				Required:            true,
			},

			// Output - common fields
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the trigger.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the trigger (e.g., 'triggerTime', 'triggerCron', 'triggerFm', 'triggerTm').",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the trigger.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the trigger.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the trigger is enabled.",
				Computed:            true,
			},
			"tasks": schema.ListAttribute{
				MarkdownDescription: "List of task names that this trigger executes.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"time_zone": schema.StringAttribute{
				MarkdownDescription: "Time zone for the trigger.",
				Computed:            true,
			},
			"calendar": schema.StringAttribute{
				MarkdownDescription: "Calendar used by the trigger.",
				Computed:            true,
			},
			"opswise_groups": schema.ListAttribute{
				MarkdownDescription: "List of business service names this trigger belongs to.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *TriggerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TriggerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TriggerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Looking up trigger", map[string]any{"name": data.Name.ValueString()})

	// Build query
	query := url.Values{}
	query.Set("triggername", data.Name.ValueString())

	// Fetch trigger
	respBody, err := d.client.Get(ctx, "/resources/trigger", query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Trigger",
			fmt.Sprintf("Could not read trigger %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	var apiModel TriggerAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Trigger Response",
			fmt.Sprintf("Could not parse trigger response: %s", err),
		)
		return
	}

	// Map API response to model
	data.SysId = types.StringValue(apiModel.SysId)
	data.Type = types.StringValue(apiModel.Type)
	data.Version = types.Int64Value(apiModel.Version)
	data.Enabled = types.BoolValue(apiModel.Enabled)

	if apiModel.Description != "" {
		data.Description = types.StringValue(apiModel.Description)
	} else {
		data.Description = types.StringNull()
	}

	if apiModel.TimeZone != "" {
		data.TimeZone = types.StringValue(apiModel.TimeZone)
	} else {
		data.TimeZone = types.StringNull()
	}

	if apiModel.Calendar != "" {
		data.Calendar = types.StringValue(apiModel.Calendar)
	} else {
		data.Calendar = types.StringNull()
	}

	// Handle tasks
	if len(apiModel.Tasks) > 0 {
		tasks, diags := types.ListValueFrom(ctx, types.StringType, apiModel.Tasks)
		resp.Diagnostics.Append(diags...)
		data.Tasks = tasks
	} else {
		data.Tasks = types.ListNull(types.StringType)
	}

	// Handle opswise_groups
	if len(apiModel.OpswiseGroups) > 0 {
		groups, diags := types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		resp.Diagnostics.Append(diags...)
		data.OpswiseGroups = groups
	} else {
		data.OpswiseGroups = types.ListNull(types.StringType)
	}

	tflog.Debug(ctx, "Found trigger", map[string]any{
		"sys_id": apiModel.SysId,
		"type":   apiModel.Type,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
