package provider

import (
	"context"

	"github.com/alexprogrammr/appstore-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &appDataSource{}
	_ datasource.DataSourceWithConfigure = &appDataSource{}
)

type appDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	BundleID types.String `tfsdk:"bundle_id"`
	SKU      types.String `tfsdk:"sku"`
}

type appDataSource struct {
	client *appstore.Client
}

func NewAppDataSource() datasource.DataSource {
	return &appDataSource{}
}

func (d *appDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (d *appDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*appstore.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *appstore.Client, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	d.client = client
}

func (d *appDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches app information from the App Store Connect.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the app.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the app.",
				Computed:    true,
			},
			"bundle_id": schema.StringAttribute{
				Description: "Bundle identifier of the app.",
				Computed:    true,
			},
			"sku": schema.StringAttribute{
				Description: "Stock keeping unit of the app.",
				Computed:    true,
			},
		},
	}
}

func (d *appDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := appDataSourceModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := state.ID.ValueString()
	if appId == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'id' is required to fetch app information.",
		)
		return
	}

	app, err := d.client.GetApp(ctx, appId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read App",
			err.Error(),
		)
		return
	}

	state.Name = types.StringValue(app.Attr.Name)
	state.SKU = types.StringValue(app.Attr.SKU)
	state.BundleID = types.StringValue(app.Attr.BundleID)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
