package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecloud "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func resourceCouchbaseCloudDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Cloud Users.",

		CreateContext: resourceCouchbaseCloudDatabaseUserCreate,
		ReadContext:   resourceCouchbaseCloudDatabaseUserRead,
		UpdateContext: resourceCouchbaseCloudDatabaseUserUpdate,
		DeleteContext: resourceCouchbaseCloudDatabaseUserDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Description: "Cluster ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"username": {
				Description: "Database user username",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "Database user password",
				Type:        schema.TypeString,
				Required:    true,
			},
			"buckets": {
				Description: "Database user bucket access",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bucket_access": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"all_bucket_access": {
				Description: "Database user all bucket access",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceCouchbaseCloudDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	createDatabaseUserRequest := *couchbasecloud.NewCreateDatabaseUserRequest(username, password)
	_, allBucketAccessOk := d.GetOk("all_bucket_access")
	_, bucketsOk := d.GetOk("buckets")

	if !allBucketAccessOk && !bucketsOk {
		return diag.Errorf("No bucket roles specified")
	}

	if allBucketAccessOk && !bucketsOk {
		allBucketAccess := couchbasecloud.BucketRoleTypes(d.Get("all_bucket_access").(string))
		createDatabaseUserRequest.SetAllBucketsAccess(allBucketAccess)
	}

	if !allBucketAccessOk && bucketsOk {
		buckets := expandBuckets(d)
		createDatabaseUserRequest.SetBuckets(buckets)
	}

	if allBucketAccessOk && bucketsOk {
		return diag.Errorf("Please specify only specific buckets or all buckets")
	}

	_, err := client.ClustersApi.ClustersCreateUser(auth, clusterId).CreateDatabaseUserRequest(createDatabaseUserRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(username)

	return nil
}

func resourceCouchbaseCloudDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)
	users, _, err := client.ClustersApi.ClustersListUsers(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, user := range users {
		if user.Username == d.Get("username") {
			return nil
		}
	}
	return diag.FromErr(err)
}

func resourceCouchbaseCloudDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	updateDatabaseUserRequest := *couchbasecloud.NewUpdateDatabaseUserRequest()

	if d.HasChange("all_bucket_access") {
		allBucketAccess := couchbasecloud.BucketRoleTypes(d.Get("all_bucket_access").(string))
		updateDatabaseUserRequest.SetAllBucketsAccess(allBucketAccess)
	} else if d.HasChange("buckets") {
		buckets := expandBuckets(d)
		updateDatabaseUserRequest.SetBuckets(buckets)
	}

	_, err := client.ClustersApi.ClustersUpdateUser(auth, clusterId, username).UpdateDatabaseUserRequest(updateDatabaseUserRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCouchbaseCloudDatabaseUserRead(ctx, d, meta)
}

func resourceCouchbaseCloudDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	_, err := client.ClustersApi.ClustersDeleteUser(auth, clusterId, username).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandBuckets(d *schema.ResourceData) []couchbasecloud.BucketRole {
	buckets := make([]couchbasecloud.BucketRole, 0)

	if v, ok := d.GetOk("buckets"); ok {
		for _, s := range v.(*schema.Set).List() {
			bucketMap := s.(map[string]interface{})

			bucketAccess := expandBucketAccessList(bucketMap["bucket_access"].([]interface{}))

			bucket := couchbasecloud.BucketRole{
				BucketName:   bucketMap["bucket_name"].(string),
				BucketAccess: bucketAccess,
			}
			buckets = append(buckets, bucket)
		}
	}

	return buckets
}

func expandBucketAccessList(bucketAccess []interface{}) (res []couchbasecloud.BucketRoleTypes) {
	for _, v := range bucketAccess {
		res = append(res, couchbasecloud.BucketRoleTypes(v.(string)))
	}

	return res
}
