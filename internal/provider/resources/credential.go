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
	_ resource.Resource                = &CredentialResource{}
	_ resource.ResourceWithImportState = &CredentialResource{}
)

func NewCredentialResource() resource.Resource {
	return &CredentialResource{}
}

// CredentialResource defines the resource implementation.
type CredentialResource struct {
	client *client.Client
}

// CredentialResourceModel describes the resource data model.
type CredentialResourceModel struct {
	// Identity
	SysId   types.String `tfsdk:"sys_id"`
	Name    types.String `tfsdk:"name"`
	Version types.Int64  `tfsdk:"version"`

	// Basic info
	Description types.String `tfsdk:"description"`

	// Runtime credentials
	RuntimeUser        types.String `tfsdk:"runtime_user"`
	RuntimePassword    types.String `tfsdk:"runtime_password"`
	RuntimePassphrase  types.String `tfsdk:"runtime_passphrase"`
	RuntimeToken       types.String `tfsdk:"runtime_token"`
	RuntimeKeyLocation types.String `tfsdk:"runtime_key_location"`

	// Provider (for external credential providers)
	Provider types.String `tfsdk:"provider_name"`

	// Business services
	OpswiseGroups types.List `tfsdk:"opswise_groups"`
}

// CredentialAPIModel represents the API request/response structure.
type CredentialAPIModel struct {
	SysId              string   `json:"sysId,omitempty"`
	Name               string   `json:"name"`
	Version            int64    `json:"version,omitempty"`
	Description        string   `json:"description,omitempty"`
	RuntimeUser        string   `json:"runtimeUser,omitempty"`
	RuntimePassword    string   `json:"runtimePassword,omitempty"`
	RuntimePassPhrase  string   `json:"runtimePassPhrase,omitempty"`
	RuntimeToken       string   `json:"runtimeToken,omitempty"`
	RuntimeKeyLocation string   `json:"runtimeKeyLocation,omitempty"`
	Provider           string   `json:"provider,omitempty"`
	OpswiseGroups      []string `json:"opswiseGroups,omitempty"`
}

func (r *CredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential"
}

func (r *CredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a StoneBranch Credential. Credentials store authentication information for use by tasks.",

		Attributes: map[string]schema.Attribute{
			// Identity
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the credential (assigned by StoneBranch).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the credential.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the credential (for optimistic locking).",
				Computed:            true,
			},

			// Basic info
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the credential.",
				Optional:            true,
			},

			// Runtime credentials
			"runtime_user": schema.StringAttribute{
				MarkdownDescription: "Username for authentication.",
				Optional:            true,
			},
			"runtime_password": schema.StringAttribute{
				MarkdownDescription: "Password for authentication. Note: This value is write-only and will not be returned by the API.",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"runtime_passphrase": schema.StringAttribute{
				MarkdownDescription: "Passphrase for key-based authentication. Note: This value is write-only and will not be returned by the API.",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"runtime_token": schema.StringAttribute{
				MarkdownDescription: "Token for token-based authentication. Note: This value is write-only and will not be returned by the API.",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"runtime_key_location": schema.StringAttribute{
				MarkdownDescription: "Path to the private key file for key-based authentication.",
				Optional:            true,
			},

			// Provider
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "Name of the external credential provider (for vault integrations). Defaults to 'Universal Controller'.",
				Optional:            true,
				Computed:            true,
			},

			// Business services
			"opswise_groups": schema.ListAttribute{
				MarkdownDescription: "List of business service names this credential belongs to.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *CredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CredentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating credential", map[string]any{"name": data.Name.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Create the credential
	_, err := r.client.Post(ctx, "/resources/credential", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Credential",
			fmt.Sprintf("Could not create credential %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back the created credential to get sysId and other computed fields
	err = r.readCredential(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Credential",
			fmt.Sprintf("Could not read credential %s after creation: %s", data.Name.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Created credential", map[string]any{"sys_id": data.SysId.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CredentialResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.readCredential(ctx, &data)
	if err != nil {
		// Check if credential was deleted outside of Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			tflog.Debug(ctx, "Credential not found, removing from state", map[string]any{"name": data.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Credential",
			fmt.Sprintf("Could not read credential %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CredentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state for sysId
	var state CredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sysId from state
	data.SysId = state.SysId

	tflog.Debug(ctx, "Updating credential", map[string]any{"sys_id": data.SysId.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Update the credential
	_, err := r.client.Put(ctx, "/resources/credential", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Credential",
			fmt.Sprintf("Could not update credential %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back to get updated version
	err = r.readCredential(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Credential",
			fmt.Sprintf("Could not read credential %s after update: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CredentialResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting credential", map[string]any{"sys_id": data.SysId.ValueString()})

	query := url.Values{}
	query.Set("credentialid", data.SysId.ValueString())

	_, err := r.client.Delete(ctx, "/resources/credential", query)
	if err != nil {
		// Ignore 404 errors (already deleted)
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Credential",
			fmt.Sprintf("Could not delete credential %s: %s", data.Name.ValueString(), err),
		)
		return
	}
}

func (r *CredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// readCredential fetches the credential from the API and updates the model.
func (r *CredentialResource) readCredential(ctx context.Context, data *CredentialResourceModel) error {
	query := url.Values{}
	query.Set("credentialname", data.Name.ValueString())

	respBody, err := r.client.Get(ctx, "/resources/credential", query)
	if err != nil {
		return err
	}

	var apiModel CredentialAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		return fmt.Errorf("failed to parse credential response: %w", err)
	}

	r.fromAPIModel(ctx, &apiModel, data)
	return nil
}

// toAPIModel converts the Terraform model to an API model.
func (r *CredentialResource) toAPIModel(ctx context.Context, data *CredentialResourceModel) *CredentialAPIModel {
	model := &CredentialAPIModel{
		SysId:              data.SysId.ValueString(),
		Name:               data.Name.ValueString(),
		Description:        data.Description.ValueString(),
		RuntimeUser:        data.RuntimeUser.ValueString(),
		RuntimePassword:    data.RuntimePassword.ValueString(),
		RuntimePassPhrase:  data.RuntimePassphrase.ValueString(),
		RuntimeToken:       data.RuntimeToken.ValueString(),
		RuntimeKeyLocation: data.RuntimeKeyLocation.ValueString(),
		Provider:           data.Provider.ValueString(),
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
func (r *CredentialResource) fromAPIModel(ctx context.Context, apiModel *CredentialAPIModel, data *CredentialResourceModel) {
	// Identity fields - always set
	data.SysId = types.StringValue(apiModel.SysId)
	data.Name = types.StringValue(apiModel.Name)
	data.Version = types.Int64Value(apiModel.Version)

	// Basic info
	data.Description = StringValueOrNull(apiModel.Description)

	// Runtime credentials - Note: passwords/tokens are NOT returned by the API for security
	// We preserve the values from the plan/state since the API won't return them
	data.RuntimeUser = StringValueOrNull(apiModel.RuntimeUser)
	// RuntimePassword, RuntimePassphrase, RuntimeToken are NOT read back from API
	data.RuntimeKeyLocation = StringValueOrNull(apiModel.RuntimeKeyLocation)

	// Provider
	data.Provider = StringValueOrNull(apiModel.Provider)

	// Handle opswise_groups
	if len(apiModel.OpswiseGroups) > 0 {
		groups, _ := types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		data.OpswiseGroups = groups
	}
}
