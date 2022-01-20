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

// Test to see if a project can be created, recreated and deleted successfully
func TestAccCouchbaseCapellaProject(t *testing.T) {
	var (
		project couchbasecapella.Project
	)

	projectName := fmt.Sprintf("testacc-project-%s", acctest.RandString(5))
	updateProjectName := fmt.Sprintf("testacc-project-%s", acctest.RandString(4))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaProjectConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaProjectExists("couchbasecapella_project.test", &project),
					resource.TestCheckResourceAttr("couchbasecapella_project.test", "name", projectName),
				),
			},
			{
				Config: testAccCouchbaseCapellaProjectConfig(updateProjectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaProjectExists("couchbasecapella_project.test", &project),
					resource.TestCheckResourceAttr("couchbasecapella_project.test", "name", updateProjectName),
				),
			},
		},
	})
}

// Test to see if project has been destroyed after Terraform Destory has been executed
func testAccCheckCouchbaseCapellaProjectDestroy(s *terraform.State) error {
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
		if rs.Type != "couchbasecapella_project" {
			continue
		}

		_, _, err := client.ProjectsApi.ProjectsShow(auth, rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// Test to see if project exists after Terraform Apply has been executed
func testAccCheckCouchbaseCapellaProjectExists(resourceName string, project *couchbasecapella.Project) resource.TestCheckFunc {
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
			return fmt.Errorf("no project name is set")
		}

		_, _, err := client.ProjectsApi.ProjectsShow(auth, rs.Primary.ID).Execute()
		if err == nil {
			return nil
		}

		return fmt.Errorf("project (%s) does not exist", rs.Primary.ID)
	}
}

// This is the Terraform Configuration that will be applied for the tests
func testAccCouchbaseCapellaProjectConfig(projectName string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_project" "test" {
			name   = "%s"
		}
	`, projectName)
}
