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

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Test to see if a hosted cluster is created, updated and deleted successfully.
// Once the cluster has been deployed, there will be tests to check that
// the name, support packages and server services have updated correctly.
func TestAccCouchbaseCapellaHostedCluster(t *testing.T) {
	var (
		cluster couchbasecapella.V3Cluster
	)

	resourceName := "couchbasecapella_hosted_cluster.test"
	clusterName := fmt.Sprintf("testacc-hosted-%s", acctest.RandString(5))
	updatedClusterName := fmt.Sprintf("testacc-hosted-%s", acctest.RandString(4))
	projectId := os.Getenv("CBC_PROJECT_ID")
	supportPackageType := "Basic"
	updatedSupportPackageType := "DeveloperPro"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaHostedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaHostedClusterConfig(clusterName, projectId, supportPackageType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaHostedClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "support_package.0.support_package_type", supportPackageType),
				),
			},
			{
				Config: testAccCouchbaseCapellaHostedClusterConfig(updatedClusterName, projectId, supportPackageType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaHostedClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", updatedClusterName),
					resource.TestCheckResourceAttr(resourceName, "support_package.0.support_package_type", supportPackageType),
				),
			},
			{
				Config: testAccCouchbaseCapellaHostedClusterConfig(updatedClusterName, projectId, updatedSupportPackageType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaHostedClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", updatedClusterName),
					resource.TestCheckResourceAttr(resourceName, "support_package.0.support_package_type", updatedSupportPackageType),
				),
			},
			{
				Config: testAccCouchbaseCapellaHostedClusterConfig_withUpdatedServices(updatedClusterName, projectId, updatedSupportPackageType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaHostedClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", updatedClusterName),
					resource.TestCheckResourceAttr(resourceName, "support_package.0.support_package_type", updatedSupportPackageType),
				),
			},
			{
				Config: testAccCouchbaseCapellaHostedClusterConfig(updatedClusterName, projectId, updatedSupportPackageType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaHostedClusterExists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", updatedClusterName),
					resource.TestCheckResourceAttr(resourceName, "support_package.0.support_package_type", updatedSupportPackageType),
				),
			},
		},
	})
}

// Test to see if hosted cluster has been destroyed after Terraform Destory has been executed
func testAccCheckCouchbaseCapellaHostedClusterDestroy(s *terraform.State) error {
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
		if rs.Type != "couchbasecapella_hosted_cluster" {
			continue
		}

		_, _, err := client.ClustersV3Api.ClustersV3show(auth, rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("hosted cluster (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// Test to see if hosted cluster exists after Terraform Apply has been executed
func testAccCheckCouchbaseCapellaHostedClusterExists(resourceName string, cluster *couchbasecapella.V3Cluster) resource.TestCheckFunc {
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

		_, _, err := client.ClustersV3Api.ClustersV3show(auth, rs.Primary.ID).Execute()
		if err == nil {
			return nil
		}

		return fmt.Errorf("hosted cluster (%s) does not exist", rs.Primary.ID)
	}
}

// This is the Terraform Configuration that will be applied for testing a hosted cluster can be created and updated
func testAccCouchbaseCapellaHostedClusterConfig(clusterName, projectId, supportPackageType string) string {
	return fmt.Sprintf(`
	resource "couchbasecapella_hosted_cluster" "test" {
		name        = "%s"
		project_id  = "%s"
		place {
			single_az = true
			hosted {
				provider = "aws"
				region   = "us-west-2"
				cidr     = "10.0.16.0/20"
			}
		}
		support_package {
			timezone = "GMT"
			support_package_type     = "%s"
		}
		servers {
			size     = 3
			compute  = "m5.xlarge"
			services = ["data"]
			storage {
				storage_type = "GP3"
				iops = "3000"
				storage_size = "50"
			}
		}
	}	
	`, clusterName, projectId, supportPackageType)
}

// This is the Terraform Configuration that will be applied for testing updating a hosted cluster services
func testAccCouchbaseCapellaHostedClusterConfig_withUpdatedServices(clusterName, projectId, supportPackageType string) string {
	return fmt.Sprintf(`
	resource "couchbasecapella_hosted_cluster" "test" {
		name        = "%s"
		project_id  = "%s"
		place {
			single_az = true
			hosted {
				provider = "aws"
				region   = "us-west-2"
				cidr     = "10.0.16.0/20"
			}
		}
		support_package {
			timezone = "GMT"
			support_package_type     = "%s"
		}
		servers {
			size     = 3
			compute  = "m5.xlarge"
			services = ["data", "index", "query"]
			storage {
				storage_type = "GP3"
				iops = "3000"
				storage_size = "50"
			}
		}
	}	
	`, clusterName, projectId, supportPackageType)
}
