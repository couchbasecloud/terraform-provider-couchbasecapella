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

func TestAccCouchbaseCapellaHostedCluster(t *testing.T) {
	var (
		cluster couchbasecapella.V3Cluster
	)

	clusterName := fmt.Sprintf("testacc-hosted-%s", acctest.RandString(5))
	projectId := os.Getenv("CBC_PROJECT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaHostedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaHostedClusterConfig(clusterName, projectId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaHostedClusterExists("couchbasecapella_hosted_cluster.test", &cluster),
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

// This is the Terraform Configuration that will be applied for the tests
func testAccCouchbaseCapellaHostedClusterConfig(clusterName, projectId string) string {
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
			type     = "Basic"
		}
		servers {
			size     = 3
			compute  = "m5.xlarge"
			services = ["data"]
			storage {
				type = "GP3"
				iops = "3000"
				size = "50"
			}
		}
	}	
	`, clusterName, projectId)
}
