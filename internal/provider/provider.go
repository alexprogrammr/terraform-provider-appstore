package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &appstoreProvider{}
)

type appstoreProvider struct {
	version string
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
	resp.Schema = schema.Schema{}
}

func (p *appstoreProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {

}

func (p *appstoreProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *appstoreProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
