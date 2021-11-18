package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CBC_ACCESS_KEY", nil),
				Description: "Couchbase Capella API Access Key",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CBC_SECRET_KEY", nil),
				Description: "Couchbase Capella API Secret Key",
				Sensitive:   true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"couchbasecapella_project":       resourceCouchbaseCapellaProject(),
			"couchbasecapella_cluster":       resourceCouchbaseCapellaCluster(),
			"couchbasecapella_database_user": resourceCouchbaseCapellaDatabaseUser(),
			"couchbasecapella_bucket":        resourceCouchbaseCapellaBucket(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// TODO: create a client with access/secret keys
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	configuration := couchbasecapella.NewConfiguration()
	apiClient := couchbasecapella.NewAPIClient(configuration)
	return apiClient, nil
}

func getAuth(ctx context.Context) context.Context {
	auth := context.WithValue(
		context.Background(),
		couchbasecapella.ContextAPIKeys,
		map[string]couchbasecapella.APIKey{
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
