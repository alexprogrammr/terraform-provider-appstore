package provider

import (
	"context"

	"github.com/alexprogrammr/appstore-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &gameCenterDataSource{}
	_ datasource.DataSourceWithConfigure = &gameCenterDataSource{}
)

type gameCenterDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	AppID            types.String `tfsdk:"app_id"`
	ArcadeEnabled    types.Bool   `tfsdk:"arcade_enabled"`
	ChallengeEnabled types.Bool   `tfsdk:"challenge_enabled"`
}

type gameCenterDataSource struct {
	client *appstore.Client
}

func NewGameCenterDataSource() datasource.DataSource {
	return &gameCenterDataSource{}
}

func (d *gameCenterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_game_center"
}

func (d *gameCenterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*appstore.Client)
}

func (d *gameCenterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches Game Center information from the App Store Connect.",
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Description: "Identifier of the app to fetch Game Center information for.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the game center.",
				Computed:    true,
			},
			"arcade_enabled": schema.BoolAttribute{
				Description: "Indicates whether Game Center is enabled for the app on Apple Arcade.",
				Computed:    true,
			},
			"challenge_enabled": schema.BoolAttribute{
				Description: "Indicates whether Game Center challenges are enabled for the app.",
				Computed:    true,
			},
		},
	}
}

func (d *gameCenterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := gameCenterDataSourceModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	appId := state.AppID.ValueString()
	if appId == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'app_id' is required to fetch Game Center information.",
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

	gameCenter, err := d.client.GetGameCenter(ctx, app)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Game Center",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(gameCenter.ID)
	state.ArcadeEnabled = types.BoolValue(gameCenter.Attr.ArcadeEnabled)
	state.ChallengeEnabled = types.BoolValue(gameCenter.Attr.ChallengeEnabled)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
