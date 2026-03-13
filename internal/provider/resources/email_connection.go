package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/OptionMetrics/terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &EmailConnectionResource{}
	_ resource.ResourceWithImportState = &EmailConnectionResource{}
)

func NewEmailConnectionResource() resource.Resource {
	return &EmailConnectionResource{}
}

// EmailConnectionResource defines the resource implementation.
type EmailConnectionResource struct {
	client *client.Client
}

// EmailConnectionResourceModel describes the resource data model.
type EmailConnectionResourceModel struct {
	// Identity
	SysId   types.String `tfsdk:"sys_id"`
	Name    types.String `tfsdk:"name"`
	Version types.Int64  `tfsdk:"version"`

	// SMTP Settings
	SMTP         types.String `tfsdk:"smtp"`
	SMTPPort     types.Int64  `tfsdk:"smtp_port"`
	SMTPSSL      types.Bool   `tfsdk:"smtp_ssl"`
	SMTPStartTLS types.Bool   `tfsdk:"smtp_starttls"`

	// Sender
	EmailAddr types.String `tfsdk:"email_address"`

	// Authentication
	Authentication     types.Bool   `tfsdk:"authentication"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	DefaultUser        types.String `tfsdk:"default_user"`
	DefaultPwd         types.String `tfsdk:"default_password"`
	OAuthClient        types.String `tfsdk:"oauth_client"`

	// IMAP Settings (for reading emails)
	IMAP         types.String `tfsdk:"imap"`
	IMAPPort     types.Int64  `tfsdk:"imap_port"`
	IMAPSSL      types.Bool   `tfsdk:"imap_ssl"`
	IMAPStartTLS types.Bool   `tfsdk:"imap_starttls"`
	TrashFolder  types.String `tfsdk:"trash_folder"`

	// Other
	Description types.String `tfsdk:"description"`

	// Business services
	OpswiseGroups types.List `tfsdk:"opswise_groups"`
}

// EmailConnectionAPIModel represents the API request/response structure.
type EmailConnectionAPIModel struct {
	SysId   string `json:"sysId,omitempty"`
	Name    string `json:"name"`
	Version int64  `json:"version,omitempty"`
	Type    string `json:"type,omitempty"`

	// SMTP Settings
	SMTP         string `json:"smtp,omitempty"`
	SMTPPort     int64  `json:"smtpPort,omitempty"`
	SMTPSSL      bool   `json:"smtpSsl,omitempty"`
	SMTPStartTLS bool   `json:"smtpStarttls,omitempty"`

	// Sender
	EmailAddr string `json:"emailAddr,omitempty"`

	// Authentication
	Authentication     bool   `json:"authentication,omitempty"`
	AuthenticationType string `json:"authenticationType,omitempty"`
	DefaultUser        string `json:"defaultUser,omitempty"`
	DefaultPwd         string `json:"defaultPwd,omitempty"`
	OAuthClient        string `json:"oauthClient,omitempty"`

	// IMAP Settings
	IMAP         string `json:"imap,omitempty"`
	IMAPPort     int64  `json:"imapPort,omitempty"`
	IMAPSSL      bool   `json:"imapSsl,omitempty"`
	IMAPStartTLS bool   `json:"imapStarttls,omitempty"`
	TrashFolder  string `json:"trashFolder,omitempty"`

	// Other
	Description string `json:"description,omitempty"`

	OpswiseGroups []string `json:"opswiseGroups,omitempty"`
}

func (r *EmailConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_connection"
}

func (r *EmailConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a StoneBranch Email Connection. Email connections define SMTP and IMAP server settings for sending and receiving emails.",

		Attributes: map[string]schema.Attribute{
			// Identity
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the email connection (assigned by StoneBranch).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the email connection.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the email connection (for optimistic locking).",
				Computed:            true,
			},

			// SMTP Settings
			"smtp": schema.StringAttribute{
				MarkdownDescription: "SMTP server hostname or IP address.",
				Required:            true,
			},
			"smtp_port": schema.Int64Attribute{
				MarkdownDescription: "SMTP server port (e.g., 25, 465, 587).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"smtp_ssl": schema.BoolAttribute{
				MarkdownDescription: "Whether to use SSL/TLS for SMTP connection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"smtp_starttls": schema.BoolAttribute{
				MarkdownDescription: "Whether to use STARTTLS for SMTP connection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			// Sender
			"email_address": schema.StringAttribute{
				MarkdownDescription: "Default sender email address (From address).",
				Optional:            true,
			},

			// Authentication
			"authentication": schema.BoolAttribute{
				MarkdownDescription: "Whether authentication is required for SMTP.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"authentication_type": schema.StringAttribute{
				MarkdownDescription: "Type of authentication (e.g., 'Basic', 'OAuth').",
				Optional:            true,
				Computed:            true,
			},
			"default_user": schema.StringAttribute{
				MarkdownDescription: "Username for SMTP authentication.",
				Optional:            true,
			},
			"default_password": schema.StringAttribute{
				MarkdownDescription: "Password for SMTP authentication.",
				Optional:            true,
				Sensitive:           true,
			},
			"oauth_client": schema.StringAttribute{
				MarkdownDescription: "Name of the OAuth client to use for authentication.",
				Optional:            true,
			},

			// IMAP Settings
			"imap": schema.StringAttribute{
				MarkdownDescription: "IMAP server hostname or IP address (for reading emails).",
				Optional:            true,
			},
			"imap_port": schema.Int64Attribute{
				MarkdownDescription: "IMAP server port (e.g., 143, 993).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"imap_ssl": schema.BoolAttribute{
				MarkdownDescription: "Whether to use SSL/TLS for IMAP connection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"imap_starttls": schema.BoolAttribute{
				MarkdownDescription: "Whether to use STARTTLS for IMAP connection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"trash_folder": schema.StringAttribute{
				MarkdownDescription: "Name of the trash folder for IMAP operations.",
				Optional:            true,
			},

			// Other
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the email connection.",
				Optional:            true,
			},

			// Business services
			"opswise_groups": schema.ListAttribute{
				MarkdownDescription: "List of business service names this email connection belongs to.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *EmailConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EmailConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EmailConnectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating email connection", map[string]any{"name": data.Name.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Create the email connection
	_, err := r.client.Post(ctx, "/resources/emailconnection", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Email Connection",
			fmt.Sprintf("Could not create email connection %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back the created email connection to get sysId and other computed fields
	err = r.readEmailConnection(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Email Connection",
			fmt.Sprintf("Could not read email connection %s after creation: %s", data.Name.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Created email connection", map[string]any{"sys_id": data.SysId.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EmailConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EmailConnectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.readEmailConnection(ctx, &data)
	if err != nil {
		// Check if email connection was deleted outside of Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			tflog.Debug(ctx, "Email connection not found, removing from state", map[string]any{"name": data.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Email Connection",
			fmt.Sprintf("Could not read email connection %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EmailConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EmailConnectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state for sysId
	var state EmailConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sysId from state
	data.SysId = state.SysId

	tflog.Debug(ctx, "Updating email connection", map[string]any{"sys_id": data.SysId.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Update the email connection
	_, err := r.client.Put(ctx, "/resources/emailconnection", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Email Connection",
			fmt.Sprintf("Could not update email connection %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back to get updated version
	err = r.readEmailConnection(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Email Connection",
			fmt.Sprintf("Could not read email connection %s after update: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EmailConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EmailConnectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting email connection", map[string]any{"sys_id": data.SysId.ValueString()})

	query := url.Values{}
	query.Set("connectionid", data.SysId.ValueString())

	_, err := r.client.Delete(ctx, "/resources/emailconnection", query)
	if err != nil {
		// Ignore 404 errors (already deleted)
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Email Connection",
			fmt.Sprintf("Could not delete email connection %s: %s", data.Name.ValueString(), err),
		)
		return
	}
}

func (r *EmailConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// readEmailConnection fetches the email connection from the API and updates the model.
func (r *EmailConnectionResource) readEmailConnection(ctx context.Context, data *EmailConnectionResourceModel) error {
	query := url.Values{}
	query.Set("connectionname", data.Name.ValueString())

	respBody, err := r.client.Get(ctx, "/resources/emailconnection", query)
	if err != nil {
		return err
	}

	var apiModel EmailConnectionAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		return fmt.Errorf("failed to parse email connection response: %w", err)
	}

	r.fromAPIModel(ctx, &apiModel, data)
	return nil
}

// toAPIModel converts the Terraform model to an API model.
func (r *EmailConnectionResource) toAPIModel(ctx context.Context, data *EmailConnectionResourceModel) *EmailConnectionAPIModel {
	model := &EmailConnectionAPIModel{
		SysId:       data.SysId.ValueString(),
		Name:        data.Name.ValueString(),
		Type:        "Outgoing",
		SMTP:        data.SMTP.ValueString(),
		EmailAddr:   data.EmailAddr.ValueString(),
		DefaultUser: data.DefaultUser.ValueString(),
		DefaultPwd:  data.DefaultPwd.ValueString(),
		OAuthClient: data.OAuthClient.ValueString(),
		IMAP:        data.IMAP.ValueString(),
		TrashFolder: data.TrashFolder.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Handle authentication_type
	if !data.AuthenticationType.IsNull() && !data.AuthenticationType.IsUnknown() {
		model.AuthenticationType = data.AuthenticationType.ValueString()
	}

	// Handle SMTP port
	if !data.SMTPPort.IsNull() && !data.SMTPPort.IsUnknown() {
		model.SMTPPort = data.SMTPPort.ValueInt64()
	}

	// Handle SMTP SSL
	if !data.SMTPSSL.IsNull() && !data.SMTPSSL.IsUnknown() {
		model.SMTPSSL = data.SMTPSSL.ValueBool()
	}

	// Handle SMTP STARTTLS
	if !data.SMTPStartTLS.IsNull() && !data.SMTPStartTLS.IsUnknown() {
		model.SMTPStartTLS = data.SMTPStartTLS.ValueBool()
	}

	// Handle authentication
	if !data.Authentication.IsNull() && !data.Authentication.IsUnknown() {
		model.Authentication = data.Authentication.ValueBool()
	}

	// Handle IMAP port
	if !data.IMAPPort.IsNull() && !data.IMAPPort.IsUnknown() {
		model.IMAPPort = data.IMAPPort.ValueInt64()
	}

	// Handle IMAP SSL
	if !data.IMAPSSL.IsNull() && !data.IMAPSSL.IsUnknown() {
		model.IMAPSSL = data.IMAPSSL.ValueBool()
	}

	// Handle IMAP STARTTLS
	if !data.IMAPStartTLS.IsNull() && !data.IMAPStartTLS.IsUnknown() {
		model.IMAPStartTLS = data.IMAPStartTLS.ValueBool()
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
func (r *EmailConnectionResource) fromAPIModel(ctx context.Context, apiModel *EmailConnectionAPIModel, data *EmailConnectionResourceModel) {
	// Identity fields - always set
	data.SysId = types.StringValue(apiModel.SysId)
	data.Name = types.StringValue(apiModel.Name)
	data.Version = types.Int64Value(apiModel.Version)

	// SMTP Settings
	data.SMTP = StringValueOrNull(apiModel.SMTP)
	data.SMTPPort = types.Int64Value(apiModel.SMTPPort)
	data.SMTPSSL = types.BoolValue(apiModel.SMTPSSL)
	data.SMTPStartTLS = types.BoolValue(apiModel.SMTPStartTLS)

	// Sender
	data.EmailAddr = StringValueOrNull(apiModel.EmailAddr)

	// Authentication
	data.Authentication = types.BoolValue(apiModel.Authentication)
	data.AuthenticationType = StringValueOrNull(apiModel.AuthenticationType)
	data.DefaultUser = StringValueOrNull(apiModel.DefaultUser)
	// Note: Password is typically not returned by the API for security
	// Only set if we have a value (to preserve state)
	if apiModel.DefaultPwd != "" {
		data.DefaultPwd = StringValueOrNull(apiModel.DefaultPwd)
	}
	data.OAuthClient = StringValueOrNull(apiModel.OAuthClient)

	// IMAP Settings
	data.IMAP = StringValueOrNull(apiModel.IMAP)
	data.IMAPPort = types.Int64Value(apiModel.IMAPPort)
	data.IMAPSSL = types.BoolValue(apiModel.IMAPSSL)
	data.IMAPStartTLS = types.BoolValue(apiModel.IMAPStartTLS)
	data.TrashFolder = StringValueOrNull(apiModel.TrashFolder)

	// Other
	data.Description = StringValueOrNull(apiModel.Description)

	// Handle opswise_groups
	if len(apiModel.OpswiseGroups) > 0 {
		groups, _ := types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		data.OpswiseGroups = groups
	}
}
