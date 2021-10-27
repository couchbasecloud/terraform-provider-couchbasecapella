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
		// UpdateContext: resourceCouchbaseCloudDatabaseUserUpdate,
		DeleteContext: resourceCouchbaseCloudDatabaseUserDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Description: "Cluster ID",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"username": {
				Description: "Database user username",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"password": {
				Description: "Database user password",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"buckets": {
				Description: "Database user bucket access",
				Type:        schema.TypeSet,
				ForceNew:    true,
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_role_type": {
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
			"all_bucket_access": {
				Description: "Database user all bucket access",
				Type:        schema.TypeString,
				ForceNew:    true,
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
	buckets := expandBuckets(d)

	// need to add logic to check for value
	//allBucketAccess := couchbasecloud.BucketRoleTypes(d.Get("all_bucket_access").(string))

	createDatabaseUserRequest := *couchbasecloud.NewCreateDatabaseUserRequest(username, password)
	createDatabaseUserRequest.SetBuckets(buckets)
	//createDatabaseUserRequest.SetAllBucketsAccess(allBucketAccess)

	_, err := client.ClustersApi.ClustersCreateUser(auth, clusterId).CreateDatabaseUserRequest(createDatabaseUserRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterId)

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

// func resourceCouchbaseCloudUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	auth := context.WithValue(
// 		context.Background(),
// 		couchbasecloud.ContextAPIKeys,
// 		map[string]couchbasecloud.APIKey{
// 			"accessKey": {
// 				Key: os.Getenv("CBC_ACCESS_KEY"),
// 			},
// 			"secretKey": {
// 				Key: os.Getenv("CBC_SECRET_KEY"),
// 			},
// 		},
// 	)
// 	client := meta.(*couchbasecloud.APIClient)
// 	clusterId := d.Get("cluster_id").(string)
// 	username := d.Get("username").(string)

// 	_, err := client.ClustersApi.ClustersUpdateUser(auth, clusterId, username).Execute()
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId(clusterId)

// 	return resourceCouchbaseCloudUserRead(ctx, d, meta)
// }

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
	var buckets []couchbasecloud.BucketRole

	if v, ok := d.GetOk("buckets"); ok {
		if rs := v.(*schema.Set); rs.Len() > 0 {
			buckets = make([]couchbasecloud.BucketRole, rs.Len())

			for k, r := range rs.List() {
				bucketsMap := r.(map[string]interface{})

				buckets[k] = couchbasecloud.BucketRole{
					BucketName:   bucketsMap["bucket_name"].(string),
					BucketAccess: expandBucketAccess(d),
				}
			}
		}
	}

	return buckets
}

func expandBucketAccess(d *schema.ResourceData) []couchbasecloud.BucketRoleTypes {
	var bucketAccess []couchbasecloud.BucketRoleTypes

	if v, ok := d.GetOk("bucket_access"); ok {
		if rs := v.(*schema.Set); rs.Len() > 0 {
			bucketAccess = make([]couchbasecloud.BucketRoleTypes, rs.Len())

			for k, r := range rs.List() {
				rolesMap := r.(map[string]couchbasecloud.BucketRoleTypes)
				bucketAccess[k] = couchbasecloud.BucketRoleTypes(rolesMap["bucket_role_type"])
			}
		}
	}

	return bucketAccess
}
