package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCouchbaseCapellaBucket() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Buckets.",

		CreateContext: resourceCouchbaseCapellaBucketCreate,
		ReadContext:   resourceCouchbaseCapellaBucketRead,
		UpdateContext: resourceCouchbaseCapellaBucketUpdate,
		DeleteContext: resourceCouchbaseCapellaBucketDelete,

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
				Description: "Bucket Memory quota.",
				Type:        schema.TypeInt,
				Optional:    false,
			},
			"replicas": {
				Description: "replicas.",
				Type:        schema.TypeInt,
				Optional:    false,
			},
			"conflict_resolution": {
				Description: "replicas.",
				Type:        schema.TypeString,
				Optional:    false,
			},
		},
	}
}

func resourceCouchbaseCapellaBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	idFromAPI := "my-id"
	d.SetId(idFromAPI)

	return diag.Errorf("not implemented")
}

func resourceCouchbaseCapellaBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceCouchbaseCapellaBucketUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceCouchbaseCapellaBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}
