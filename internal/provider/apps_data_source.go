package provider

import (
	"context"
	"fmt"

	"github.com/alexprogrammr/appstore-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &appsDataSource{}
	_ datasource.DataSourceWithConfigure = &appsDataSource{}
)

type appsDataSourceModel struct {
	Apps []appsModel `tfsdk:"apps"`
}

type appsModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	BundleID types.String `tfsdk:"bundle_id"`
	SKU      types.String `tfsdk:"sku"`
}

type appsDataSource struct {
	client *appstore.Client
}

func NewAppsDataSource() datasource.DataSource {
	return &appsDataSource{}
}

func (d *appsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apps"
}

func (d *appsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*appstore.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *appstore.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *appsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of apps from the App Store Connect.",
		Attributes: map[string]schema.Attribute{
			"apps": schema.ListNestedAttribute{
				Description: "List of apps",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Identifier of the app.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *appsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := appsDataSourceModel{}

	apps, err := d.client.ListApps(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Apps",
			err.Error(),
		)
		return
	}

	for _, app := range apps {
		appState := appsModel{
			ID:       types.StringValue(app.ID),
			Name:     types.StringValue(app.Attr.Name),
			SKU:      types.StringValue(app.Attr.SKU),
			BundleID: types.StringValue(app.Attr.BundleID),
		}

		state.Apps = append(state.Apps, appState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
