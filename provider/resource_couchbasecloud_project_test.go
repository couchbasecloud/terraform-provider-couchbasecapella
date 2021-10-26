package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	couchbasecloud "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCouchbaseCloudProject_basic(t *testing.T) {
	var (
		project couchbasecloud.Project
	)

	projectName := fmt.Sprintf("testacc-project-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCloudProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCloudProjectConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCloudProjectExists("couchbasecloud_project.test", &project),
				),
			},
		},
	})
}

func testAccCheckCouchbaseCloudProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*couchbasecloud.APIClient)
	auth := context.WithValue(
		context.Background(),
		couchbasecloud.ContextAPIKeys,
		map[string]couchbasecloud.APIKey{
			"accessKey": {
				Key: os.Getenv("CBC_ACCESS_KEY"),
			},
			"secretKey": {
				Key: os.Getenv("CBC_SECRET_KEY"),
			},
		},
	)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "couchbasecloud_project" {
			continue
		}

		_, resp, _ := client.ProjectsApi.ProjectsShow(auth, rs.Primary.ID).Execute()
		if resp != nil {
			return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCouchbaseCloudProjectExists(resourceName string, project *couchbasecloud.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*couchbasecloud.APIClient)
		auth := context.WithValue(
			context.Background(),
			couchbasecloud.ContextAPIKeys,
			map[string]couchbasecloud.APIKey{
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

		log.Printf("[DEBUG] projectID: %s", rs.Primary.ID)

		_, _, err := client.ProjectsApi.ProjectsShow(auth, rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCouchbaseCloudProjectConfig(projectName string) string {
	return fmt.Sprintf(`
		resource "couchbasecloud_project" "test" {
			name   = "%s"
		}
	`, projectName)
}
