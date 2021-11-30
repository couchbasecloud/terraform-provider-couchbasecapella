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

func TestAccCouchbaseCapellaDatabaseUser(t *testing.T) {
	var (
		databaseUser couchbasecapella.CreateDatabaseUserRequest
	)

	testClusterId := os.Getenv("CBC_CLUSTER_ID")
	username := fmt.Sprintf("testacc-user-%s", acctest.RandString(5))
	password := "Password123!"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaDatabaseUserConfig(testClusterId, username, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaDatabaseUserExists("couchbasecapella_database_user.test", &databaseUser),
				),
			},
		},
	})
}

func testAccCheckCouchbaseCapellaDatabaseUserDestroy(s *terraform.State) error {
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
		if rs.Type != "couchbasecapella_database_user" {
			continue
		}

		users, _, err := client.ClustersApi.ClustersListUsers(auth, rs.Primary.Attributes["cluster_id"]).Execute()
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		for _, user := range users {
			if user.Username == rs.Primary.ID {
				return fmt.Errorf("database user still exists")
			}
		}
	}

	return nil
}

func testAccCheckCouchbaseCapellaDatabaseUserExists(resourceName string, databaseUser *couchbasecapella.CreateDatabaseUserRequest) resource.TestCheckFunc {
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
			return fmt.Errorf("no username is set")
		}

		users, _, err := client.ClustersApi.ClustersListUsers(auth, rs.Primary.Attributes["cluster_id"]).Execute()
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		for _, user := range users {
			if user.Username == rs.Primary.Attributes["username"] {
				return nil
			}
		}

		return fmt.Errorf("database user does not exist")
	}
}

func testAccCouchbaseCapellaDatabaseUserConfig(clusterId, username, password string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_database_user" "test" {
			cluster_id   = "%s"
			username = "%s"
			password = "%s"
			all_bucket_access = "data_reader"
		}
	`, clusterId, username, password)
}
