package provider

import (
	"context"

	"github.com/d-asmaa/couchbase-cloud-go-client/couchbasecloud"
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
	client := meta.(*couchbasecloud.CouchbaseCloudClient)

	cloudId := d.Id()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := client.GetCloud(cloudId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
