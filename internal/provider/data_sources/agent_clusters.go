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
var _ datasource.DataSource = &AgentClustersDataSource{}

func NewAgentClustersDataSource() datasource.DataSource {
	return &AgentClustersDataSource{}
}

// AgentClustersDataSource defines the data source implementation.
type AgentClustersDataSource struct {
	client *client.Client
}

// AgentClustersDataSourceModel describes the data source data model.
type AgentClustersDataSourceModel struct {
	// Filter inputs
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	BusinessServices types.String `tfsdk:"business_services"`

	// Output
	AgentClusters types.List `tfsdk:"agent_clusters"`
}

// AgentClusterModel describes a single agent cluster in the results.
type AgentClusterModel struct {
	SysId         types.String `tfsdk:"sys_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Type          types.String `tfsdk:"type"`
	Version       types.Int64  `tfsdk:"version"`
	Distribution  types.String `tfsdk:"distribution"`
	Suspended     types.Bool   `tfsdk:"suspended"`
	LimitType     types.String `tfsdk:"limit_type"`
	LimitAmount   types.Int64  `tfsdk:"limit_amount"`
	OpswiseGroups types.List   `tfsdk:"opswise_groups"`
}

// AgentClusterAPIModel represents the API response structure.
type AgentClusterAPIModel struct {
	SysId         string   `json:"sysId"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Type          string   `json:"type"`
	Version       int64    `json:"version,omitempty"`
	Distribution  string   `json:"distribution,omitempty"`
	Suspended     bool     `json:"suspended,omitempty"`
	LimitType     string   `json:"limitType,omitempty"`
	LimitAmount   int64    `json:"limitAmount,omitempty"`
	OpswiseGroups []string `json:"opswiseGroups,omitempty"`
}

func (d *AgentClustersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_clusters"
}

func (d *AgentClustersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a list of agent clusters from StoneBranch Universal Controller.",

		Attributes: map[string]schema.Attribute{
			// Filter inputs
			"name": schema.StringAttribute{
				MarkdownDescription: "Filter agent clusters by name.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Filter agent clusters by type. Values: 'Windows', 'Linux/Unix'.",
				Optional:            true,
			},
			"business_services": schema.StringAttribute{
				MarkdownDescription: "Filter agent clusters by business service name.",
				Optional:            true,
			},

			// Output
			"agent_clusters": schema.ListNestedAttribute{
				MarkdownDescription: "List of agent clusters matching the filter criteria.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sys_id": schema.StringAttribute{
							MarkdownDescription: "System ID of the agent cluster.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the agent cluster.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the agent cluster.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the agent cluster (Windows, Linux/Unix).",
							Computed:            true,
						},
						"version": schema.Int64Attribute{
							MarkdownDescription: "Version number of the agent cluster.",
							Computed:            true,
						},
						"distribution": schema.StringAttribute{
							MarkdownDescription: "Distribution method for task assignment.",
							Computed:            true,
						},
						"suspended": schema.BoolAttribute{
							MarkdownDescription: "Whether the agent cluster is suspended.",
							Computed:            true,
						},
						"limit_type": schema.StringAttribute{
							MarkdownDescription: "Type of limit applied to the cluster.",
							Computed:            true,
						},
						"limit_amount": schema.Int64Attribute{
							MarkdownDescription: "Limit amount for the cluster.",
							Computed:            true,
						},
						"opswise_groups": schema.ListAttribute{
							MarkdownDescription: "List of business service names this agent cluster belongs to.",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *AgentClustersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AgentClustersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AgentClustersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading agent clusters list")

	// Build query parameters from filters
	query := url.Values{}
	if !data.Name.IsNull() {
		query.Set("agentclustername", data.Name.ValueString())
	}
	if !data.Type.IsNull() {
		query.Set("type", data.Type.ValueString())
	}
	if !data.BusinessServices.IsNull() {
		query.Set("businessServices", data.BusinessServices.ValueString())
	}

	// Make API call
	respBody, err := d.client.Get(ctx, "/resources/agentcluster/listadv", query)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Agent Clusters",
			fmt.Sprintf("Could not read agent clusters: %s", err),
		)
		return
	}

	// Parse response (array of agent clusters)
	var apiModels []AgentClusterAPIModel
	if err := json.Unmarshal(respBody, &apiModels); err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Response",
			fmt.Sprintf("Could not parse agent clusters response: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Read agent clusters", map[string]any{"count": len(apiModels)})

	// Convert to Terraform model
	agentClusters, diags := d.fromAPIModels(ctx, apiModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.AgentClusters = agentClusters

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// fromAPIModels converts a list of API models to a Terraform list.
func (d *AgentClustersDataSource) fromAPIModels(ctx context.Context, apiModels []AgentClusterAPIModel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	clusterType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"sys_id":         types.StringType,
			"name":           types.StringType,
			"description":    types.StringType,
			"type":           types.StringType,
			"version":        types.Int64Type,
			"distribution":   types.StringType,
			"suspended":      types.BoolType,
			"limit_type":     types.StringType,
			"limit_amount":   types.Int64Type,
			"opswise_groups": types.ListType{ElemType: types.StringType},
		},
	}

	if len(apiModels) == 0 {
		return types.ListValueMust(clusterType, []attr.Value{}), diags
	}

	clusters := make([]attr.Value, len(apiModels))
	for i, apiModel := range apiModels {
		// Handle opswise_groups
		var opswiseGroups types.List
		if len(apiModel.OpswiseGroups) > 0 {
			opswiseGroups, _ = types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		} else {
			opswiseGroups = types.ListValueMust(types.StringType, []attr.Value{})
		}

		clusterObj, objDiags := types.ObjectValue(clusterType.AttrTypes, map[string]attr.Value{
			"sys_id":         types.StringValue(apiModel.SysId),
			"name":           types.StringValue(apiModel.Name),
			"description":    stringValueOrNull(apiModel.Description),
			"type":           types.StringValue(apiModel.Type),
			"version":        types.Int64Value(apiModel.Version),
			"distribution":   stringValueOrNull(apiModel.Distribution),
			"suspended":      types.BoolValue(apiModel.Suspended),
			"limit_type":     stringValueOrNull(apiModel.LimitType),
			"limit_amount":   types.Int64Value(apiModel.LimitAmount),
			"opswise_groups": opswiseGroups,
		})
		diags.Append(objDiags...)
		if diags.HasError() {
			return types.ListNull(clusterType), diags
		}
		clusters[i] = clusterObj
	}

	return types.ListValueMust(clusterType, clusters), diags
}
