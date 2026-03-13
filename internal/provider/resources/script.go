package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/OptionMetrics/terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ScriptResource{}
	_ resource.ResourceWithImportState = &ScriptResource{}
)

func NewScriptResource() resource.Resource {
	return &ScriptResource{}
}

// ScriptResource defines the resource implementation.
type ScriptResource struct {
	client *client.Client
}

// ScriptResourceModel describes the resource data model.
type ScriptResourceModel struct {
	// Identity
	SysId   types.String `tfsdk:"sys_id"`
	Name    types.String `tfsdk:"name"`
	Version types.Int64  `tfsdk:"version"`

	// Script content
	ScriptType types.String `tfsdk:"script_type"`
	Content    types.String `tfsdk:"content"`

	// Optional fields
	Description      types.String `tfsdk:"description"`
	ResolveVariables types.Bool   `tfsdk:"resolve_variables"`

	// Business services
	OpswiseGroups types.List `tfsdk:"opswise_groups"`
}

// ScriptAPIModel represents the API request/response structure.
type ScriptAPIModel struct {
	SysId            string   `json:"sysId,omitempty"`
	ScriptName       string   `json:"scriptName"`
	Version          int64    `json:"version,omitempty"`
	ScriptType       string   `json:"scriptType,omitempty"`
	Description      string   `json:"description,omitempty"`
	Content          string   `json:"content,omitempty"`
	ResolveVariables bool     `json:"resolveVariables,omitempty"`
	OpswiseGroups    []string `json:"opswiseGroups,omitempty"`
}

func (r *ScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script"
}

func (r *ScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a StoneBranch Script resource. Scripts can be referenced by Unix/Linux tasks using the `script` attribute.",

		Attributes: map[string]schema.Attribute{
			// Identity
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the script (assigned by StoneBranch).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the script. This name is used to reference the script from tasks.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the script (for optimistic locking).",
				Computed:            true,
			},

			// Script content
			"script_type": schema.StringAttribute{
				MarkdownDescription: "Type of script (e.g., 'Unix', 'Windows').",
				Optional:            true,
				Computed:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The actual script content/code.",
				Required:            true,
			},

			// Optional fields
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the script.",
				Optional:            true,
			},
			"resolve_variables": schema.BoolAttribute{
				MarkdownDescription: "Whether to resolve variables in the script content at runtime.",
				Optional:            true,
				Computed:            true,
			},

			// Business services
			"opswise_groups": schema.ListAttribute{
				MarkdownDescription: "List of business service names this script belongs to.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *ScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ScriptResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating script", map[string]any{"name": data.Name.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Create the script
	_, err := r.client.Post(ctx, "/resources/script", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Script",
			fmt.Sprintf("Could not create script %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back the created script to get sysId and other computed fields
	err = r.readScript(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Script",
			fmt.Sprintf("Could not read script %s after creation: %s", data.Name.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Created script", map[string]any{"sys_id": data.SysId.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ScriptResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.readScript(ctx, &data)
	if err != nil {
		// Check if script was deleted outside of Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			tflog.Debug(ctx, "Script not found, removing from state", map[string]any{"name": data.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Script",
			fmt.Sprintf("Could not read script %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ScriptResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state for sysId
	var state ScriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sysId from state
	data.SysId = state.SysId

	tflog.Debug(ctx, "Updating script", map[string]any{"sys_id": data.SysId.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Update the script
	_, err := r.client.Put(ctx, "/resources/script", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Script",
			fmt.Sprintf("Could not update script %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back to get updated version
	err = r.readScript(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Script",
			fmt.Sprintf("Could not read script %s after update: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ScriptResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting script", map[string]any{"sys_id": data.SysId.ValueString()})

	query := url.Values{}
	query.Set("scriptid", data.SysId.ValueString())

	_, err := r.client.Delete(ctx, "/resources/script", query)
	if err != nil {
		// Ignore 404 errors (already deleted)
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Script",
			fmt.Sprintf("Could not delete script %s: %s", data.Name.ValueString(), err),
		)
		return
	}
}

func (r *ScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// readScript fetches the script from the API and updates the model.
func (r *ScriptResource) readScript(ctx context.Context, data *ScriptResourceModel) error {
	query := url.Values{}
	query.Set("scriptname", data.Name.ValueString())

	respBody, err := r.client.Get(ctx, "/resources/script", query)
	if err != nil {
		return err
	}

	var apiModel ScriptAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		return fmt.Errorf("failed to parse script response: %w", err)
	}

	r.fromAPIModel(ctx, &apiModel, data)
	return nil
}

// toAPIModel converts the Terraform model to an API model.
func (r *ScriptResource) toAPIModel(ctx context.Context, data *ScriptResourceModel) *ScriptAPIModel {
	model := &ScriptAPIModel{
		SysId:            data.SysId.ValueString(),
		ScriptName:       data.Name.ValueString(),
		ScriptType:       data.ScriptType.ValueString(),
		Description:      data.Description.ValueString(),
		Content:          data.Content.ValueString(),
		ResolveVariables: data.ResolveVariables.ValueBool(),
	}

	// Handle opswise_groups list
	if !data.OpswiseGroups.IsNull() && !data.OpswiseGroups.IsUnknown() {
		var groups []string
		data.OpswiseGroups.ElementsAs(ctx, &groups, false)
		model.OpswiseGroups = groups
	}

	return model
}

// fromAPIModel converts an API model to the Terraform model.
func (r *ScriptResource) fromAPIModel(ctx context.Context, apiModel *ScriptAPIModel, data *ScriptResourceModel) {
	// Identity fields - always set
	data.SysId = types.StringValue(apiModel.SysId)
	data.Name = types.StringValue(apiModel.ScriptName)
	data.Version = types.Int64Value(apiModel.Version)

	// Script content
	data.ScriptType = types.StringValue(apiModel.ScriptType)
	data.Content = types.StringValue(apiModel.Content)

	// Optional fields
	data.Description = StringValueOrNull(apiModel.Description)
	data.ResolveVariables = types.BoolValue(apiModel.ResolveVariables)

	// Handle opswise_groups
	if len(apiModel.OpswiseGroups) > 0 {
		groups, _ := types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		data.OpswiseGroups = groups
	}
}
