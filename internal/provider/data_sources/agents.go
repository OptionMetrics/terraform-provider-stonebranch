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

	"github.com/OptionMetrics/terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AgentsDataSource{}

func NewAgentsDataSource() datasource.DataSource {
	return &AgentsDataSource{}
}

// AgentsDataSource defines the data source implementation.
type AgentsDataSource struct {
	client *client.Client
}

// AgentsDataSourceModel describes the data source data model.
type AgentsDataSourceModel struct {
	// Filter inputs
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	BusinessServices types.String `tfsdk:"business_services"`

	// Output
	Agents types.List `tfsdk:"agents"`
}

// AgentModel describes a single agent in the results.
type AgentModel struct {
	SysId          types.String `tfsdk:"sys_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	HostName       types.String `tfsdk:"host_name"`
	IpAddress      types.String `tfsdk:"ip_address"`
	Status         types.String `tfsdk:"status"`
	Version        types.String `tfsdk:"version"`
	Os             types.String `tfsdk:"os"`
	OsRelease      types.String `tfsdk:"os_release"`
	CpuLoad        types.Int64  `tfsdk:"cpu_load"`
	Suspended      types.Bool   `tfsdk:"suspended"`
	Decommissioned types.Bool   `tfsdk:"decommissioned"`
	OpswiseGroups  types.List   `tfsdk:"opswise_groups"`
}

// AgentAPIModel represents the API response structure.
type AgentAPIModel struct {
	SysId          string   `json:"sysId"`
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	Type           string   `json:"type"`
	HostName       string   `json:"hostName,omitempty"`
	IpAddress      string   `json:"ipAddress,omitempty"`
	Status         string   `json:"status,omitempty"`
	Version        string   `json:"version,omitempty"`
	Os             string   `json:"os,omitempty"`
	OsRelease      string   `json:"osRelease,omitempty"`
	CpuLoad        int32    `json:"cpuLoad,omitempty"`
	Suspended      bool     `json:"suspended,omitempty"`
	Decommissioned bool     `json:"decommissioned,omitempty"`
	OpswiseGroups  []string `json:"opswiseGroups,omitempty"`
}

func (d *AgentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agents"
}

func (d *AgentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of agents from StoneBranch Universal Controller.",

		Attributes: map[string]schema.Attribute{
			// Filter inputs
			"name": schema.StringAttribute{
				MarkdownDescription: "Filter agents by name (supports wildcards).",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Filter agents by type. Values: 'Windows', 'Linux/Unix', 'z/OS'.",
				Optional:            true,
			},
			"business_services": schema.StringAttribute{
				MarkdownDescription: "Filter agents by business service name.",
				Optional:            true,
			},

			// Output
			"agents": schema.ListNestedAttribute{
				MarkdownDescription: "List of agents matching the filter criteria.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sys_id": schema.StringAttribute{
							MarkdownDescription: "System ID of the agent.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the agent.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the agent.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the agent (Windows, Linux/Unix, z/OS).",
							Computed:            true,
						},
						"host_name": schema.StringAttribute{
							MarkdownDescription: "Hostname of the agent machine.",
							Computed:            true,
						},
						"ip_address": schema.StringAttribute{
							MarkdownDescription: "IP address of the agent.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Current status of the agent.",
							Computed:            true,
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "Agent software version.",
							Computed:            true,
						},
						"os": schema.StringAttribute{
							MarkdownDescription: "Operating system of the agent.",
							Computed:            true,
						},
						"os_release": schema.StringAttribute{
							MarkdownDescription: "Operating system release version.",
							Computed:            true,
						},
						"cpu_load": schema.Int64Attribute{
							MarkdownDescription: "Current CPU load on the agent.",
							Computed:            true,
						},
						"suspended": schema.BoolAttribute{
							MarkdownDescription: "Whether the agent is suspended.",
							Computed:            true,
						},
						"decommissioned": schema.BoolAttribute{
							MarkdownDescription: "Whether the agent is decommissioned.",
							Computed:            true,
						},
						"opswise_groups": schema.ListAttribute{
							MarkdownDescription: "List of business service names this agent belongs to.",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *AgentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AgentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AgentsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading agents list")

	// Build query parameters from filters
	query := url.Values{}
	if !data.Name.IsNull() {
		query.Set("agentname", data.Name.ValueString())
	}
	if !data.Type.IsNull() {
		query.Set("type", data.Type.ValueString())
	}
	if !data.BusinessServices.IsNull() {
		query.Set("businessServices", data.BusinessServices.ValueString())
	}

	// Make API call
	respBody, err := d.client.Get(ctx, "/resources/agent/listadv", query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Agents",
			fmt.Sprintf("Could not read agents: %s", err),
		)
		return
	}

	// Parse response (array of agents)
	var apiModels []AgentAPIModel
	if err := json.Unmarshal(respBody, &apiModels); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Response",
			fmt.Sprintf("Could not parse agents response: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Read agents", map[string]any{"count": len(apiModels)})

	// Convert to Terraform model
	agents, diags := d.fromAPIModels(ctx, apiModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Agents = agents

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// fromAPIModels converts a list of API models to a Terraform list.
func (d *AgentsDataSource) fromAPIModels(ctx context.Context, apiModels []AgentAPIModel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	agentType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"sys_id":         types.StringType,
			"name":           types.StringType,
			"description":    types.StringType,
			"type":           types.StringType,
			"host_name":      types.StringType,
			"ip_address":     types.StringType,
			"status":         types.StringType,
			"version":        types.StringType,
			"os":             types.StringType,
			"os_release":     types.StringType,
			"cpu_load":       types.Int64Type,
			"suspended":      types.BoolType,
			"decommissioned": types.BoolType,
			"opswise_groups": types.ListType{ElemType: types.StringType},
		},
	}

	if len(apiModels) == 0 {
		return types.ListValueMust(agentType, []attr.Value{}), diags
	}

	agents := make([]attr.Value, len(apiModels))
	for i, apiModel := range apiModels {
		// Handle opswise_groups
		var opswiseGroups types.List
		if len(apiModel.OpswiseGroups) > 0 {
			opswiseGroups, _ = types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		} else {
			opswiseGroups = types.ListValueMust(types.StringType, []attr.Value{})
		}

		agentObj, objDiags := types.ObjectValue(agentType.AttrTypes, map[string]attr.Value{
			"sys_id":         types.StringValue(apiModel.SysId),
			"name":           types.StringValue(apiModel.Name),
			"description":    stringValueOrNull(apiModel.Description),
			"type":           types.StringValue(apiModel.Type),
			"host_name":      stringValueOrNull(apiModel.HostName),
			"ip_address":     stringValueOrNull(apiModel.IpAddress),
			"status":         stringValueOrNull(apiModel.Status),
			"version":        stringValueOrNull(apiModel.Version),
			"os":             stringValueOrNull(apiModel.Os),
			"os_release":     stringValueOrNull(apiModel.OsRelease),
			"cpu_load":       types.Int64Value(int64(apiModel.CpuLoad)),
			"suspended":      types.BoolValue(apiModel.Suspended),
			"decommissioned": types.BoolValue(apiModel.Decommissioned),
			"opswise_groups": opswiseGroups,
		})
		diags.Append(objDiags...)
		if diags.HasError() {
			return types.ListNull(agentType), diags
		}
		agents[i] = agentObj
	}

	return types.ListValueMust(agentType, agents), diags
}

// stringValueOrNull returns a StringValue if s is non-empty, otherwise StringNull.
func stringValueOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
