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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func resourceCouchbaseCapellaProject() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Projects.",

		CreateContext: resourceCouchbaseCapellaProjectCreate,
		ReadContext:   resourceCouchbaseCapellaProjectRead,
		DeleteContext: resourceCouchbaseCapellaProjectDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Project id.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
			},
			"name": {
				Description: "Project name.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
		},
	}
}

func resourceCouchbaseCapellaProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	createProjectRequest := *couchbasecapella.NewCreateProjectRequest(d.Get("name").(string))

	project, r, err := client.ProjectsApi.ProjectsCreate(auth).CreateProjectRequest(createProjectRequest).Execute()
	if err != nil {
		return manageErrors(err, *r, "Create Project")
	}

	d.SetId(project.Id)

	return resourceCouchbaseCapellaProjectRead(ctx, d, meta)
}

func resourceCouchbaseCapellaProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)
	projectId := d.Id()

	project, resp, err := client.ProjectsApi.ProjectsShow(auth, projectId).Execute()

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if err := d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCouchbaseCapellaProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	projectId := d.Id()

	r, err := client.ProjectsApi.ProjectsDelete(auth, projectId).Execute()
	if err != nil {
		return manageErrors(err, *r, "Delete Project")
	}

	return nil
}
