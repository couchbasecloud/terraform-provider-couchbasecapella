// Couchbase, Inc. licenses this to you under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0.

// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package provider

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func resourceCouchbaseCapellaHostedCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Capella hosted clusters.",

		CreateContext: resourceCouchbaseCapellaHostedClusterCreate,
		ReadContext:   resourceCouchbaseCapellaHostedClusterRead,
		UpdateContext: resourceCouchbaseCapellaHostedClusterUpdate,
		DeleteContext: resourceCouchbaseCapellaHostedClusterDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Cluster's id.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
			},
			"name": {
				Description: "Cluster's name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Cluster's description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"project_id": {
				Description: "Project's Id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"place": {
				Description: "Cluster's place",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"single_az": {
							Type:        schema.TypeBool,
							Description: "Single AZ",
							Required:    true,
						},
						"hosted": {
							Type:        schema.TypeSet,
							Description: "Support package type",
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider": {
										Description: "Cloud provider",
										Type:        schema.TypeString,
										Required:    true,
									},
									"region": {
										Description: "A valid region for the cloudProvider",
										Type:        schema.TypeString,
										Required:    true,
									},
									"cidr": {
										Description: "CIDR block",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"support_package": {
				Description: "Suport package for cluster",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timezone": {
							Type:        schema.TypeString,
							Description: "Timezone",
							Required:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Support package type",
							Required:    true,
						},
					},
				},
			},
			"servers": {
				Description: "Cluster servers configuration",
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Description: "Number of nodes",
							Required:    true,
						},
						"compute": {
							Type:        schema.TypeString,
							Description: "Compute instance type",
							Required:    true,
						},
						"services": {
							Type:        schema.TypeList,
							Description: "Services",
							Required:    true,
							MinItems:    1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"storage": {
							Description: "Storage configuration.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Storage type",
										Type:        schema.TypeString,
										Required:    true,
									},
									"iops": {
										Description: "IOPS",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"size": {
										Description: "size(Gb).",
										Type:        schema.TypeInt,
										Required:    true,
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
			Update: schema.DefaultTimeout(25 * time.Minute),
		},
	}
}

func resourceCouchbaseCapellaHostedClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	environment := "hosted"
	clusterName := d.Get("name").(string)
	projectId := d.Get("project_id").(string)
	servers := expandHostedServersSet(d.Get("servers").(*schema.Set))
	supportPackage := expandHostedSupportPackageSet((d.Get("support_package")).(*schema.Set))
	place := expandHostedPlaceSet(d.Get("place").(*schema.Set))

	// force Single AZ to true if support package is Basic
	if supportPackage.Type == couchbasecapella.V3_BASIC {
		place.SingleAZ = true
	}

	newClusterRequest := *couchbasecapella.NewV3CreateClusterRequest(couchbasecapella.V3Environment(environment), clusterName, projectId,
		place, servers, supportPackage)

	if d.Get("description") != nil {
		description := d.Get("description").(string)
		newClusterRequest.SetDescription(description)
	}

	// Create the cluster
	response, err := client.ClustersV3Api.ClustersV3create(auth).V3CreateClusterRequest(newClusterRequest).Execute()
	if err != nil {
		return manageErrors(err, *response, "Create Hosted Cluster")
	}

	// TODO: need to be changed after cloud api fix!
	location := string(response.Header.Get("Location"))
	urlparts := strings.Split(location, "/")
	clusterId := urlparts[len(urlparts)-1]
	d.SetId(clusterId)

	defer response.Body.Close()

	// Wait for the cluster to deploy
	createStateConf := &resource.StateChangeConf{
		Pending: []string{"deploying"},
		Target:  []string{"healthy"},
		Refresh: func() (interface{}, string, error) {
			statusResp, _, err := client.ClustersV3Api.ClustersV3status(auth, clusterId).Execute()
			if err != nil {
				return 0, "Error", err
			}
			return statusResp, string(statusResp.Status), nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Minute,
		MinTimeout: 30 * time.Second,
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for cluster (%s) to be created: %s", d.Id(), err)
	}

	return resourceCouchbaseCapellaHostedClusterRead(ctx, d, meta)
}

func resourceCouchbaseCapellaHostedClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)
	clusterId := d.Get("id").(string)

	cluster, resp, err := client.ClustersV3Api.ClustersV3show(auth, clusterId).Execute()

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if err := d.Set("name", cluster.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCouchbaseCapellaHostedClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("id").(string)

	// Name and Description Update
	if d.HasChange("name") || d.HasChange("description") {
		v3UpdateClusterMetaRequest := *couchbasecapella.NewV3UpdateClusterMetaRequest()
		v3UpdateClusterMetaRequest.SetName(d.Get("name").(string))
		v3UpdateClusterMetaRequest.SetDescription((d.Get("description").(string)))
		_, err := client.ClustersV3Api.ClustersV3updateMeta(auth, clusterId).V3UpdateClusterMetaRequest(v3UpdateClusterMetaRequest).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Support Package Update
	if d.HasChange("support_package") {
		supportPackage := expandHostedSupportPackageSet((d.Get("support_package")).(*schema.Set))
		v3UpdateClusterSupportRequestSupportPackage := couchbasecapella.V3UpdateClusterSupportRequestSupportPackage{
			Timezone: &supportPackage.Timezone,
			Type:     supportPackage.Type,
		}
		v3UpdateClusterSupportRequest := couchbasecapella.V3UpdateClusterSupportRequest{
			SupportPackage: v3UpdateClusterSupportRequestSupportPackage,
		}
		_, err := client.ClustersV3Api.ClustersV3updateSupport(auth, clusterId).V3UpdateClusterSupportRequest(v3UpdateClusterSupportRequest).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Servers Update
	if d.HasChange("servers") {
		servers := expandHostedServersSet(d.Get("servers").(*schema.Set))
		v3UpdateClusterServersRequest := *couchbasecapella.NewV3UpdateClusterServersRequest(servers) // V3UpdateClusterServersRequest |  (optional)
		_, err := client.ClustersV3Api.ClustersV3updateServers(auth, clusterId).V3UpdateClusterServersRequest(v3UpdateClusterServersRequest).Execute()
		if err != nil {
			return diag.FromErr(err)
		}

		// Wait for the cluster to deploy
		updateStateConf := &resource.StateChangeConf{
			Pending: []string{"deploying"},
			Target:  []string{"healthy"},
			Refresh: func() (interface{}, string, error) {
				statusResp, _, err := client.ClustersV3Api.ClustersV3status(auth, clusterId).Execute()
				if err != nil {
					return 0, "Error", err
				}
				return statusResp, string(statusResp.Status), nil
			},
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      2 * time.Minute,
			MinTimeout: 30 * time.Second,
		}
		_, err = updateStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error waiting for hosted cluster (%s) to be updated: %s", d.Id(), err)
		}
	}

	return resourceCouchbaseCapellaHostedClusterRead(ctx, d, meta)
}

func resourceCouchbaseCapellaHostedClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	clusterId := d.Get("id").(string)

	// Check that Cluster is ready to be destroyed
	statusResp, _, err := client.ClustersV3Api.ClustersV3status(auth, clusterId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	if statusResp.Status != couchbasecapella.V3_HEALTHY {
		return diag.Errorf("Cluster is not ready to be deleted. Cluster Status: %s", statusResp.Status)
	}

	r, err2 := client.ClustersV3Api.ClustersV3delete(auth, clusterId).Execute()
	if err2 != nil {
		return manageErrors(err2, *r, "Hosted Cluster Delete")
	}

	// Wait for the cluster to be destroyed
	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{"destroying"},
		Target:  []string{""},
		Refresh: func() (interface{}, string, error) {
			statusResp, _, _ := client.ClustersV3Api.ClustersV3status(auth, clusterId).Execute()
			return statusResp, string(statusResp.Status), nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Minute,
		MinTimeout: 5 * time.Second,
	}
	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for hosted cluster (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func expandHostedServersSet(servers *schema.Set) []couchbasecapella.V3Servers {
	result := make([]couchbasecapella.V3Servers, servers.Len())

	for i, value := range servers.List() {
		v := value.(map[string]interface{})
		result[i] = createHostedServer(v)
	}

	return result
}

func expandHostedServiceList(services []interface{}) (res []couchbasecapella.V3CouchbaseServices) {
	for _, v := range services {
		res = append(res, couchbasecapella.V3CouchbaseServices(v.(string)))
	}

	return res
}

func createHostedServer(v map[string]interface{}) couchbasecapella.V3Servers {
	var server couchbasecapella.V3Servers
	for _, storages := range v["storage"].(*schema.Set).List() {
		storage, ok := storages.(map[string]interface{})
		if ok {
			server = couchbasecapella.V3Servers{
				Size:     int32(v["size"].(int)),
				Compute:  v["compute"].(string),
				Services: expandHostedServiceList(v["services"].([]interface{})),
				Storage: couchbasecapella.V3ServersStorage{
					Type: couchbasecapella.V3StorageType((storage["type"].(string))),
					IOPS: int32(storage["iops"].(int)),
					Size: int32(storage["size"].(int)),
				},
			}
		}
	}

	return server
}

func expandHostedPlaceSet(place *schema.Set) couchbasecapella.V3Place {
	result := make([]couchbasecapella.V3Place, place.Len())

	for i, value := range place.List() {
		v := value.(map[string]interface{})
		result[i] = createHostedPlace(v)
	}

	return result[0]
}

func createHostedPlace(v map[string]interface{}) couchbasecapella.V3Place {
	var place couchbasecapella.V3Place
	for _, hosteds := range v["hosted"].(*schema.Set).List() {
		hosted, ok := hosteds.(map[string]interface{})
		if ok {
			place = couchbasecapella.V3Place{
				SingleAZ: v["single_az"].(bool),
				Hosted: &couchbasecapella.V3PlaceHosted{
					Provider: couchbasecapella.V3Provider((hosted["provider"].(string))),
					Region:   hosted["region"].(string),
					CIDR:     hosted["cidr"].(string),
				},
			}
		}
	}

	return place
}

func expandHostedSupportPackageSet(supportPackage *schema.Set) couchbasecapella.V3SupportPackage {
	result := make([]couchbasecapella.V3SupportPackage, supportPackage.Len())

	for i, value := range supportPackage.List() {
		v := value.(map[string]interface{})
		result[i] = createHostedSupportPackage(v)
	}

	return result[0]
}

func createHostedSupportPackage(v map[string]interface{}) couchbasecapella.V3SupportPackage {
	var supportPackage couchbasecapella.V3SupportPackage
	supportPackage = couchbasecapella.V3SupportPackage{
		Timezone: couchbasecapella.V3SupportPackageTimezones((v["timezone"].(string))),
		Type:     couchbasecapella.V3SupportPackageType((v["type"].(string))),
	}

	return supportPackage
}
