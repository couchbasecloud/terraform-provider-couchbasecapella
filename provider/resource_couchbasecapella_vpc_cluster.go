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
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	couchbasecapella "github.com/couchbasecloud/couchbase-capella-api-go-client"
)

func resourceCouchbaseCapellaVpcCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Capella vpc clusters.",

		CreateContext: resourceCouchbaseCapellaVpcClusterCreate,
		ReadContext:   resourceCouchbaseCapellaVpcClusterRead,
		DeleteContext: resourceCouchbaseCapellaVpcClusterDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of the Cluster",
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
			},
			"name": {
				Description:  "Name of the Cluster",
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateClusterName,
			},
			"cloud_id": {
				Description:  "ID of the Cloud the Cluster will be deployed in",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"project_id": {
				Description:  "ID of the Project",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"servers": {
				Description: "Server Configuration of the Cluster",
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:         schema.TypeInt,
							Description:  "Number of nodes",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateSize,
						},
						"services": {
							Type:        schema.TypeList,
							Description: "Couchbase Services",
							Required:    true,
							ForceNew:    true,
							MinItems:    1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateService,
							},
						},
						"aws": {
							Description: "Aws configuration",
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_size": {
										Description:  "Aws instance",
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validateAwsInstance,
									},
									"ebs_size_gib": {
										Description:  "Aws volume size (Gb)",
										Type:         schema.TypeInt,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validateAwsVolumeSize,
									},
								},
							},
						},
						"azure": {
							Description: "Azure configuration",
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_size": {
										Description:  "Azure instance",
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validateAzureInstance,
									},
									"volume_type": {
										Description:  "Azure volume size (Gb)",
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validateAzureVolume,
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(25 * time.Minute),
			Delete: schema.DefaultTimeout(25 * time.Minute),
		},
	}
}

// resourceCouchbaseCapellaVpcClusterCreate is responsible for creating a
// vpc cluster in Couchbase Capella using the Terraform resource data.
func resourceCouchbaseCapellaVpcClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterName := d.Get("name").(string)
	cloudId := d.Get("cloud_id").(string)
	projectId := d.Get("project_id").(string)

	newClusterRequest := *couchbasecapella.NewCreateClusterRequest(clusterName, cloudId, projectId)

	// Get The cloud
	cloud, resp, err := client.CloudsApi.CloudsShow(auth, cloudId).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return diag.FromErr(fmt.Errorf("404: the cloud doesn't exist. Please verify your cloud_id"))
		}
		return diag.FromErr(err)
	}
	providerName := string(cloud.Provider)
	// add Servers + Check servers Vs Cloud provider
	if servers, ok := d.GetOk("servers"); ok {
		// check server providers
		providers := getVpcServersProvider(servers.(*schema.Set))
		if len(providers) > 1 {
			return diag.FromErr(fmt.Errorf(VpcClusterServerDoesNotMatchProvider))
		}
		if len(providers) == 1 && !Has(providers, providerName) {
			return diag.FromErr(fmt.Errorf(VpcClusterServerDoesNotMatchProvider))
		}
		newClusterRequest.SetServers(expandVpcServersSet(servers.(*schema.Set)))
	}

	// Create the cluster
	response, err := client.ClustersApi.ClustersCreate(auth).CreateClusterRequest(newClusterRequest).Execute()
	if err != nil {
		return manageErrors(err, *response, "Create Cluster")
	}

	// TODO: need to be changed after cloud api fix!
	location := string(response.Header.Get("Location"))
	urlparts := strings.Split(location, "/")
	clusterId := urlparts[len(urlparts)-1]
	d.SetId(clusterId)

	defer response.Body.Close()

	// Wait for the cluster to deploy
	createStateConf := &resource.StateChangeConf{
		Pending: []string{"deploying", "deploy_succeeded"},
		Target:  []string{"ready"},
		Refresh: func() (interface{}, string, error) {
			statusResp, _, err := client.ClustersApi.ClustersStatus(auth, clusterId).Execute()
			if err != nil {
				return 0, "Error", err
			}
			return statusResp, string(statusResp.Status), nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Minute,
		MinTimeout: 30 * time.Second,
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vpc cluster (%s) to be created: %s", d.Id(), err)
	}

	return resourceCouchbaseCapellaVpcClusterRead(ctx, d, meta)
}

// resourceCouchbaseCapellaVpcClusterRead is responsible for reading a
// vpc cluster in Couchbase Capella using the Terraform resource data.
func resourceCouchbaseCapellaVpcClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)
	clusterId := d.Id()

	_, resp, err := client.ClustersApi.ClustersShow(auth, clusterId).Execute()

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

// resourceCouchbaseCapellaVpcClusterDelete is responsible for deleting a
// vpc cluster in Couchbase Capella using the Terraform resource data.
func resourceCouchbaseCapellaVpcClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Id()

	// Check that Cluster is ready to be destroyed
	statusResp, _, err := client.ClustersApi.ClustersStatus(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if statusResp.Status != couchbasecapella.CLUSTERSTATUS_READY {
		return diag.Errorf("VPC Cluster is not ready to be deleted. Cluster Status: %s", statusResp.Status)
	}

	r, err2 := client.ClustersApi.ClustersDelete(auth, clusterId).Execute()
	if err2 != nil {
		return manageErrors(err2, *r, "VPC Cluster Delete")
	}

	// Wait for the cluster to be destroyed
	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{"destroying", "destroy_succeeded"},
		Target:  []string{""},
		Refresh: func() (interface{}, string, error) {
			statusResp, _, _ := client.ClustersApi.ClustersStatus(auth, clusterId).Execute()
			return statusResp, string(statusResp.Status), nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Minute,
		MinTimeout: 5 * time.Second,
	}
	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vpc cluster (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

// expandVpcServersSet is responsible for converting the servers set into
// a slice of type Server
func expandVpcServersSet(servers *schema.Set) []couchbasecapella.Server {
	result := make([]couchbasecapella.Server, servers.Len())

	for i, value := range servers.List() {
		v := value.(map[string]interface{})
		result[i] = createVpcServer(v)
	}

	return result
}

// expandVpcServicesList is responsible for converting the services interface into
// a slice of type CouchbaseServices
func expandVpcServiceList(services []interface{}) (res []couchbasecapella.CouchbaseServices) {
	for _, v := range services {
		res = append(res, couchbasecapella.CouchbaseServices(v.(string)))
	}

	return res
}

func getVpcServersProvider(servers *schema.Set) []string {
	providers := make([]string, 0)

	for _, value := range servers.List() {
		server := value.(map[string]interface{})
		for k, v := range server {
			if k == "aws" && len(v.(*schema.Set).List()) > 0 {
				if !Has(providers, "aws") {
					providers = append(providers, "aws")
				}
			}
			if k == "azure" && len(v.(*schema.Set).List()) > 0 {
				if !Has(providers, "azure") {
					providers = append(providers, "azure")
				}
			}
		}
	}
	return providers
}

func createVpcServer(v map[string]interface{}) couchbasecapella.Server {
	var server couchbasecapella.Server
	for _, awss := range v["aws"].(*schema.Set).List() {
		aws, ok := awss.(map[string]interface{})
		if ok {
			server = couchbasecapella.Server{
				Size:     int32(v["size"].(int)),
				Services: expandVpcServiceList(v["services"].([]interface{})),
				Aws: &couchbasecapella.ServerAws{
					InstanceSize: couchbasecapella.AwsInstances(aws["instance_size"].(string)),
					EbsSizeGib:   int32(aws["ebs_size_gib"].(int)),
				},
			}
		}
	}
	for _, azures := range v["azure"].(*schema.Set).List() {
		azure, ok := azures.(map[string]interface{})
		if ok {
			server = couchbasecapella.Server{
				Size:     int32(v["size"].(int)),
				Services: expandVpcServiceList(v["services"].([]interface{})),
				Azure: &couchbasecapella.ServerAzure{
					InstanceSize: couchbasecapella.AzureInstances(azure["instance_size"].(string)),
					VolumeType:   couchbasecapella.AzureVolumeTypes(azure["volume_type"].(string)),
				},
			}
		}
	}

	return server
}
