package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-stonebranch/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DatabaseConnectionResource{}
	_ resource.ResourceWithImportState = &DatabaseConnectionResource{}
)

func NewDatabaseConnectionResource() resource.Resource {
	return &DatabaseConnectionResource{}
}

// DatabaseConnectionResource defines the resource implementation.
type DatabaseConnectionResource struct {
	client *client.Client
}

// DatabaseConnectionResourceModel describes the resource data model.
type DatabaseConnectionResourceModel struct {
	// Identity
	SysId   types.String `tfsdk:"sys_id"`
	Name    types.String `tfsdk:"name"`
	Version types.Int64  `tfsdk:"version"`

	// Database connection details
	DBType        types.String `tfsdk:"db_type"`
	DBUrl         types.String `tfsdk:"db_url"`
	DBDriver      types.String `tfsdk:"db_driver"`
	DBMaxRows     types.Int64  `tfsdk:"db_max_rows"`
	DBDescription types.String `tfsdk:"description"`

	// Credentials
	Credentials types.String `tfsdk:"credentials"`

	// Business services
	OpswiseGroups types.List `tfsdk:"opswise_groups"`
}

// DatabaseConnectionAPIModel represents the API request/response structure.
type DatabaseConnectionAPIModel struct {
	SysId         string   `json:"sysId,omitempty"`
	Name          string   `json:"name"`
	Version       int64    `json:"version,omitempty"`
	DBType        string   `json:"dbType,omitempty"`
	DBUrl         string   `json:"dbUrl,omitempty"`
	DBDriver      string   `json:"dbDriver,omitempty"`
	DBMaxRows     int64    `json:"dbMaxRows,omitempty"`
	DBDescription string   `json:"dbDescription,omitempty"`
	Credentials   string   `json:"credentials,omitempty"`
	OpswiseGroups []string `json:"opswiseGroups,omitempty"`
}

func (r *DatabaseConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_connection"
}

func (r *DatabaseConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a StoneBranch Database Connection. Database connections define how to connect to databases for SQL tasks.",

		Attributes: map[string]schema.Attribute{
			// Identity
			"sys_id": schema.StringAttribute{
				MarkdownDescription: "System ID of the database connection (assigned by StoneBranch).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the database connection.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version number of the database connection (for optimistic locking).",
				Computed:            true,
			},

			// Database connection details
			"db_type": schema.StringAttribute{
				MarkdownDescription: "Type of database (e.g., 'MySQL', 'Oracle', 'PostgreSQL', 'SQL Server', 'Other').",
				Optional:            true,
				Computed:            true,
			},
			"db_url": schema.StringAttribute{
				MarkdownDescription: "JDBC URL for the database connection (e.g., 'jdbc:mysql://localhost:3306/mydb').",
				Required:            true,
			},
			"db_driver": schema.StringAttribute{
				MarkdownDescription: "Fully qualified Java class name of the JDBC driver (e.g., 'com.mysql.cj.jdbc.Driver').",
				Required:            true,
			},
			"db_max_rows": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of rows to return from queries. Use 0 for unlimited.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the database connection.",
				Optional:            true,
			},

			// Credentials
			"credentials": schema.StringAttribute{
				MarkdownDescription: "Name of the credential to use for database authentication.",
				Optional:            true,
			},

			// Business services
			"opswise_groups": schema.ListAttribute{
				MarkdownDescription: "List of business service names this database connection belongs to.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *DatabaseConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseConnectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating database connection", map[string]any{"name": data.Name.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Create the database connection
	_, err := r.client.Post(ctx, "/resources/databaseconnection", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Database Connection",
			fmt.Sprintf("Could not create database connection %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back the created database connection to get sysId and other computed fields
	err = r.readDatabaseConnection(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Database Connection",
			fmt.Sprintf("Could not read database connection %s after creation: %s", data.Name.ValueString(), err),
		)
		return
	}

	tflog.Debug(ctx, "Created database connection", map[string]any{"sys_id": data.SysId.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseConnectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.readDatabaseConnection(ctx, &data)
	if err != nil {
		// Check if database connection was deleted outside of Terraform
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			tflog.Debug(ctx, "Database connection not found, removing from state", map[string]any{"name": data.Name.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Database Connection",
			fmt.Sprintf("Could not read database connection %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseConnectionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state for sysId
	var state DatabaseConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sysId from state
	data.SysId = state.SysId

	tflog.Debug(ctx, "Updating database connection", map[string]any{"sys_id": data.SysId.ValueString()})

	// Build API model
	apiModel := r.toAPIModel(ctx, &data)

	// Update the database connection
	_, err := r.client.Put(ctx, "/resources/databaseconnection", apiModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Database Connection",
			fmt.Sprintf("Could not update database connection %s: %s", data.Name.ValueString(), err),
		)
		return
	}

	// Read back to get updated version
	err = r.readDatabaseConnection(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Database Connection",
			fmt.Sprintf("Could not read database connection %s after update: %s", data.Name.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseConnectionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting database connection", map[string]any{"sys_id": data.SysId.ValueString()})

	query := url.Values{}
	query.Set("connectionid", data.SysId.ValueString())

	_, err := r.client.Delete(ctx, "/resources/databaseconnection", query)
	if err != nil {
		// Ignore 404 errors (already deleted)
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Database Connection",
			fmt.Sprintf("Could not delete database connection %s: %s", data.Name.ValueString(), err),
		)
		return
	}
}

func (r *DatabaseConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// readDatabaseConnection fetches the database connection from the API and updates the model.
func (r *DatabaseConnectionResource) readDatabaseConnection(ctx context.Context, data *DatabaseConnectionResourceModel) error {
	query := url.Values{}
	query.Set("connectionname", data.Name.ValueString())

	respBody, err := r.client.Get(ctx, "/resources/databaseconnection", query)
	if err != nil {
		return err
	}

	var apiModel DatabaseConnectionAPIModel
	if err := json.Unmarshal(respBody, &apiModel); err != nil {
		return fmt.Errorf("failed to parse database connection response: %w", err)
	}

	r.fromAPIModel(ctx, &apiModel, data)
	return nil
}

// toAPIModel converts the Terraform model to an API model.
func (r *DatabaseConnectionResource) toAPIModel(ctx context.Context, data *DatabaseConnectionResourceModel) *DatabaseConnectionAPIModel {
	model := &DatabaseConnectionAPIModel{
		SysId:         data.SysId.ValueString(),
		Name:          data.Name.ValueString(),
		DBType:        data.DBType.ValueString(),
		DBUrl:         data.DBUrl.ValueString(),
		DBDriver:      data.DBDriver.ValueString(),
		DBDescription: data.DBDescription.ValueString(),
		Credentials:   data.Credentials.ValueString(),
	}

	// Handle db_max_rows - only set if specified
	if !data.DBMaxRows.IsNull() && !data.DBMaxRows.IsUnknown() {
		model.DBMaxRows = data.DBMaxRows.ValueInt64()
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
func (r *DatabaseConnectionResource) fromAPIModel(ctx context.Context, apiModel *DatabaseConnectionAPIModel, data *DatabaseConnectionResourceModel) {
	// Identity fields - always set
	data.SysId = types.StringValue(apiModel.SysId)
	data.Name = types.StringValue(apiModel.Name)
	data.Version = types.Int64Value(apiModel.Version)

	// Database connection details
	data.DBType = StringValueOrNull(apiModel.DBType)
	data.DBUrl = StringValueOrNull(apiModel.DBUrl)
	data.DBDriver = StringValueOrNull(apiModel.DBDriver)
	data.DBMaxRows = types.Int64Value(apiModel.DBMaxRows)
	data.DBDescription = StringValueOrNull(apiModel.DBDescription)

	// Credentials
	data.Credentials = StringValueOrNull(apiModel.Credentials)

	// Handle opswise_groups
	if len(apiModel.OpswiseGroups) > 0 {
		groups, _ := types.ListValueFrom(ctx, types.StringType, apiModel.OpswiseGroups)
		data.OpswiseGroups = groups
	}
}
