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
	"regexp"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCouchbaseCapellaBucket() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Capella Buckets",

		CreateContext: resourceCouchbaseCapellaBucketCreate,
		ReadContext:   resourceCouchbaseCapellaBucketRead,
		UpdateContext: resourceCouchbaseCapellaBucketUpdate,
		DeleteContext: resourceCouchbaseCapellaBucketDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Description: "ID of the Cluster",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					idIsValid := IsValidUUID(val.(string))
					if !idIsValid {
						errs = append(errs, fmt.Errorf("please enter a valid cluster uuid"))
					}
					return
				},
			},
			"name": {
				Description: "Name of the Bucket",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9_.]*$`).MatchString
					var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString
					name := val.(string)
					nameValidate := isStringAlphabetic(name) && len(name) > 0 && len(name) < 100 && isAlphaNumeric(name[0:1])
					if !nameValidate {
						errs = append(errs, fmt.Errorf("use letters, numbers, periods (.) or dashes (- ). Bucket names cannot exceed 100 characters and must begin with a letter or a number"))
					}
					return
				},
			},
			"memory_quota": {
				Description: "Bucket Memory quota in Mb",
				Type:        schema.TypeInt,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					memory := val.(int)
					if memory < 100 {
						errs = append(errs, fmt.Errorf("please enter a memory value greater than 100 MiB"))
					}
					return
				},
			},
			"replicas": {
				Description: "Number of bucket replicas.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"conflict_resolution": {
				Description: "Conflict resolution for bucket",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					conflict := val.(string)
					conflictValidation := couchbasecapella.ConflictResolution(conflict).IsValid()
					if !conflictValidation {
						errs = append(errs, fmt.Errorf("please enter a valid value for conflict resolution {lww, seqno}"))

					}
					return
				},
			},
		},
	}
}

/**
*** Creating the Bucket
**/
func resourceCouchbaseCapellaBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)

	// Check if the Cluster is inVPC to create the bucket
	// Managing buckets is not available for hosted clusters
	_, _, clusterError := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if clusterError != nil {
		// Check V3Cluster :: Need to be fixed in next versions
		_, _, err3 := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()
		if err3 != nil {
			return diag.FromErr(fmt.Errorf("a problem occured while accessing to the cluster"))
		}
		return diag.FromErr(fmt.Errorf("sorry, managing buckets is not available for hosted clusters"))
	}

	bucketName := d.Get("name").(string)
	memoryQuota := int32(d.Get("memory_quota").(int))
	replicas := int32(d.Get("replicas").(int))
	conflictResolution := couchbasecapella.ConflictResolution(d.Get("conflict_resolution").(string))

	couchbaseBucketSpec := couchbasecapella.NewCouchbaseBucketSpec(bucketName, memoryQuota)
	couchbaseBucketSpec.SetReplicas(replicas)
	couchbaseBucketSpec.SetConflictResolution(conflictResolution)

	_, r, err := client.ClustersApi.ClustersCreateBucket(auth, clusterId).CouchbaseBucketSpec(*couchbaseBucketSpec).Execute()
	if err != nil {
		if r != nil {
			// TODO:  HANDLE ERRORS GRACEFULLY HERE AND REPORT AN ERROR THAT
			//  MAKES SENSE TO THE USER SO THEY KNOW HOW TO FIX THE PROBLEM
			// diag.FromErr((handleResponse(r))
		}
		return diag.FromErr(fmt.Errorf("problem occured :: %v", r))
	}

	d.SetId(bucketName)

	return resourceCouchbaseCapellaBucketRead(ctx, d, meta)
}

func resourceCouchbaseCapellaBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)

	// Check if the Cluster is inVPC to read the bucket list
	// Managing buckets is not available for hosted clusters
	_, _, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		// Check V3Cluster :: Need to be fixed in next versions
		_, _, err3 := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()
		if err3 != nil {
			return diag.FromErr(fmt.Errorf("a problem occured while accessing to the cluster"))
		}
		return diag.FromErr(fmt.Errorf("sorry, managing buckets is not available for hosted clusters"))
	}
	buckets, _, err := client.ClustersApi.ClustersListBuckets(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	bucketExists := false
	for _, bucket := range buckets {
		if bucket.Name == d.Id() {
			bucketExists = true
		}
	}
	if !bucketExists {
		bucketName := d.Id()
		d.SetId("")
		return diag.Errorf("Error 404: Failed to find the bucket %s ", bucketName)
	}
	return nil
}

/**
*** Updating the Bucket
**/
func resourceCouchbaseCapellaBucketUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)

	// Check if the Cluster is inVPC to update the bucket
	// Managing buckets is not available for hosted clusters
	_, _, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		// Check V3Cluster :: Need to be fixed in next versions
		_, _, err3 := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()
		if err3 != nil {
			return diag.FromErr(fmt.Errorf("a problem occured while accessing to the cluster"))
		}
		return diag.FromErr(fmt.Errorf("sorry, managing buckets is not available for hosted clusters"))
	}

	bucketName := d.Get("name").(string)
	memoryQuota := int32(d.Get("memory_quota").(int))

	// List buckets and iterate through to find bucket ID
	buckets, _, err := client.ClustersApi.ClustersListBuckets(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, bucket := range buckets {
		if bucket.Name == bucketName {
			bucketId := string(bucket.Id)

			updateBucketRequest := *couchbasecapella.NewUpdateBucketRequest(memoryQuota)

			// Update bucket with bucket ID
			_, err := client.ClustersApi.ClustersUpdateSingleBucket(auth, clusterId, bucketId).UpdateBucketRequest(updateBucketRequest).Execute()
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceCouchbaseCapellaBucketRead(ctx, d, meta)
}

/**
*** Deleting the Bucket
**/
func resourceCouchbaseCapellaBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("cluster_id").(string)

	// Check if the Cluster is inVPC to delete the bucket
	// Managing buckets is not available for hosted clusters
	_, _, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		// Check V3Cluster :: Need to be fixed in next versions
		_, _, err3 := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()
		if err3 != nil {
			return diag.FromErr(fmt.Errorf("a problem occured while accessing to the cluster"))
		}
		return diag.FromErr(fmt.Errorf("sorry, managing buckets is not available for hosted clusters"))
	}
	bucketName := d.Get("name").(string)

	deleteBucketRequest := *couchbasecapella.NewDeleteBucketRequest(bucketName)

	_, deleteError := client.ClustersApi.ClustersDeleteBucket(auth, clusterId).DeleteBucketRequest(deleteBucketRequest).Execute()
	if deleteError != nil {
		return diag.FromErr(deleteError)
	}
	return nil
}
