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
	_ resource.Resource                = &VariableResource{}
	_ resource.ResourceWithImportState = &VariableResource{}
)

func NewVariableResource() resource.Resource {
	return &VariableResource{}
}

// VariableResource defines the resource implementation.
type VariableResource struct {
	client *client.Client
}

// VariableResourceModel describes the resource data model.
type VariableResourceModel struct {
	// Identity
	SysId   types.String `tfsdk:"sys_id"`
	Name    types.String `tfsdk:"name"`
	Version types.Int64  `tfsdk:"version"`

	// Content
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`

	// Business services
	OpswiseGroups types.List `tfsdk:"opswise_groups"`
}

// VariableAPIModel represents the API request/response structure.
type VariableAPIModel struct {
	SysId         string   `json:"sysId,omitempty"`
	Name          string   `json:"name"`
	Version       int64    `json:"version,omitempty"`
	Value         string   `json:"value,omitempty"`
	Description   string   `json:"description,omitempty"`
	OpswiseGroups []string `json:"opswiseGroups,omitempty"`
}

func (r *VariableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (r *VariableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a StoneBranch Global Variable. Variables store reusable values that can be referenced by tasks and triggers.",

		Attributes: map[string]schema.Attribute{
			// Identity
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the variable (assigned by StoneBranch).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the variable. Must begin with a letter. Only alphanumerics and underscores allowed (no hyphens or spaces). Names are not case-sensitive. Do not use the prefix `ops_` (reserved for built-in variables).",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the variable (for optimistic locking).",
				Computed:            true,
			},

			// Content
			"value": schema.StringAttribute{
				MarkdownDescription: "Value of the variable.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the variable.",
				Optional:            true,
			},

			// Business services
			"opswise_groups": schema.ListAttribute{
				MarkdownDescription: "List of business service names this variable belongs to.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *VariableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VariableResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating variable", map[string]any{"name": data.Name.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Create the variable
	_, err := r.client.Post(ctx, "/resources/variable", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Variable",
			fmt.Sprintf("Could not create variable %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back the created variable to get sysId and other computed fields
	err = r.readVariable(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Variable",
			fmt.Sprintf("Could not read variable %s after creation: %s", data.Name.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Created variable", map[string]any{"sys_id": data.SysId.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VariableResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.readVariable(ctx, &data)
	if err != nil {
		// Check if variable was deleted outside of Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			tflog.Debug(ctx, "Variable not found, removing from state", map[string]any{"name": data.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Variable",
			fmt.Sprintf("Could not read variable %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VariableResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state for sysId
	var state VariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sysId from state
	data.SysId = state.SysId

	tflog.Debug(ctx, "Updating variable", map[string]any{"sys_id": data.SysId.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Update the variable
	_, err := r.client.Put(ctx, "/resources/variable", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Variable",
			fmt.Sprintf("Could not update variable %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back to get updated version
	err = r.readVariable(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Variable",
			fmt.Sprintf("Could not read variable %s after update: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VariableResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting variable", map[string]any{"sys_id": data.SysId.ValueString()})

	query := url.Values{}
	query.Set("variableid", data.SysId.ValueString())

	_, err := r.client.Delete(ctx, "/resources/variable", query)
	if err != nil {
		// Ignore 404 errors (already deleted)
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Variable",
			fmt.Sprintf("Could not delete variable %s: %s", data.Name.ValueString(), err),
		)
		return
	}
}

func (r *VariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// readVariable fetches the variable from the API and updates the model.
func (r *VariableResource) readVariable(ctx context.Context, data *VariableResourceModel) error {
	query := url.Values{}
	query.Set("variablename", data.Name.ValueString())

	respBody, err := r.client.Get(ctx, "/resources/variable", query)
	if err != nil {
		return err
	}

	var apiModel VariableAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		return fmt.Errorf("failed to parse variable response: %w", err)
	}

	r.fromAPIModel(ctx, &apiModel, data)
	return nil
}

// toAPIModel converts the Terraform model to an API model.
func (r *VariableResource) toAPIModel(ctx context.Context, data *VariableResourceModel) *VariableAPIModel {
	model := &VariableAPIModel{
		SysId:       data.SysId.ValueString(),
		Name:        data.Name.ValueString(),
		Value:       data.Value.ValueString(),
		Description: data.Description.ValueString(),
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
func (r *VariableResource) fromAPIModel(ctx context.Context, apiModel *VariableAPIModel, data *VariableResourceModel) {
	// Identity fields - always set
	data.SysId = types.StringValue(apiModel.SysId)
	data.Name = types.StringValue(apiModel.Name)
	data.Version = types.Int64Value(apiModel.Version)

	// Content
	data.Value = StringValueOrNull(apiModel.Value)
	data.Description = StringValueOrNull(apiModel.Description)

	// Handle opswise_groups
	if len(apiModel.OpswiseGroups) > 0 {
		groups, _ := types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		data.OpswiseGroups = groups
	}
}
