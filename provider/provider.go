// Couchbase, Inc. licenses this to you under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0.

// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package provider

import (
	"context"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CBC_ACCESS_KEY", nil),
				Description: "Couchbase Capella API Access Key",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CBC_SECRET_KEY", nil),
				Description: "Couchbase Capella API Secret Key",
				Sensitive:   true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"couchbasecapella_project":       resourceCouchbaseCapellaProject(),
			"couchbasecapella_cluster":       resourceCouchbaseCapellaCluster(),
			"couchbasecapella_database_user": resourceCouchbaseCapellaDatabaseUser(),
			"couchbasecapella_bucket":        resourceCouchbaseCapellaBucket(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// TODO: create a client with access/secret keys
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	configuration := couchbasecapella.NewConfiguration()
	apiClient := couchbasecapella.NewAPIClient(configuration)
	return apiClient, nil
}
