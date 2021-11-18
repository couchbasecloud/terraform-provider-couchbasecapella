package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCouchbaseCapellaBucket_basic(t *testing.T) {
	var (
		bucket couchbasecapella.CouchbaseBucketSpec
	)

	testClusterId := os.Getenv("CBC_CLUSTER_ID")
	bucketName := fmt.Sprintf("testacc-bucket-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaBucketDestroy,
		Steps: []resource.TestStep{
			{
				// TODO: Seems to be an issue with bucket config - test bucket is not created - or is not returned when buckets listed
				Config: testAccCouchbaseCapellaBucketConfig(testClusterId, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaBucketExists("couchbasecapella_bucket.test", &bucket),
				),
			},
		},
	})
}

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

		log.Printf("[DEBUG] bucketID: %s", rs.Primary.ID)

		buckets, _, err := client.ClustersApi.ClustersListBuckets(auth, rs.Primary.Attributes["cluster_id"]).Execute()
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		for _, bucket := range buckets {
			if bucket.Name == rs.Primary.Attributes["name"] {
				return nil
			}
		}
		return fmt.Errorf("bucket does not exist")
	}
}

func testAccCouchbaseCapellaBucketConfig(clusterId, bucketName string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_bucket" "test" {
			cluster_id = "%s"
			name   = "%s"
			memory_quota = "128"
			replicas = "1"
			conflict_resolution = "seqno"
		}
	`, clusterId, bucketName)
}
