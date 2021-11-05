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

func TestAccCouchbaseCapellaCluster_basic(t *testing.T) {
	var (
		cluster couchbasecapella.Cluster
	)

	projectName := fmt.Sprintf("testacc-cluster-%s", acctest.RandString(10))
	cloudId := "5da1481d-5884-486a-ba24-a7b6410c9637"
	projectId := "8b0d5bb1-1ea3-4474-b0ad-3775fe9bfe2e"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaClusterConfig(projectName, cloudId, projectId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaClusterExists("couchbasecapella_cluster.test", &cluster),
				),
			},
		},
	})
}

// This currently doesn't work as the cluster must have `Healthy` status before it can be destroyed
// It cannot be destroyed whilst it is still `Deploying`
func testAccCheckCouchbaseCapellaClusterDestroy(s *terraform.State) error {
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
		if rs.Type != "couchbasecapella_cluster" {
			continue
		}

		_, _, err := client.ClustersApi.ClustersShow(auth, rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("cluster (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCouchbaseCapellaClusterExists(resourceName string, cluster *couchbasecapella.Cluster) resource.TestCheckFunc {
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
			return fmt.Errorf("no cluster id is set")
		}

		log.Printf("[DEBUG] clusterID: %s", rs.Primary.ID)

		_, _, err := client.ClustersApi.ClustersShow(auth, rs.Primary.ID).Execute()
		if err == nil {
			return nil
		}

		return fmt.Errorf("cluster (%s) does not exist", rs.Primary.ID)
	}
}

func testAccCouchbaseCapellaClusterConfig(clusterName, cloudId, projectId string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_cluster" "test" {
			name   = "%s"
			cloud_id = "%s"
			project_id = "%s"
			servers{
				size = 3
				services = ["data"]
				aws {
					instance_size = "m5.xlarge"
					ebs_size_gib = 50
				}
			}
		}
	`, clusterName, cloudId, projectId)
}
