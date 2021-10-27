package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecloud "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func resourceCouchbaseCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase clusters.",

		CreateContext: resourceCouchbaseClusterCreate,
		ReadContext:   resourceCouchbaseClusterRead,
		UpdateContext: resourceCouchbaseClusterUpdate,
		DeleteContext: resourceCouchbaseClusterDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Cluster's name.",
				Type:        schema.TypeString,
				Optional:    false,
			},
			"cloudId": {
				Description: "Cloud's Id.",
				Type:        schema.TypeString,
				Optional:    false,
			},
			"tenantId": {
				Description: "Tenant's Id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"projectId": {
				Description: "Project's Id.",
				Type:        schema.TypeString,
				Optional:    false,
			},
		},
	}
}

func resourceCouchbaseClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)

	newClusterRequest := *couchbasecloud.NewCreateClusterRequest(d.Get("name").(string), d.Get("cloudId").(string), d.Get("projectId").(string))

	cluster, _, err := client.ClustersApi.ClustersCreate(auth).CreateClusterRequest(newClusterRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.Id)

	return resourceCouchbaseClusterCreate(ctx, d, meta)
}

func resourceCouchbaseClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)
	clusterId := d.Get("id").(string)

	cluster, resp, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if err := d.Set("name", cluster.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCouchbaseClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceCouchbaseClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}
