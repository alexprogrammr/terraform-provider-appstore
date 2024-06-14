package provider

import (
	"context"

	"github.com/alexprogrammr/appstore-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &achievementLocalizationResource{}
	_ resource.ResourceWithConfigure = &achievementLocalizationResource{}
)

type achievementLocalizationResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AchievementID types.String `tfsdk:"achievement_id"`
	Locale        types.String `tfsdk:"locale"`
	Name          types.String `tfsdk:"name"`
	BeforeEarned  types.String `tfsdk:"before_earned_description"`
	AfterEarned   types.String `tfsdk:"after_earned_description"`
}

type achievementLocalizationResource struct {
	client *appstore.Client
}

func NewAchievementLocalizationResource() resource.Resource {
	return &achievementLocalizationResource{}
}

func (r *achievementLocalizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_achievement_localization"
}

func (r *achievementLocalizationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*appstore.Client)
}

func (r *achievementLocalizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages game center achievement localization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the achievement localization.",
				Computed:    true,
			},
			"achievement_id": schema.StringAttribute{
				Description: "Identifier of the achievement to associate the localization with.",
				Required:    true,
			},
			"locale": schema.StringAttribute{
				Description: "Locale of the achievement localization.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the achievement.",
				Required:    true,
			},
			"before_earned_description": schema.StringAttribute{
				Description: "Description of the achievement before it is earned.",
				Required:    true,
			},
			"after_earned_description": schema.StringAttribute{
				Description: "Description of the achievement after it is earned.",
				Required:    true,
			},
		},
	}
}

func (r *achievementLocalizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state := achievementLocalizationResourceModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	achievementID := state.AchievementID.ValueString()
	if achievementID == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'achievement_id' is required to create an achievement localization.",
		)
		return
	}

	achievement, err := r.client.GetAchievementByID(ctx, achievementID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read achievement",
			err.Error(),
		)
		return
	}

	localization, err := r.client.CreateAchievementLocalization(ctx, achievement, &appstore.AchievementLocalization{
		Locale:                  state.Locale.ValueString(),
		Name:                    state.Name.ValueString(),
		BeforeEarnedDescription: state.BeforeEarned.ValueString(),
		AfterEarnedDescription:  state.AfterEarned.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create achievement localization",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(localization.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *achievementLocalizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := achievementLocalizationResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	localization, err := r.client.GetAchievementLocalizationByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read achievement localization",
			err.Error(),
		)
		return
	}

	state.Locale = types.StringValue(localization.Attr.Locale)
	state.Name = types.StringValue(localization.Attr.Name)
	state.BeforeEarned = types.StringValue(localization.Attr.BeforeEarnedDescription)
	state.AfterEarned = types.StringValue(localization.Attr.AfterEarnedDescription)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *achievementLocalizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *achievementLocalizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := achievementLocalizationResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAchievementLocalizationByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to aelete achievement localization",
			err.Error(),
		)
		return
	}
}
