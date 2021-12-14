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
	"strings"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/JamesWilkinsonCB/terraform-provider-couchbasecapella/internal/buckets"
)

// Provider is responsible for configuring the TF resources and provider
// NOTE: could do with a much better comment here
type Provider struct {
	bucketGW *buckets.Gateway
}

// NewProvider returns an instantiated instance of a provider. A provides
// manages the resource handlers for the TF provider. It has the following
// dependencies:
//
// bucket gateway - for managing bucket resources
func NewProvider(bucketGW *buckets.Gateway) (*Provider, error) {
	p := Provider{
		bucketGW: bucketGW,
	}

	if err := p.validate(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Provider) validate() error {
	var missingDeps []string

	for _, tc := range []struct {
		dep string
		chk func() bool
	}{
		{
			dep: "bucket gateway",
			chk: func() bool { return p.bucketGW != nil },
		},
	} {
		if !tc.chk() {
			missingDeps = append(missingDeps, tc.dep)
		}
	}

	if len(missingDeps) > 0 {
		return fmt.Errorf(
			"unable to initailize a provider due to (%d) missing dependencies: %s",
			len(missingDeps),
			strings.Join(missingDeps, ","),
		)
	}

	return nil
}

// Provider returns a Terraform provider with the resources options provided
// TODO: better comment here
// TODO: figure out what to do with data resources
func (p *Provider) Provider(resources ...ResourceOption) *schema.Provider {
	tfProvider := &schema.Provider{
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
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}

	for i := range resources {
		resources[i](tfProvider)
	}

	return tfProvider
}

type ResourceOption func(p *schema.Provider)

// WithCouchbaseCapellaBucketResource returns an option that sets the underlying
// resources map to the Couchbase Capella bucket resource handled by the bucket
// gateway
func (p *Provider) WithCouchbaseCapellaBucketResource() ResourceOption {
	return func(tfProvider *schema.Provider) {
		tfProvider.ResourcesMap[buckets.Resource] = buckets.NewCouchbaseCapellaBucketResource(p.bucketGW)
	}
}

// TODO: create a client with access/secret keys
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	configuration := couchbasecapella.NewConfiguration()
	apiClient := couchbasecapella.NewAPIClient(configuration)
	return apiClient, nil
}

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}
