package provider

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/alexprogrammr/appstore-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &appstoreProvider{}
)

type appstoreProvider struct {
	version string
}

type appstoreProviderModel struct {
	KeyID      types.String `tfsdk:"key_id"`
	IssuerID   types.String `tfsdk:"issuer_id"`
	PrivateKey types.String `tfsdk:"private_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &appstoreProvider{
			version: version,
		}
	}
}

func (p *appstoreProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "appstore"
	resp.Version = p.version
}

func (p *appstoreProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with App Store Connect.",
		Attributes: map[string]schema.Attribute{
			"key_id": schema.StringAttribute{
				Description: "Private key ID from App Store Connect, for example, 2X9R4HXF34.",
				Required:    true,
			},
			"issuer_id": schema.StringAttribute{
				Description: "Issuer ID from the API Keys page in App Store Connect, for example, 57246542-96fe-1a63-e053-0824d011072a.",
				Required:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "PEM-encoded private key from App Store Connect. " +
					"Keep your API keys secure and private. Don’t share your keys, store keys in a code repository, or include keys in client-side code. " +
					"If the key becomes lost or compromised, remember to revoke it immediately.",
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *appstoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring App Store Connect API client")

	config := appstoreProviderModel{}
	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.KeyID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key_id"),
			"Unknown Key ID",
			"The provider cannot create the App Store Connect API client as there is an unknown configuration value for the Key ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KEY_ID environment variable.",
		)
	}
	if config.IssuerID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("issuer_id"),
			"Unknown Issuer ID",
			"The provider cannot create the App Store Connect API client as there is an unknown configuration value for the Issuer ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ISSUER_ID environment variable.",
		)
	}
	if config.PrivateKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("private_key"),
			"Unknown Private Key",
			"The provider cannot create the App Store Connect API client as there is an unknown configuration value for the Private Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PRIVATE_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	keyID := os.Getenv("KEY_ID")
	issuerID := os.Getenv("ISSUER_ID")
	privateKey := os.Getenv("PRIVATE_KEY")

	if !config.KeyID.IsNull() {
		keyID = config.KeyID.ValueString()
	}
	if !config.IssuerID.IsNull() {
		issuerID = config.IssuerID.ValueString()
	}
	if !config.PrivateKey.IsNull() {
		privateKey = config.PrivateKey.ValueString()
	}

	if keyID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("key_id"),
			"Missing Key ID",
			"The provider cannot create the App Store Connect API client as there is a missing or empty value for the Key ID. "+
				"Set the host value in the configuration or use the KEY_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if issuerID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("issuer_id"),
			"Missing Issuer ID",
			"The provider cannot create the App Store Connect API client as there is a missing or empty value for the Issuer ID. "+
				"Set the host value in the configuration or use the ISSUER_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if privateKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("private_key"),
			"Missing Private Key",
			"The provider cannot create the App Store Connect API client as there is a missing or empty value for the Private Key. "+
				"Set the host value in the configuration or use the PRIVATE_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "key_id", keyID)
	ctx = tflog.SetField(ctx, "issuer_id", issuerID)
	ctx = tflog.SetField(ctx, "private_key", privateKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "private_key")

	tflog.Debug(ctx, "Creating App Store Connect API client")

	source, err := appstore.NewTokenSource(appstore.Config{
		KeyID:       keyID,
		IssuerID:    issuerID,
		PrivateKey:  []byte(privateKey),
		ExpireAfter: time.Minute,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create App Store Connect API client",
			"An unexpected error occurred when creating the App Store Connect API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"App Store Connect Client Error: "+err.Error(),
		)
		return
	}

	client := appstore.NewClient(http.DefaultClient, source)

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured App Store Connect API client")
}

func (p *appstoreProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAppsDataSource,
		NewAppDataSource,
		NewGameCenterDataSource,
	}
}

func (p *appstoreProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAchievementResource,
		NewAchievementLocalizationResource,
		NewAchievementImageResource,
	}
}
