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

func TestAccCouchbaseCapellaCluster(t *testing.T) {
	var (
		cluster couchbasecapella.Cluster
	)

	clusterName := fmt.Sprintf("testacc-cluster-%s", acctest.RandString(5))
	cloudId := os.Getenv("CBC_CLOUD_ID")
	projectId := os.Getenv("CBC_PROJECT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaClusterConfig(clusterName, cloudId, projectId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaClusterExists("couchbasecapella_cluster.test", &cluster),
				),
			},
		},
	})
}

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
