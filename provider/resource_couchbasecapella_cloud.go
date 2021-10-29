package provider

import (
	"context"
	"net/http"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCouchbaseCapellaCloud() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Capella Clouds.",

		ReadContext: resourceCouchbaseCapellaCloudRead,

		Schema: map[string]*schema.Schema{
			"cloud_id": {
				Description: "Cloud id.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCouchbaseCapellaCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	cloudId := d.Get("cloud_id").(string)

	cloud, resp, err := client.CloudsApi.CloudsShow(auth, cloudId).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if err := d.Set("name", cloud.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("provider", cloud.Provider); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
