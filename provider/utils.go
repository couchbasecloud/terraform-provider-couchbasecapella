// Couchbase, Inc. licenses this to you under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0.

// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package provider

import (
	"context"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	couchbasecapella "github.com/couchbasecloud/couchbase-capella-api-go-client"
)

func Has(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// TODO: Get Auth with both env variable and terraform ones
func getAuth(ctx context.Context) context.Context {
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
	return auth
}

func manageErrors(err error, r http.Response, functionality string) diag.Diagnostics {
	if err != nil {
		switch r.StatusCode {
		case 403:
			return diag.Errorf("You don't have the required access to apply this function " + functionality)
		case 401:
			return diag.Errorf("Please verify the validity of your Access key and Secret key")
		case 422:
			body, _ := io.ReadAll(r.Body)
			return diag.Errorf("Failed to create resource; API response: \n%s", body)
		default:
			return diag.FromErr(err)
		}
	}
	return nil
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
