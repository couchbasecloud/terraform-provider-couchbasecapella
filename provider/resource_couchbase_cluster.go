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
				Required:    true,
			},
			"cloud_id": {
				Description: "Cloud's Id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_id": {
				Description: "Project's Id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"servers": {
				Description: "Cluster servers configuration",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Description: "Number of nodes",
							Required:    true,
						},
						"services": {
							Type:        schema.TypeList,
							Description: "Services",
							Required:    true,
							MinItems:    1,
							Default:     []couchbasecloud.CouchbaseServices{"data"},
						},
						"aws": {
							Description: "Aws configuration.",
							Type:        schema.TypeSet,
							Required:    false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_size": {
										Description: "Aws instance.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"ebs_size_gib": {
										Description: "Aws size(Gb).",
										Type:        schema.TypeInt,
										Required:    true,
									},
								},
							},
						},
						"Azure": {
							Description: "Azure configuration.",
							Type:        schema.TypeSet,
							Required:    false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_size": {
										Description: "Azure instance.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"ebs_size_gib": {
										Description: "Aws size(Gb).",
										Type:        schema.TypeInt,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceCouchbaseClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)

	clusterName := d.Get("name").(string)
	cloudId := d.Get("cloud_id").(string)
	projectId := d.Get("project_id").(string)

	newClusterRequest := *couchbasecloud.NewCreateClusterRequest(clusterName, cloudId, projectId)

	cluster, _, err := client.ClustersApi.ClustersCreate(auth).CreateClusterRequest(newClusterRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.Id)

	return resourceCouchbaseClusterRead(ctx, d, meta)
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
