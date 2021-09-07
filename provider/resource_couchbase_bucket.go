package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCouchbaseBucket() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Buckets.",

		CreateContext: resourceCouchbaseBucketCreate,
		ReadContext:   resourceCouchbaseBucketRead,
		UpdateContext: resourceCouchbaseBucketUpdate,
		DeleteContext: resourceCouchbaseBucketDelete,

		Schema: map[string]*schema.Schema{
		    "id": {
                Description: "Cluster's id.",
                Type:        schema.TypeString,
                Optional:    false,
        	},
			"name": {
				Description: "Bucket's name.",
				Type:        schema.TypeString,
				Optional:    false,
			},
			"memoryQuota": {
                Description: "Memory's quote.",
                Type:        schema.TypeInt,
                Optional:    false,
            },
            "replicas": {
                Description: "replicas.",
                Type:        schema.TypeInt,
                Optional:    false,
            },
		},
	}
}

func resourceCouchbaseBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	idFromAPI := "my-id"
	d.SetId(idFromAPI)

	return diag.Errorf("not implemented")
}

func resourceCouchbaseBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceCouchbaseBucketUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourcCouchbaseBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}