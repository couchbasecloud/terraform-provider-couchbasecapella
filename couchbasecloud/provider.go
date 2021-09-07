package couchbasecloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"template_resource": resourceTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"template_data_source": dataSourceTemplate(),
		},
	}
}
