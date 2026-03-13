package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/OptionMetrics/terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &AgentClusterResource{}
	_ resource.ResourceWithImportState = &AgentClusterResource{}
)

func NewAgentClusterResource() resource.Resource {
	return &AgentClusterResource{}
}

// AgentClusterResource defines the resource implementation.
type AgentClusterResource struct {
	client *client.Client
}

// AgentClusterResourceModel describes the resource data model.
type AgentClusterResourceModel struct {
	// Identity
	SysId   types.String `tfsdk:"sys_id"`
	Name    types.String `tfsdk:"name"`
	Version types.Int64  `tfsdk:"version"`

	// Basic
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`

	// Distribution
	Distribution types.String `tfsdk:"distribution"`

	// Network alias (for load balancing)
	NetworkAlias     types.String `tfsdk:"network_alias"`
	NetworkAliasPort types.Int64  `tfsdk:"network_alias_port"`

	// Task execution limits
	LimitType   types.String `tfsdk:"limit_type"`
	LimitAmount types.Int64  `tfsdk:"limit_amount"`

	// Per-agent limits
	AgentLimitType   types.String `tfsdk:"agent_limit_type"`
	AgentLimitAmount types.Int64  `tfsdk:"agent_limit_amount"`

	// Agent selection options
	IgnoreInactiveAgents  types.Bool `tfsdk:"ignore_inactive_agents"`
	IgnoreSuspendedAgents types.Bool `tfsdk:"ignore_suspended_agents"`

	// Business services
	OpswiseGroups types.List `tfsdk:"opswise_groups"`
}

// AgentClusterAPIModel represents the API request/response structure.
type AgentClusterAPIModel struct {
	SysId                 string   `json:"sysId,omitempty"`
	Name                  string   `json:"name"`
	Version               int64    `json:"version,omitempty"`
	Type                  string   `json:"type"`
	Description           string   `json:"description,omitempty"`
	Distribution          string   `json:"distribution,omitempty"`
	NetworkAlias          string   `json:"networkAlias,omitempty"`
	NetworkAliasPort      int64    `json:"networkAliasPort,omitempty"`
	LimitType             string   `json:"limitType,omitempty"`
	LimitAmount           int64    `json:"limitAmount,omitempty"`
	AgentLimitType        string   `json:"agentLimitType,omitempty"`
	AgentLimitAmount      int64    `json:"agentLimitAmount,omitempty"`
	IgnoreInactiveAgents  bool     `json:"ignoreInactiveAgents,omitempty"`
	IgnoreSuspendedAgents bool     `json:"ignoreSuspendedAgents,omitempty"`
	OpswiseGroups         []string `json:"opswiseGroups,omitempty"`
}

// Type mapping between user-friendly names and API discriminator values
var agentClusterTypeToAPI = map[string]string{
	"Linux/Unix": "unixAgentCluster",
	"Windows":    "windowsAgentCluster",
	"z/OS":       "ibmiAgentCluster",
}

var agentClusterTypeFromAPI = map[string]string{
	"unixAgentCluster":    "Linux/Unix",
	"windowsAgentCluster": "Windows",
	"ibmiAgentCluster":    "z/OS",
}

func (r *AgentClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_cluster"
}

func (r *AgentClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a StoneBranch Agent Cluster. Agent clusters group multiple agents together for load distribution and high availability.",

		Attributes: map[string]schema.Attribute{
			// Identity
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the agent cluster (assigned by StoneBranch).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the agent cluster.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the agent cluster (for optimistic locking).",
				Computed:            true,
			},

			// Basic
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of agents in the cluster. Values: 'Linux/Unix', 'Windows'.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the agent cluster.",
				Optional:            true,
			},

			// Distribution
			"distribution": schema.StringAttribute{
				MarkdownDescription: "How tasks are distributed to agents in the cluster. Values: 'Any', 'All', 'Lowest CPU Utilization', 'Round Robin'.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Any"),
			},

			// Network alias
			"network_alias": schema.StringAttribute{
				MarkdownDescription: "Network alias hostname for the cluster (used for network-level load balancing).",
				Optional:            true,
			},
			"network_alias_port": schema.Int64Attribute{
				MarkdownDescription: "Port number for the network alias.",
				Optional:            true,
			},

			// Task execution limits
			"limit_type": schema.StringAttribute{
				MarkdownDescription: "Type of task execution limit. Values: 'Unlimited', 'Limited'. Defaults to 'Unlimited'.",
				Optional:            true,
				Computed:            true,
			},
			"limit_amount": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of concurrent tasks across the cluster (when limit_type is 'Limited'). Server default is 5.",
				Optional:            true,
				Computed:            true,
			},

			// Per-agent limits
			"agent_limit_type": schema.StringAttribute{
				MarkdownDescription: "Type of per-agent task limit. Values: 'Unlimited', 'Limited'. Defaults to 'Unlimited'.",
				Optional:            true,
				Computed:            true,
			},
			"agent_limit_amount": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of concurrent tasks per agent (when agent_limit_type is 'Limited'). Server default is 10.",
				Optional:            true,
				Computed:            true,
			},

			// Agent selection options
			"ignore_inactive_agents": schema.BoolAttribute{
				MarkdownDescription: "Whether to skip inactive agents when distributing tasks.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"ignore_suspended_agents": schema.BoolAttribute{
				MarkdownDescription: "Whether to skip suspended agents when distributing tasks.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},

			// Business services
			"opswise_groups": schema.ListAttribute{
				MarkdownDescription: "List of business service names this agent cluster belongs to.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *AgentClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AgentClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AgentClusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating agent cluster", map[string]any{"name": data.Name.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Create the agent cluster
	_, err := r.client.Post(ctx, "/resources/agentcluster", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Agent Cluster",
			fmt.Sprintf("Could not create agent cluster %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back the created agent cluster to get sysId and other computed fields
	err = r.readAgentCluster(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Agent Cluster",
			fmt.Sprintf("Could not read agent cluster %s after creation: %s", data.Name.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Created agent cluster", map[string]any{"sys_id": data.SysId.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AgentClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AgentClusterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.readAgentCluster(ctx, &data)
	if err != nil {
		// Check if agent cluster was deleted outside of Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			tflog.Debug(ctx, "Agent cluster not found, removing from state", map[string]any{"name": data.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Agent Cluster",
			fmt.Sprintf("Could not read agent cluster %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AgentClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AgentClusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state for sysId
	var state AgentClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sysId from state
	data.SysId = state.SysId

	tflog.Debug(ctx, "Updating agent cluster", map[string]any{"sys_id": data.SysId.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Update the agent cluster
	_, err := r.client.Put(ctx, "/resources/agentcluster", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Agent Cluster",
			fmt.Sprintf("Could not update agent cluster %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back to get updated version
	err = r.readAgentCluster(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Agent Cluster",
			fmt.Sprintf("Could not read agent cluster %s after update: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AgentClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AgentClusterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting agent cluster", map[string]any{"sys_id": data.SysId.ValueString()})

	query := url.Values{}
	query.Set("agentclusterid", data.SysId.ValueString())

	_, err := r.client.Delete(ctx, "/resources/agentcluster", query)
	if err != nil {
		// Ignore 404 errors (already deleted)
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Agent Cluster",
			fmt.Sprintf("Could not delete agent cluster %s: %s", data.Name.ValueString(), err),
		)
		return
	}
}

func (r *AgentClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// readAgentCluster fetches the agent cluster from the API and updates the model.
func (r *AgentClusterResource) readAgentCluster(ctx context.Context, data *AgentClusterResourceModel) error {
	query := url.Values{}
	query.Set("agentclustername", data.Name.ValueString())

	respBody, err := r.client.Get(ctx, "/resources/agentcluster", query)
	if err != nil {
		return err
	}

	var apiModel AgentClusterAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		return fmt.Errorf("failed to parse agent cluster response: %w", err)
	}

	r.fromAPIModel(ctx, &apiModel, data)
	return nil
}

// toAPIModel converts the Terraform model to an API model.
func (r *AgentClusterResource) toAPIModel(ctx context.Context, data *AgentClusterResourceModel) *AgentClusterAPIModel {
	// Map user-friendly type to API discriminator
	apiType := data.Type.ValueString()
	if mapped, ok := agentClusterTypeToAPI[apiType]; ok {
		apiType = mapped
	}

	apiModel := &AgentClusterAPIModel{
		SysId:                 data.SysId.ValueString(),
		Name:                  data.Name.ValueString(),
		Type:                  apiType,
		Description:           data.Description.ValueString(),
		Distribution:          data.Distribution.ValueString(),
		NetworkAlias:          data.NetworkAlias.ValueString(),
		LimitType:             data.LimitType.ValueString(),
		AgentLimitType:        data.AgentLimitType.ValueString(),
		IgnoreInactiveAgents:  data.IgnoreInactiveAgents.ValueBool(),
		IgnoreSuspendedAgents: data.IgnoreSuspendedAgents.ValueBool(),
	}

	// Only set port if network alias is set
	if !data.NetworkAliasPort.IsNull() && !data.NetworkAliasPort.IsUnknown() {
		apiModel.NetworkAliasPort = data.NetworkAliasPort.ValueInt64()
	}

	// Only set limit amount if limit type is Limited
	if !data.LimitAmount.IsNull() && !data.LimitAmount.IsUnknown() {
		apiModel.LimitAmount = data.LimitAmount.ValueInt64()
	}

	// Only set agent limit amount if agent limit type is Limited
	if !data.AgentLimitAmount.IsNull() && !data.AgentLimitAmount.IsUnknown() {
		apiModel.AgentLimitAmount = data.AgentLimitAmount.ValueInt64()
	}

	// Handle opswise_groups
	if !data.OpswiseGroups.IsNull() && !data.OpswiseGroups.IsUnknown() {
		var groups []string
		data.OpswiseGroups.ElementsAs(ctx, &groups, false)
		apiModel.OpswiseGroups = groups
	}

	return apiModel
}

// fromAPIModel converts an API model to the Terraform model.
func (r *AgentClusterResource) fromAPIModel(ctx context.Context, apiModel *AgentClusterAPIModel, data *AgentClusterResourceModel) {
	// Identity fields - always set
	data.SysId = types.StringValue(apiModel.SysId)
	data.Name = types.StringValue(apiModel.Name)
	data.Version = types.Int64Value(apiModel.Version)

	// Basic - map API discriminator back to user-friendly type
	userType := apiModel.Type
	if mapped, ok := agentClusterTypeFromAPI[apiModel.Type]; ok {
		userType = mapped
	}
	data.Type = types.StringValue(userType)
	data.Description = StringValueOrNull(apiModel.Description)

	// Distribution
	data.Distribution = StringValueOrNull(apiModel.Distribution)

	// Network alias
	data.NetworkAlias = StringValueOrNull(apiModel.NetworkAlias)
	if apiModel.NetworkAliasPort > 0 {
		data.NetworkAliasPort = types.Int64Value(apiModel.NetworkAliasPort)
	} else {
		data.NetworkAliasPort = types.Int64Null()
	}

	// Task execution limits
	data.LimitType = StringValueOrNull(apiModel.LimitType)
	data.LimitAmount = types.Int64Value(apiModel.LimitAmount)

	// Per-agent limits
	data.AgentLimitType = StringValueOrNull(apiModel.AgentLimitType)
	data.AgentLimitAmount = types.Int64Value(apiModel.AgentLimitAmount)

	// Agent selection options
	data.IgnoreInactiveAgents = types.BoolValue(apiModel.IgnoreInactiveAgents)
	data.IgnoreSuspendedAgents = types.BoolValue(apiModel.IgnoreSuspendedAgents)

	// Business services
	if len(apiModel.OpswiseGroups) > 0 {
		groups, _ := types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		data.OpswiseGroups = groups
	} else {
		data.OpswiseGroups = types.ListNull(types.StringType)
	}
}
