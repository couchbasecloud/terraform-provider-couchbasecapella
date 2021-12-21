// Couchbase, Inc. licenses this to you under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0.

// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package provider

import (
	"context"
	"fmt"
	"time"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func resourceCouchbaseCapellaDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Capella Database Users",

		CreateContext: resourceCouchbaseCapellaDatabaseUserCreate,
		ReadContext:   resourceCouchbaseCapellaDatabaseUserRead,
		UpdateContext: resourceCouchbaseCapellaDatabaseUserUpdate,
		DeleteContext: resourceCouchbaseCapellaDatabaseUserDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Description: "ID of the Cluster",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					idIsValid := IsValidUUID(val.(string))
					if !idIsValid {
						errs = append(errs, fmt.Errorf("please enter a valid cluster uuid"))
					}
					return
				},
			},
			"username": {
				Description:  "Username for the Database User",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"password": {
				Description: "Password for the Database User",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					password := val.(string)
					passwordValidate := validatePassword(password)
					if !passwordValidate {
						errs = append(errs, fmt.Errorf("password must contain 8+ characters, 1+ lowercase, 1+ uppercase, 1+ symbols, 1+ numbers."))
					}
					return
				},
			},
			"buckets": {
				Description: "Define bucket access level for the Database User",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"bucket_access": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
									access := val.(string)
									accessValidation := couchbasecapella.BucketRoleTypes(access).IsValid()
									if !accessValidation {
										errs = append(errs, fmt.Errorf("please enter a valid value for bucket access {data_reader, data_writer}"))
									}
									return
								},
							},
						},
					},
				},
			},
			"all_bucket_access": {
				Description: "Define all bucket access for the Database User",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					access := val.(string)
					accessValidation := couchbasecapella.BucketRoleTypes(access).IsValid()
					if !accessValidation {
						errs = append(errs, fmt.Errorf("please enter a valid value for all bucket access {data_reader, data_writer}"))

					}
					return
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(25 * time.Minute),
		},
	}
}

// resourceCouchbaseCapellaDatabaseUserCreate is responsible for creating a
// database user in a Couchbase Capella VPC Cluster using the Terraform resource data.
// WARNING: Creating database users is only supported for VPC Clusters in this current
// release.
func resourceCouchbaseCapellaDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)

	// Check if the Cluster is inVPC to create the users
	// Managing buckets is not available for hosted clusters
	_, _, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		// Check V3Cluster :: Need to be fixed in next versions
		_, _, err3 := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()
		if err3 != nil {
			return diag.FromErr(fmt.Errorf("a problem occurred while accessing to the cluster"))
		}
		return diag.FromErr(fmt.Errorf("sorry, managing database users is not available for hosted clusters"))
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Check to see if a user with the same name already exists in the cluster. If a user
	// already has the name, an error is thrown. If not, then proceeds with creation.
	users, _, err := client.ClustersApi.ClustersListUsers(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, user := range users {
		if user.Username == username {
			return diag.Errorf("Failed to create: A user already exists with that name")
		}
	}

	createDatabaseUserRequest := *couchbasecapella.NewCreateDatabaseUserRequest(username, password)
	_, allBucketAccessExists := d.GetOk("all_bucket_access")
	_, bucketsExists := d.GetOk("buckets")

	// Only one of `buckets` or `all_bucket_access` can be specified in the terraform configuration file.
	// The follwing checks handle errors where none or both are specified. If only one of `buckets` or
	// `all_bucket_access` is specified, the bucket access roles will be set accordingly.
	if !allBucketAccessExists && !bucketsExists {
		return diag.Errorf("No bucket access roles specified")
	}

	if allBucketAccessExists && !bucketsExists {
		allBucketAccess := couchbasecapella.BucketRoleTypes(d.Get("all_bucket_access").(string))
		createDatabaseUserRequest.SetAllBucketsAccess(allBucketAccess)
	}

	if !allBucketAccessExists && bucketsExists {
		buckets := expandBuckets(d)
		createDatabaseUserRequest.SetBuckets(buckets)
	}

	if allBucketAccessExists && bucketsExists {
		return diag.Errorf("Please specify only access for specific buckets or access for all buckets")
	}

	r, err := client.ClustersApi.ClustersCreateUser(auth, clusterId).CreateDatabaseUserRequest(createDatabaseUserRequest).Execute()
	if r == nil {
		return diag.Errorf("Pointer to database user create http.Response is nil")
	}
	if err != nil {
		return manageErrors(err, *r, "Create Database User")
	}

	d.SetId(username)

	return resourceCouchbaseCapellaDatabaseUserRead(ctx, d, meta)
}

// resourceCouchbaseCapellaDatabaseUserRead is responsible for reading a Couchbase
// Capella database user using the Terraform resource data.
func resourceCouchbaseCapellaDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)
	clusterId := d.Get("cluster_id").(string)

	// Check if the Cluster is inVPC to read the db users
	// Managing buckets is not available for hosted clusters
	_, _, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		// Check V3Cluster :: Need to be fixed in next versions
		_, _, err3 := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()
		if err3 != nil {
			return diag.FromErr(fmt.Errorf("a problem occurred while accessing to the cluster"))
		}
		return diag.FromErr(fmt.Errorf("sorry, managing database users is not available for hosted clusters"))
	}

	// The current version of the Capella API doesn't support getting a singular
	// database user. To obtain the database user, we need to iterate the trough
	// the list of all database users and check if it exists. If the user
	// is not present in the list of users and error is thrown.
	users, _, err := client.ClustersApi.ClustersListUsers(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	userExists := false
	for _, user := range users {
		if user.Username == d.Id() {
			userExists = true
		}
	}
	if !userExists {
		username := d.Id()
		d.SetId("")
		return diag.Errorf("Error 404: Failed to find the username %s ", username)
	}
	return nil
}

// resourceCouchbaseCapellaDatabaseUserUpdate is responsible for updating a
// database user in a Couchbase Capella VPC Cluster using the Terraform resource data.
// WARNING: Updating database users is only supported for VPC Clusters in this current
// release.
func resourceCouchbaseCapellaDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)
	clusterId := d.Get("cluster_id").(string)

	// Check if the Cluster is inVPC to update the users
	// Managing buckets is not available for hosted clusters
	_, _, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		// Check V3Cluster :: Need to be fixed in next versions
		_, _, err3 := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()
		if err3 != nil {
			return diag.FromErr(fmt.Errorf("a problem occurred while accessing to the cluster"))
		}
		return diag.FromErr(fmt.Errorf("sorry, managing database users is not available for hosted clusters"))
	}

	username := d.Get("username").(string)

	updateDatabaseUserRequest := *couchbasecapella.NewUpdateDatabaseUserRequest()

	// Check to see if either `all_bucket_access` or `buckets` has changed as only one should
	// be present in the resource data. If there has been a change, then update accordingly.
	if d.HasChange("all_bucket_access") {
		allBucketAccess := couchbasecapella.BucketRoleTypes(d.Get("all_bucket_access").(string))
		updateDatabaseUserRequest.SetAllBucketsAccess(allBucketAccess)
	} else if d.HasChange("buckets") {
		buckets := expandBuckets(d)
		updateDatabaseUserRequest.SetBuckets(buckets)
	}

	r, err := client.ClustersApi.ClustersUpdateUser(auth, clusterId, username).UpdateDatabaseUserRequest(updateDatabaseUserRequest).Execute()
	if r == nil {
		return diag.Errorf("Pointer to database user update http.Response is nil")
	}
	if err != nil {
		return manageErrors(err, *r, "Update Database User")
	}

	return resourceCouchbaseCapellaDatabaseUserRead(ctx, d, meta)
}

// resourceCouchbaseCapellaDatabaseUserDelete is responsible for deleting a
// database user in a Couchbase Capella VPC Cluster using the Terraform resource data.
// WARNING: Deleting database users is only supported for VPC Clusters in this current
// release.
func resourceCouchbaseCapellaDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)

	// Check if the Cluster is inVPC to delete the users
	// Managing buckets is not available for hosted clusters
	_, _, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		// Check V3Cluster :: Need to be fixed in next versions
		_, _, err3 := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()
		if err3 != nil {
			return diag.FromErr(fmt.Errorf("a problem occurred while accessing to the cluster"))
		}
		return diag.FromErr(fmt.Errorf("sorry, managing database users is not available for hosted clusters"))
	}

	username := d.Get("username").(string)

	// Check to see if database user exists in list of database users. If the database user
	// exists, it will be deleted from the Cluster. If the database user does not appear in the list of users,
	// likely being deleted elsewhere, an error is thrown.
	users, _, err := client.ClustersApi.ClustersListUsers(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, user := range users {
		if user.Username == username {
			r, err := client.ClustersApi.ClustersDeleteUser(auth, clusterId, username).Execute()
			if r == nil {
				return diag.Errorf("Pointer to database user delete http.Response is nil")
			}
			if err != nil {
				return manageErrors(err, *r, "Delete Database User")
			}
			return nil
		}
	}
	return diag.Errorf("Failed to delete: Database User doesn't exist in list of users")
}

// expandBuckets is responsible for converting the bucket interface into
// a slice of type BucketRole
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

// expandBucketAccessList is responsible for converting the bucketAccess interface into
// a slice of type BucketRoleTypes
func expandBucketAccessList(bucketAccess []interface{}) (roles []couchbasecapella.BucketRoleTypes) {
	for _, v := range bucketAccess {
		roles = append(roles, couchbasecapella.BucketRoleTypes(v.(string)))
	}

	return roles
}

// validatePassword is responsible for checking if a password string matches the required
// format. A password must contain 8+ characters, 1+ lowercase, 1+ uppercase, 1+ symbols, 1+ numbers.
// If the password matches the required format, the function will return true.
// If the password fails to match the required format, the function will return false.
func validatePassword(password string) bool {
	chars, lower, upper, symbol, number := false, false, false, false, false
	letters := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsLower(c):
			lower = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			symbol = true
		case unicode.IsLetter(c):
			letters++
		case c == ' ':
			return false
		default:
			return false
		}
	}
	chars = letters >= 8
	if chars && lower && upper && symbol && number {
		return true
	} else {
		return false
	}
}
