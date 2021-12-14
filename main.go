// Couchbase, Inc. licenses this to you under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0.

// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"go.uber.org/zap"

	"github.com/JamesWilkinsonCB/terraform-provider-couchbasecapella/internal/buckets"
	"github.com/JamesWilkinsonCB/terraform-provider-couchbasecapella/provider"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("unable to get logger: %v", err)
	}

	bGW, err := bucketGateway(logger)
	if err != nil {
		log.Fatalf("unable to get bucket gateway:%v", err)
	}

	p, err := provider.NewProvider(bGW)
	if err != nil {
		log.Fatalf("unable to get bucket resource: %v", err)
	}

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return p.Provider(
				p.WithCouchbaseCapellaBucketResource(),
			)
		},
	})
}

func bucketGateway(logger *zap.Logger) (*buckets.Gateway, error) {
	gw, err := buckets.NewGateway(logger)
	if err != nil {
		return nil, fmt.Errorf("unable to get bucket gateway: %w", err)
	}

	return gw, nil
}
