// Couchbase, Inc. licenses this to you under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0.

// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func resourceCouchbaseCapellaDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Capella Users.",

		CreateContext: resourceCouchbaseCapellaDatabaseUserCreate,
		ReadContext:   resourceCouchbaseCapellaDatabaseUserRead,
		UpdateContext: resourceCouchbaseCapellaDatabaseUserUpdate,
		DeleteContext: resourceCouchbaseCapellaDatabaseUserDelete,

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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(25 * time.Minute),
		},
	}
}

func resourceCouchbaseCapellaDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	createDatabaseUserRequest := *couchbasecapella.NewCreateDatabaseUserRequest(username, password)
	_, allBucketAccessOk := d.GetOk("all_bucket_access")
	_, bucketsOk := d.GetOk("buckets")

	if !allBucketAccessOk && !bucketsOk {
		return diag.Errorf("No bucket roles specified")
	}

	if allBucketAccessOk && !bucketsOk {
		allBucketAccess := couchbasecapella.BucketRoleTypes(d.Get("all_bucket_access").(string))
		createDatabaseUserRequest.SetAllBucketsAccess(allBucketAccess)
	}

	if !allBucketAccessOk && bucketsOk {
		buckets := expandBuckets(d)
		createDatabaseUserRequest.SetBuckets(buckets)
	}

	if allBucketAccessOk && bucketsOk {
		return diag.Errorf("Please specify only specific buckets or all buckets")
	}

	r, err := client.ClustersApi.ClustersCreateUser(auth, clusterId).CreateDatabaseUserRequest(createDatabaseUserRequest).Execute()
	if err != nil {
		return manageErrors(err, *r, "Create Database User")
	}

	d.SetId(username)

	return resourceCouchbaseCapellaDatabaseUserRead(ctx, d, meta)
}

func resourceCouchbaseCapellaDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
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

func resourceCouchbaseCapellaDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)
	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	updateDatabaseUserRequest := *couchbasecapella.NewUpdateDatabaseUserRequest()

	if d.HasChange("all_bucket_access") {
		allBucketAccess := couchbasecapella.BucketRoleTypes(d.Get("all_bucket_access").(string))
		updateDatabaseUserRequest.SetAllBucketsAccess(allBucketAccess)
	} else if d.HasChange("buckets") {
		buckets := expandBuckets(d)
		updateDatabaseUserRequest.SetBuckets(buckets)
	}

	r, err := client.ClustersApi.ClustersUpdateUser(auth, clusterId, username).UpdateDatabaseUserRequest(updateDatabaseUserRequest).Execute()
	if err != nil {
		return manageErrors(err, *r, "Update Database User")
	}

	return resourceCouchbaseCapellaDatabaseUserRead(ctx, d, meta)
}

func resourceCouchbaseCapellaDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)
	username := d.Get("username").(string)

	r, err := client.ClustersApi.ClustersDeleteUser(auth, clusterId, username).Execute()
	if err != nil {
		return manageErrors(err, *r, "Delete Database User")
	}

	return nil
}

func expandBuckets(d *schema.ResourceData) []couchbasecapella.BucketRole {
	buckets := make([]couchbasecapella.BucketRole, 0)

	if v, ok := d.GetOk("buckets"); ok {
		for _, s := range v.(*schema.Set).List() {
			bucketMap := s.(map[string]interface{})

			bucketAccess := expandBucketAccessList(bucketMap["bucket_access"].([]interface{}))

			bucket := couchbasecapella.BucketRole{
				BucketName:   bucketMap["bucket_name"].(string),
				BucketAccess: bucketAccess,
			}
			buckets = append(buckets, bucket)
		}
	}

	return buckets
}

func expandBucketAccessList(bucketAccess []interface{}) (res []couchbasecapella.BucketRoleTypes) {
	for _, v := range bucketAccess {
		res = append(res, couchbasecapella.BucketRoleTypes(v.(string)))
	}

	return res
}
