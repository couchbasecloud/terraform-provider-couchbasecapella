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

func TestAccCouchbaseCapellaProject(t *testing.T) {
	var (
		project couchbasecapella.Project
	)

	projectName := fmt.Sprintf("testacc-project-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaProjectConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaProjectExists("couchbasecapella_project.test", &project),
				),
			},
		},
	})
}

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

func testAccCouchbaseCapellaProjectConfig(projectName string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_project" "test" {
			name   = "%s"
		}
	`, projectName)
}
