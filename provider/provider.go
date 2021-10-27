package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecloud "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

// Provider returns the provider to be use by the code.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CBC_ACCESS_KEY", nil),
				Description: "Couchbase Cloud API Access Key",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CBC_SECRET_KEY", nil),
				Description: "Couchbase Cloud API Secret Key",
				Sensitive:   true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"couchbasecloud_project": resourceCouchbaseCloudProject(),
			"couchbasecloud_cluster": resourceCouchbaseCluster(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// TODO: create a client with access/secret keys
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	configuration := couchbasecloud.NewConfiguration()
	apiClient := couchbasecloud.NewAPIClient(configuration)
	return apiClient, nil
}

// TODO: Get Auth with both env variable and terraform ones
func getAuth(ctx context.Context) context.Context {
	auth := context.WithValue(
		context.Background(),
		couchbasecloud.ContextAPIKeys,
		map[string]couchbasecloud.APIKey{
			"accessKey": {
				Key: os.Getenv("CBC_ACCESS_KEY"),
			},
			"secretKey": {
				Key: os.Getenv("CBC_SECRET_KEY"),
			},
		},
	)
	return auth
}