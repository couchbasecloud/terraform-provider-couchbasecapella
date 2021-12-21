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

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Description:  "Project name.",
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

// resourceCouchbaseCapellaProjectCreate is responsible for creating a
// project in Couchbase Capella using the Terraform resource data.
func resourceCouchbaseCapellaProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)
	projectName := d.Get("name").(string)

	createProjectRequest := *couchbasecapella.NewCreateProjectRequest(projectName)

	project, r, err := client.ProjectsApi.ProjectsCreate(auth).CreateProjectRequest(createProjectRequest).Execute()
	if r == nil {
		return diag.Errorf("Pointer to database project create http.Response is nil")
	}
	if err != nil {
		return manageErrors(err, *r, "Create Project")
	}

	d.SetId(project.Id)

	return resourceCouchbaseCapellaProjectRead(ctx, d, meta)
}

// resourceCouchbaseCapellaProjectRead is responsible for reading a
// project in Couchbase Capella using the Terraform resource data.
func resourceCouchbaseCapellaProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)
	projectId := d.Id()

	_, resp, err := client.ProjectsApi.ProjectsShow(auth, projectId).Execute()

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

// resourceCouchbaseCapellaProjectDelete is responsible for deleting a
// project in Couchbase Capella using the Terraform resource data.
func resourceCouchbaseCapellaProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	projectId := d.Id()

	// Check to see if project exists in Capella. If the project
	// exists, it will be deleted. If the project does not exist in Capella,
	// likely being deleted elsewhere, an error is thrown.
	_, _, err := client.ProjectsApi.ProjectsShow(auth, projectId).Execute()
	if err != nil {
		return diag.Errorf("Failed to delete: Project doesn't exist Capella")
	}
	r, err := client.ProjectsApi.ProjectsDelete(auth, projectId).Execute()
	if r == nil {
		return diag.Errorf("Pointer to project delete http.Response is nil")
	}
	if err != nil {
		return manageErrors(err, *r, "Delete Project")
	}
	return nil
}
