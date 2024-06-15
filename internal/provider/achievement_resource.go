package provider

import (
	"context"

	"github.com/alexprogrammr/appstore-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &achievementResource{}
	_ resource.ResourceWithConfigure = &achievementResource{}
)

type achievementResourceModel struct {
	ID               types.String `tfsdk:"id"`
	GameCenterID     types.String `tfsdk:"game_center_id"`
	ReferenceName    types.String `tfsdk:"reference_name"`
	VendorID         types.String `tfsdk:"vendor_id"`
	Points           types.Int64  `tfsdk:"points"`
	Repeatable       types.Bool   `tfsdk:"repeatable"`
	ShowBeforeEarned types.Bool   `tfsdk:"show_before_earned"`
}

type achievementResource struct {
	client *appstore.Client
}

func NewAchievementResource() resource.Resource {
	return &achievementResource{}
}

func (r *achievementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_achievement"
}

func (d *achievementResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*appstore.Client)
}

func (r *achievementResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages game center achievement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the achievement.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"game_center_id": schema.StringAttribute{
				Description: "Identifier of the game center to associate the achievement with. Resource will be re-created if this value is changed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"reference_name": schema.StringAttribute{
				Description: "An internal name of the achievement.",
				Required:    true,
			},
			"vendor_id": schema.StringAttribute{
				Description: "A chosen alphanumeric identifier of the achievement. Resource will be re-created if this value is changed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"points": schema.Int64Attribute{
				Description: "The points that each achievement is worth.",
				Required:    true,
			},
			"repeatable": schema.BoolAttribute{
				Description: "An indication of whether the player can earn the achievement multiple times.",
				Required:    true,
			},
			"show_before_earned": schema.BoolAttribute{
				Description: "An indication of whether the achievement is visible to the player before it is earned.",
				Required:    true,
			},
		},
	}
}

func (r *achievementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state := achievementResourceModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gameCenterId := state.GameCenterID.ValueString()
	if gameCenterId == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'game_center_id' is required to create an achievement.",
		)
		return
	}

	gameCenter, err := r.client.GetGameCenterByID(ctx, gameCenterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read game center",
			err.Error(),
		)
		return
	}

	achievement := appstore.Achievement{
		ReferenceName:    state.ReferenceName.ValueString(),
		VendorIdentifier: state.VendorID.ValueString(),
		Points:           int(state.Points.ValueInt64()),
		Repeatable:       state.Repeatable.ValueBool(),
		ShowBeforeEarned: state.ShowBeforeEarned.ValueBool(),
	}

	response, err := r.client.CreateAchievement(ctx, gameCenter, &achievement)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create achievement",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(response.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *achievementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := achievementResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	achievement, err := r.client.GetAchievementByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read achievement",
			err.Error(),
		)
		return
	}

	state.ReferenceName = types.StringValue(achievement.Attr.ReferenceName)
	state.VendorID = types.StringValue(achievement.Attr.VendorIdentifier)
	state.Points = types.Int64Value(int64(achievement.Attr.Points))
	state.Repeatable = types.BoolValue(achievement.Attr.Repeatable)
	state.ShowBeforeEarned = types.BoolValue(achievement.Attr.ShowBeforeEarned)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *achievementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := achievementResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateAchievement(ctx, appstore.AchievementUpdate{
		ID:               plan.ID.ValueString(),
		ReferenceName:    plan.ReferenceName.ValueString(),
		Points:           int(plan.Points.ValueInt64()),
		Repeatable:       plan.Repeatable.ValueBool(),
		ShowBeforeEarned: plan.ShowBeforeEarned.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update achievement",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *achievementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := achievementResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAchievementByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete achievement",
			err.Error(),
		)
		return
	}
}
