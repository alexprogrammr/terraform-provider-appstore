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
		Attributes: map[string]schema.Attribute{
			"key_id": schema.StringAttribute{
				Required: true,
			},
			"issuer_id": schema.StringAttribute{
				Required: true,
			},
			"private_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *appstoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
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

	source, err := appstore.NewTokenSource(appstore.Config{
		KeyID:       keyID,
		IssuerID:    issuerID,
		PrivateKey:  []byte(privateKey),
		ExpireAfter: time.Minute,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create App Store Connect API client",
			"An unexpected error occurred when creating the App Store Connect API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"App Store Connect Client Error: "+err.Error(),
		)
		return
	}

	client := appstore.NewClient(http.DefaultClient, source)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *appstoreProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAppsDataSource,
	}
}

func (p *appstoreProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
