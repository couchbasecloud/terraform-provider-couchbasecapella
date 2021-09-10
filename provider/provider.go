package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Couchbase Cloud Base URL",
			},
			"acesss_key": {
				Type:        schema.TypeString,
				Required:    true,
				Default:     "",
				Description: "Couchbase Cloud API Access Key",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				Default:     "",
				Description: "Couchbase Cloud API Secret Key",
				Sensitive:   true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"couchbase_data_source": dataSourceCouchbase(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"couchbase_resource": resourceCouchbase(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := Client{
		accessKey: d.Get("access_key").(string),
		secretKey: d.Get("secret_key").(string),
		baseURL:   d.Get("base_url").(string),
	}
	return NewClient(baseURL, accessKey, secretKey), nil
}

type Client struct {
	baseURL    string
	accessKey  string
	secretKey  string
	httpClient *http.Client
}

func NewClient(baseURL, access, secret string) *Client {
	return &Client{
		baseURL:    baseURL,
		accessKey:  access,
		secretKey:  secret,
		httpClient: http.DefaultClient,
	}
}
