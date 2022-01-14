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
	"os"
	"testing"
	"time"

	couchbasecapella "github.com/couchbasecloud/couchbase-capella-api-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Test to see if a bucket with sequential number conflict resolution can be created, updated and deleted
// successfully
func TestAccCouchbaseCapellaBucket_withSequentialNumberResolution(t *testing.T) {
	var (
		bucket couchbasecapella.CouchbaseBucketSpec
	)

	testClusterId := os.Getenv("CBC_CLUSTER_ID")
	bucketName := fmt.Sprintf("testacc-bucket-%s", acctest.RandString(5))
	memoryQuota := "128"
	updatedMemoryQuota := "256"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaBucketConfig_withSequentialNumberResolution(testClusterId, bucketName, memoryQuota),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaBucketExists("couchbasecapella_bucket.test", &bucket),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "memory_quota", memoryQuota),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "replicas", "1"),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "conflict_resolution", "seqno"),
				),
			},
			{
				Config: testAccCouchbaseCapellaBucketConfig_withSequentialNumberResolution(testClusterId, bucketName, updatedMemoryQuota),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaBucketExists("couchbasecapella_bucket.test", &bucket),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "memory_quota", updatedMemoryQuota),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "replicas", "1"),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "conflict_resolution", "seqno"),
				),
			},
		},
	})
}

// Test to see if a bucket with last write wins conflict resolution can be created, updated and deleted
// successfully
func TestAccCouchbaseCapellaBucket_withLastWriteWinsResolution(t *testing.T) {
	var (
		bucket couchbasecapella.CouchbaseBucketSpec
	)

	testClusterId := os.Getenv("CBC_CLUSTER_ID")
	bucketName := fmt.Sprintf("testacc-bucket-%s", acctest.RandString(5))
	memoryQuota := "128"
	updatedMemoryQuota := "256"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaBucketConfig_withLastWriteWinsResolution(testClusterId, bucketName, memoryQuota),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaBucketExists("couchbasecapella_bucket.test", &bucket),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "memory_quota", memoryQuota),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "replicas", "1"),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "conflict_resolution", "lww"),
				),
			},
			{
				Config: testAccCouchbaseCapellaBucketConfig_withLastWriteWinsResolution(testClusterId, bucketName, updatedMemoryQuota),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaBucketExists("couchbasecapella_bucket.test", &bucket),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "name", bucketName),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "memory_quota", updatedMemoryQuota),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "replicas", "1"),
					resource.TestCheckResourceAttr("couchbasecapella_bucket.test", "conflict_resolution", "lww"),
				),
			},
		},
	})
}

// Test to see if bucket has been destroyed after Terraform Destory has been executed
func testAccCheckCouchbaseCapellaBucketDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*couchbasecapella.APIClient)
	auth := context.WithValue(
		context.Background(),
		couchbasecapella.ContextAPIKeys,
		map[string]couchbasecapella.APIKey{
			"accessKey": {
				Key: os.Getenv("CBC_ACCESS_KEY"),
			},
			"secretKey": {
				Key: os.Getenv("CBC_SECRET_KEY"),
			},
		},
	)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchbasecapella_bucket" {
			continue
		}

		buckets, _, err := client.ClustersApi.ClustersListBuckets(auth, rs.Primary.Attributes["cluster_id"]).Execute()
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		for _, bucket := range buckets {
			if bucket.Name == rs.Primary.Attributes["name"] {
				return fmt.Errorf("bucket still exists")
			}
		}
	}

	return nil
}

// Test to see if bucket exists after Terraform Apply has been executed
func testAccCheckCouchbaseCapellaBucketExists(resourceName string, bucket *couchbasecapella.CouchbaseBucketSpec) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*couchbasecapella.APIClient)
		auth := context.WithValue(
			context.Background(),
			couchbasecapella.ContextAPIKeys,
			map[string]couchbasecapella.APIKey{
				"accessKey": {
					Key: os.Getenv("CBC_ACCESS_KEY"),
				},
				"secretKey": {
					Key: os.Getenv("CBC_SECRET_KEY"),
				},
			},
		)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no bucket id is set")
		}

		// NOTE: There is a delay for retrieving a newly created bucket from Capella's list of buckets.
		// This poll will check at regular intervals if the newly created bucket is in the list of buckets
		// until the timeout period expires. If the timeout is reached, then the newly created bucket is not in the list of buckets.
		timeout := time.NewTimer(time.Second * 60)
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-timeout.C:
				return fmt.Errorf("bucket does not exist")
			case <-ticker.C:
				buckets, _, err := client.ClustersApi.ClustersListBuckets(auth, rs.Primary.Attributes["cluster_id"]).Execute()
				if err != nil {
					return fmt.Errorf("%s", err)
				}
				for _, bucket := range buckets {
					if bucket.Name == rs.Primary.ID {
						return nil
					}
				}
			}
		}
	}
}

// This is the Terraform Configuration that will be applied for testing a bucket with sequential number
// conflict resolution
func testAccCouchbaseCapellaBucketConfig_withSequentialNumberResolution(clusterId, bucketName, memoryQuota string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_bucket" "test" {
			cluster_id = "%s"
			name   = "%s"
			memory_quota = "%s"
			replicas = "1"
			conflict_resolution = "seqno"
		}
	`, clusterId, bucketName, memoryQuota)
}

// This is the Terraform Configuration that will be applied for testing a bucket with last write wins
// conflict resolution
func testAccCouchbaseCapellaBucketConfig_withLastWriteWinsResolution(clusterId, bucketName, memoryQuota string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_bucket" "test" {
			cluster_id = "%s"
			name   = "%s"
			memory_quota = "%s"
			replicas = "1"
			conflict_resolution = "lww"
		}
	`, clusterId, bucketName, memoryQuota)
}
