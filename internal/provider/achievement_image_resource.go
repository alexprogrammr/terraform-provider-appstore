package provider

import (
	"context"
	"os"
	"path/filepath"

	"github.com/alexprogrammr/appstore-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &achievementImageResource{}
	_ resource.ResourceWithConfigure = &achievementImageResource{}
)

type achievementImageResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AchievementID types.String `tfsdk:"achievement_localization_id"`
	File          types.String `tfsdk:"file"`
	Checksum      types.String `tfsdk:"checksum"`
}

type achievementImageResource struct {
	client *appstore.Client
}

func NewAchievementImageResource() resource.Resource {
	return &achievementImageResource{}
}

func (r *achievementImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_achievement_image"
}

func (r *achievementImageResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*appstore.Client)
}

func (r *achievementImageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages game center achievement localization images.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the achievement image.",
				Computed:    true,
			},
			"achievement_localization_id": schema.StringAttribute{
				Description: "Identifier of the achievement localization to associate the image with.",
				Required:    true,
			},
			"file": schema.StringAttribute{
				Description: "Path to the image file.",
				Required:    true,
			},
			"checksum": schema.StringAttribute{
				Description: "MD5 checksum of the image.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *achievementImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state := achievementImageResourceModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	achievementId := state.AchievementID.ValueString()
	if achievementId == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'achievement_localization_id' is required to create an achievement image.",
		)
		return
	}

	filePath := state.File.ValueString()
	if filePath == "" {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Attribute 'file' is required to create an achievement image.",
		)
		return
	}

	image, err := os.ReadFile(filePath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read image file",
			err.Error(),
		)
		return
	}

	achievement, err := r.client.GetAchievementLocalizationByID(ctx, achievementId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get achievement localization",
			err.Error(),
		)
		return
	}

	asset, err := r.client.CreateAchievementImage(ctx, achievement, filepath.Base(filePath), image)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create achievement image",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(asset.ID)
	state.Checksum = types.StringValue(checksum(image))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *achievementImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := achievementImageResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	image, err := os.ReadFile(state.File.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read image file",
			err.Error(),
		)
		return
	}

	if state.Checksum.ValueString() != checksum(image) {
		state.File = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *achievementImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *achievementImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := achievementImageResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAchievementImageByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete achievement image",
			err.Error(),
		)
		return
	}
}
