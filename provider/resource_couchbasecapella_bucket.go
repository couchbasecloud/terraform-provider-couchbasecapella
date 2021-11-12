package provider

import (
	"context"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
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
			"cluster_id": {
				Description: "Cluster's id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Bucket's name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"memory_quota": {
				Description: "Bucket Memory quota.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"replicas": {
				Description: "replicas.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"conflict_resolution": {
				Description: "replicas.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCouchbaseCapellaBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)
	bucketName := d.Get("name").(string)
	memoryQuota := int32(d.Get("memory_quota").(int))
	replicas := int32(d.Get("replicas").(int))
	conflictResolution := couchbasecapella.ConflictResolution(d.Get("conflict_resolution").(string))

	couchbaseBucketSpec := *couchbasecapella.NewCouchbaseBucketSpec(bucketName, memoryQuota)
	couchbaseBucketSpec.SetReplicas(replicas)
	couchbaseBucketSpec.SetConflictResolution(conflictResolution)

	_, _, err := client.ClustersApi.ClustersCreateBucket(auth, clusterId).CouchbaseBucketSpec(couchbaseBucketSpec).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bucketName)

	return resourceCouchbaseCapellaBucketRead(ctx, d, meta)
}

func resourceCouchbaseCapellaBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)
	buckets, _, err := client.ClustersApi.ClustersListBuckets(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, bucket := range buckets {
		if bucket.Name == d.Get("name") {
			return nil
		}
	}
	return diag.FromErr(err)
}

func resourceCouchbaseCapellaBucketUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("This current release of terraform provider doesn't support updating the bucket, Please log to your Capella UI account")
}

func resourceCouchbaseCapellaBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)
	bucketName := d.Get("name").(string)

	deleteBucketRequest := *couchbasecapella.NewDeleteBucketRequest(bucketName)

	_, err := client.ClustersApi.ClustersDeleteBucket(auth, clusterId).DeleteBucketRequest(deleteBucketRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
