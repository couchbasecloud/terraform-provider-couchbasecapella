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

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Description:  "ID of the Cluster",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Name of the Bucket",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateBucketName,
			},
			"memory_quota": {
				Description:  "Bucket Memory quota in Mb",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateMemoryQuota,
			},
			"conflict_resolution": {
				Description:  "Conflict resolution for bucket",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateConflictResolution,
			},
		},
	}
}

// resourceCouchbaseCapellaBucketCreate is responsible for creating a
// bucket in a Couchbase Capella VPC Cluster using the Terraform resource data.
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
			return diag.FromErr(fmt.Errorf(ClusterProblemAccessing))
		}
		return diag.FromErr(fmt.Errorf(BucketHostedNotSupported))
	}

	bucketName := d.Get("name").(string)
	memoryQuota := int32(d.Get("memory_quota").(int))
	conflictResolution := couchbasecapella.ConflictResolution(d.Get("conflict_resolution").(string))

	couchbaseBucketSpec := couchbasecapella.NewCouchbaseBucketSpec(bucketName, memoryQuota)
	couchbaseBucketSpec.SetConflictResolution(conflictResolution)

	_, r, err := client.ClustersApi.ClustersCreateBucket(auth, clusterId).CouchbaseBucketSpec(*couchbaseBucketSpec).Execute()
	if r == nil {
		return diag.Errorf("Pointer to bucket create http.Response is nil")
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("problem ocurred :: %v", r))
	}

	d.SetId(bucketName)

	return resourceCouchbaseCapellaBucketRead(ctx, d, meta)
}

// resourceCouchbaseCapellaBucketRead is responsible for reading a
// bucket in a Couchbase Capella VPC Cluster using the Terraform resource data.
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
			return diag.FromErr(fmt.Errorf(ClusterProblemAccessing))
		}
		return diag.FromErr(fmt.Errorf(BucketHostedNotSupported))
	}

	// NOTE: There is a delay for retrieving a newly created bucket from Capella's list of buckets.
	// This poll will check at regular intervals if the newly created bucket is in the list of buckets
	// until the timeout period expires. If the timeout is reached, then the newly created bucket is not in the list of buckets.
	timeout := time.NewTimer(time.Second * 180)
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-timeout.C:
			bucketName := d.Id()
			d.SetId("")
			return diag.Errorf("Error 404: Failed to find the bucket %s ", bucketName)
		case <-ticker.C:
			buckets, _, err := client.ClustersApi.ClustersListBuckets(auth, clusterId).Execute()
			if err != nil {
				return diag.FromErr(err)
			}
			for _, bucket := range buckets {
				if bucket.Name == d.Id() {
					return nil
				}
			}
		}
	}
}

// resourceCouchbaseCapellaBucketUpdate is responsible for updating a
// bucket in a Couchbase Capella VPC Cluster using the Terraform resource data.
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
			return diag.FromErr(fmt.Errorf(ClusterProblemAccessing))
		}
		return diag.FromErr(fmt.Errorf(BucketHostedNotSupported))
	}

	bucketName := d.Get("name").(string)

	// List buckets and iterate through to find bucket ID
	buckets, _, err := client.ClustersApi.ClustersListBuckets(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	for _, bucket := range buckets {
		if bucket.Name == bucketName {
			d.Set("conflict_resolution", bucket.ConflictResolution)
			d.Set("memory_quota", bucket.MemoryQuota)
		}
	}
	return diag.FromErr(fmt.Errorf(BucketHostedNotSupported))
	// return resourceCouchbaseCapellaBucketRead(ctx, d, meta)
}

// resourceCouchbaseCapellaBucketDelete is responsible for deleting a
// bucket in a Couchbase Capella VPC Cluster using the Terraform resource data.
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
			return diag.FromErr(fmt.Errorf(ClusterProblemAccessing))
		}
		return diag.FromErr(fmt.Errorf(BucketHostedNotSupported))
	}
	bucketName := d.Get("name").(string)

	deleteBucketRequest := *couchbasecapella.NewDeleteBucketRequest(bucketName)

	_, deleteError := client.ClustersApi.ClustersDeleteBucket(auth, clusterId).DeleteBucketRequest(deleteBucketRequest).Execute()
	if deleteError != nil {
		return diag.FromErr(deleteError)
	}
	return nil
}
