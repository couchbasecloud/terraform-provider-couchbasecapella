package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCouchbaseCloud() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Clouds.",

		ReadContext: resourceCouchbaseCloudRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Cloud id.",
				Type:        schema.TypeString,
				Optional:    false,
			},
		},
	}
}

func resourceCouchbaseCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
