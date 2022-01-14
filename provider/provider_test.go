// Couchbase, Inc. licenses this to you under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0.

// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"couchbasecapella": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("CBC_ACCESS_KEY"); err == "" {
		t.Fatal("CBC_ACCESS_KEY must be set for acceptance tests")
	}
	if err := os.Getenv("CBC_SECRET_KEY"); err == "" {
		t.Fatal("CBC_SECRET_KEY must be set for acceptance tests")
	}
	if err := os.Getenv("CBC_AWS_CLOUD_ID"); err == "" {
		t.Fatal("CBC_AWS_CLOUD_ID must be set for acceptance tests")
	}
	if err := os.Getenv("CBC_AZURE_CLOUD_ID"); err == "" {
		t.Fatal("CBC_AZURE_CLOUD_ID must be set for acceptance tests")
	}
	if err := os.Getenv("CBC_PROJECT_ID"); err == "" {
		t.Fatal("CBC_PROJECT_ID must be set for acceptance tests")
	}
	if err := os.Getenv("CBC_CLUSTER_ID"); err == "" {
		t.Fatal("CBC_CLUSTER_ID must be set for acceptance tests")
	}
	if err := os.Getenv("CBC_CLUSTER_CIDR"); err == "" {
		t.Fatal("CBC_CLUSTER_CIDR must be set for acceptance tests")
	}
	if err := os.Getenv("CBC_BUCKET_NAME"); err == "" {
		t.Fatal("CBC_BUCKET_NAME must be set for acceptance tests")
	}
}
