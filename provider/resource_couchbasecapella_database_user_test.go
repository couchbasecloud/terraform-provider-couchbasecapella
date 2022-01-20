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

// Test to see if a databaser user with all bucket access can be created, updated, recreated and deleted
// successfully
func TestAccCouchbaseCapellaDatabaseUser_allBucketAccess(t *testing.T) {
	var (
		databaseUser couchbasecapella.CreateDatabaseUserRequest
	)

	testClusterId := os.Getenv("CBC_CLUSTER_ID")
	resourceName := "couchbasecapella_database_user.test"
	username := fmt.Sprintf("testacc-user-%s", acctest.RandString(5))
	updatedUsername := fmt.Sprintf("testacc-user-%s", acctest.RandString(4))
	password := "Password123!"
	allBucketAccess := "data_reader"
	updatedAllBucketAccess := "data_writer"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaDatabaseUserConfig_allBucketAccess(testClusterId, username, password, allBucketAccess),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaDatabaseUserExists(resourceName, &databaseUser),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "all_bucket_access", allBucketAccess),
				),
			},
			{
				Config: testAccCouchbaseCapellaDatabaseUserConfig_allBucketAccess(testClusterId, username, password, updatedAllBucketAccess),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaDatabaseUserExists(resourceName, &databaseUser),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "all_bucket_access", updatedAllBucketAccess),
				),
			},
			{
				Config: testAccCouchbaseCapellaDatabaseUserConfig_allBucketAccess(testClusterId, updatedUsername, password, allBucketAccess),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaDatabaseUserExists(resourceName, &databaseUser),
					resource.TestCheckResourceAttr(resourceName, "username", updatedUsername),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "all_bucket_access", allBucketAccess),
				),
			},
		},
	})
}

// Test to see if a databaser user with specific bucket access can be created, updated, recreated and deleted
// successfully
func TestAccCouchbaseCapellaDatabaseUser_specificBucketAccess(t *testing.T) {
	var (
		databaseUser couchbasecapella.CreateDatabaseUserRequest
	)

	testClusterId := os.Getenv("CBC_CLUSTER_ID")

	resourceName := "couchbasecapella_database_user.test"
	username := fmt.Sprintf("testacc-user-%s", acctest.RandString(5))
	updatedUsername := fmt.Sprintf("testacc-user-%s", acctest.RandString(4))
	password := "Password123!"
	testBucketName := os.Getenv("CBC_BUCKET_NAME")
	bucketAccess := "data_reader"
	updatedBucketAccess := "data_writer"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCouchbaseCapellaDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCouchbaseCapellaDatabaseUserConfig_specificBucketAccess(testClusterId, username, password, testBucketName, bucketAccess),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaDatabaseUserExists(resourceName, &databaseUser),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
				),
			},
			{
				Config: testAccCouchbaseCapellaDatabaseUserConfig_specificBucketAccess(testClusterId, username, password, testBucketName, updatedBucketAccess),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaDatabaseUserExists(resourceName, &databaseUser),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
				),
			},
			{
				Config: testAccCouchbaseCapellaDatabaseUserConfig_specificBucketAccess(testClusterId, updatedUsername, password, testBucketName, updatedBucketAccess),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCouchbaseCapellaDatabaseUserExists(resourceName, &databaseUser),
					resource.TestCheckResourceAttr(resourceName, "username", updatedUsername),
					resource.TestCheckResourceAttr(resourceName, "password", password),
				),
			},
		},
	})
}

// Test to see if database user has been destroyed after Terraform Destory has been executed
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

// Test to see if database user exists after Terraform Apply has been executed
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

// This is the Terraform Configuration that will be applied for testing a database user with all bucket access
func testAccCouchbaseCapellaDatabaseUserConfig_allBucketAccess(clusterId, username, password, allBucketAccess string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_database_user" "test" {
			cluster_id   = "%s"
			username = "%s"
			password = "%s"
			all_bucket_access = "%s"
		}
	`, clusterId, username, password, allBucketAccess)
}

// This is the Terraform Configuration that will be applied for testing a database user with specific bucket access
func testAccCouchbaseCapellaDatabaseUserConfig_specificBucketAccess(clusterId, username, password, bucketName, bucketAccess string) string {
	return fmt.Sprintf(`
		resource "couchbasecapella_database_user" "test" {
			cluster_id = "%s"
			username = "%s"
			password = "%s"
			buckets{
				bucket_name = "%s"
				bucket_access = ["%s"]
			}
		}
	`, clusterId, username, password, bucketName, bucketAccess)
}
