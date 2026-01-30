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

	"terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &BusinessServiceResource{}
	_ resource.ResourceWithImportState = &BusinessServiceResource{}
)

func NewBusinessServiceResource() resource.Resource {
	return &BusinessServiceResource{}
}

// BusinessServiceResource defines the resource implementation.
type BusinessServiceResource struct {
	client *client.Client
}

// BusinessServiceResourceModel describes the resource data model.
type BusinessServiceResourceModel struct {
	// Identity
	SysId   types.String `tfsdk:"sys_id"`
	Name    types.String `tfsdk:"name"`
	Version types.Int64  `tfsdk:"version"`

	// Content
	Description types.String `tfsdk:"description"`
}

// BusinessServiceAPIModel represents the API request/response structure.
type BusinessServiceAPIModel struct {
	SysId       string `json:"sysId,omitempty"`
	Name        string `json:"name"`
	Version     int64  `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r *BusinessServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_business_service"
}

func (r *BusinessServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a StoneBranch Business Service. Business Services are used to group and organize resources such as tasks, triggers, and variables.",

		Attributes: map[string]schema.Attribute{
			// Identity
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the business service (assigned by StoneBranch).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the business service.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the business service (for optimistic locking).",
				Computed:            true,
			},

			// Content
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the business service.",
				Optional:            true,
			},
		},
	}
}

func (r *BusinessServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BusinessServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BusinessServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating business service", map[string]any{"name": data.Name.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(&data)

	// Create the business service
	_, err := r.client.Post(ctx, "/resources/businessservice", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Business Service",
			fmt.Sprintf("Could not create business service %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back the created business service to get sysId and other computed fields
	err = r.readBusinessService(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Business Service",
			fmt.Sprintf("Could not read business service %s after creation: %s", data.Name.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Created business service", map[string]any{"sys_id": data.SysId.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BusinessServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BusinessServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.readBusinessService(ctx, &data)
	if err != nil {
		// Check if business service was deleted outside of Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			tflog.Debug(ctx, "Business service not found, removing from state", map[string]any{"name": data.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Business Service",
			fmt.Sprintf("Could not read business service %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BusinessServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BusinessServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state for sysId
	var state BusinessServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sysId from state
	data.SysId = state.SysId

	tflog.Debug(ctx, "Updating business service", map[string]any{"sys_id": data.SysId.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(&data)

	// Update the business service
	_, err := r.client.Put(ctx, "/resources/businessservice", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Business Service",
			fmt.Sprintf("Could not update business service %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back to get updated version
	err = r.readBusinessService(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Business Service",
			fmt.Sprintf("Could not read business service %s after update: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BusinessServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BusinessServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting business service", map[string]any{"sys_id": data.SysId.ValueString()})

	query := url.Values{}
	query.Set("busserviceid", data.SysId.ValueString())

	_, err := r.client.Delete(ctx, "/resources/businessservice", query)
	if err != nil {
		// Ignore 404 errors (already deleted)
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Business Service",
			fmt.Sprintf("Could not delete business service %s: %s", data.Name.ValueString(), err),
		)
		return
	}
}

func (r *BusinessServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// readBusinessService fetches the business service from the API and updates the model.
func (r *BusinessServiceResource) readBusinessService(ctx context.Context, data *BusinessServiceResourceModel) error {
	query := url.Values{}
	query.Set("busservicename", data.Name.ValueString())

	respBody, err := r.client.Get(ctx, "/resources/businessservice", query)
	if err != nil {
		return err
	}

	var apiModel BusinessServiceAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		return fmt.Errorf("failed to parse business service response: %w", err)
	}

	r.fromAPIModel(&apiModel, data)
	return nil
}

// toAPIModel converts the Terraform model to an API model.
func (r *BusinessServiceResource) toAPIModel(data *BusinessServiceResourceModel) *BusinessServiceAPIModel {
	return &BusinessServiceAPIModel{
		SysId:       data.SysId.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}
}

// fromAPIModel converts an API model to the Terraform model.
func (r *BusinessServiceResource) fromAPIModel(apiModel *BusinessServiceAPIModel, data *BusinessServiceResourceModel) {
	// Identity fields - always set
	data.SysId = types.StringValue(apiModel.SysId)
	data.Name = types.StringValue(apiModel.Name)
	data.Version = types.Int64Value(apiModel.Version)

	// Content
	data.Description = StringValueOrNull(apiModel.Description)
}
