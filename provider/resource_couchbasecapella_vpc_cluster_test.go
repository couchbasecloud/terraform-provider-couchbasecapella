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

	couchbasecapella "github.com/couchbasecloud/couchbase-capella-api-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Test to see if a vpc cluster can be created, exists and is deleted successfully in AWS
func TestAccCouchbaseCapellaVpcCluster_AWS(t *testing.T) {
	var (
		cluster couchbasecapella.Cluster
	)

	resourceName := "couchbasecapella_vpc_cluster.test"
	clusterName := fmt.Sprintf("testacc-vpc-%s", acctest.RandString(5))
	cloudId := os.Getenv("CBC_AWS_CLOUD_ID")
	projectId := os.Getenv("CBC_PROJECT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaVpcClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaVpcClusterConfig_AWS(clusterName, cloudId, projectId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaVpcClusterExists(resourceName, &cluster),
				),
			},
		},
	})
}

// Test to see if a vpc cluster can be created, exists and is deleted successfully in Azure
func TestAccCouchbaseCapellaVpcCluster_Azure(t *testing.T) {
	var (
		cluster couchbasecapella.Cluster
	)

	resourceName := "couchbasecapella_vpc_cluster.test"
	clusterName := fmt.Sprintf("testacc-vpc-%s", acctest.RandString(5))
	cloudId := os.Getenv("CBC_AZURE_CLOUD_ID")
	projectId := os.Getenv("CBC_PROJECT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaVpcClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaVpcClusterConfig_Azure(clusterName, cloudId, projectId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaVpcClusterExists(resourceName, &cluster),
				),
			},
		},
	})
}

// Test to see if vpc cluster has been destroyed after Terraform Destory has been executed
func testAccCheckCouchbaseCapellaVpcClusterDestroy(s *terraform.State) error {
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
		if rs.Type != "couchbasecapella_vpc_cluster" {
			continue
		}

		_, _, err := client.ClustersApi.ClustersShow(auth, rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("vpc cluster (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// Test to see if vpc cluster exists after Terraform Apply has been executed
func testAccCheckCouchbaseCapellaVpcClusterExists(resourceName string, cluster *couchbasecapella.Cluster) resource.TestCheckFunc {
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
			return fmt.Errorf("no vpc cluster id is set")
		}

		_, _, err := client.ClustersApi.ClustersShow(auth, rs.Primary.ID).Execute()
		if err == nil {
			return nil
		}

		return fmt.Errorf("vpc cluster (%s) does not exist", rs.Primary.ID)
	}
}

// This is the Terraform Configuration that will be applied for the testing a cluster deployed in AWS
func testAccCouchbaseCapellaVpcClusterConfig_AWS(clusterName, cloudId, projectId string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_vpc_cluster" "test" {
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

// This is the Terraform Configuration that will be applied for the testing a cluster deployed in Azure
func testAccCouchbaseCapellaVpcClusterConfig_Azure(clusterName, cloudId, projectId string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_vpc_cluster" "test" {
			name       = "%s"
			cloud_id   = "%s"
			project_id = "%s"
			servers {
				size     = 3
				services = ["data", "query", "index"]
				azure {
					instance_size = "Standard_F4s_v2"
					volume_type  = "P4"
				}
			}
		}
	`, clusterName, cloudId, projectId)
}
