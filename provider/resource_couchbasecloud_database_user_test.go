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

func TestAccCouchbaseCloudDatabaseUser_basic(t *testing.T) {
	var (
		databaseUser couchbasecloud.CreateDatabaseUserRequest
	)

	testClusterId := ""
	username := fmt.Sprintf("testacc-user-%s", acctest.RandString(10))
	password := fmt.Sprintf("%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCloudDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCloudDatabaseUserConfig(testClusterId, username, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCloudDatabaseUserExists("couchbasecloud_database_user.test", &databaseUser),
				),
			},
		},
	})
}

func testAccCheckCouchbaseCloudDatabaseUserDestroy(s *terraform.State) error {
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
		if rs.Type != "couchbasecloud_database_user" {
			continue
		}

		_, resp, _ := client.ProjectsApi.ProjectsShow(auth, rs.Primary.ID).Execute()
		if resp != nil {
			return fmt.Errorf("database user (%s) still exists", rs.Primary.ID)
		}

		users, _, err := client.ClustersApi.ClustersListUsers(auth, rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("database user still exists")
		}
		for _, user := range users {
			if user.Username == rs.Primary.Attributes["username"] {
				return nil
			}
		}
	}

	return nil
}

func testAccCheckCouchbaseCloudDatabaseUserExists(resourceName string, databaseUser *couchbasecloud.CreateDatabaseUserRequest) resource.TestCheckFunc {
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
			return fmt.Errorf("no username is set")
		}

		log.Printf("[DEBUG] Database User: %s", rs.Primary.ID)

		users, _, err := client.ClustersApi.ClustersListUsers(auth, rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("database user still exists")
		}
		for _, user := range users {
			if user.Username == rs.Primary.Attributes["username"] {
				return nil
			}
		}

		return nil
	}
}

func testAccCouchbaseCloudDatabaseUserConfig(clusterId, username, password string) string {
	return fmt.Sprintf(`
		resource "couchbasecloud_database_user" "test" {
			cluster_id   = "%s"
			username = "%s"
			password = "%s"
			all_bucket_access = "data_reader"
		}
	`, clusterId, username, password)
}
