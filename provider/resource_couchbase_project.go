package provider

import (
	"context"

	couchbasecloud "github.com/d-asmaa/couchbase-cloud-go-client/couchbasecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCouchbaseProject() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Projects.",

		CreateContext: resourceCouchbaseProjectCreate,
		ReadContext:   resourceCouchbaseProjectRead,
		DeleteContext: resourceCouchbaseProjectDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Project's id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "Project's name.",
				Type:        schema.TypeString,
				Optional:    false,
			},
		},
	}
}

func resourceCouchbaseProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.CouchbaseCloudClient)

	payload := &couchbasecloud.CreateProjectPayload{
		Name: d.Get("name").(string),
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := client.CreateProject(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	resourceCouchbaseProjectRead(ctx, d, meta)

	return diags
}

func resourceCouchbaseProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.CouchbaseCloudClient)

	projectId := d.Id()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := client.GetProject(projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCouchbaseProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.CouchbaseCloudClient)

	projectId := d.Id()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := client.DeleteProject(projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
