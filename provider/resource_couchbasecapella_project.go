package provider

import (
	"context"
	"net/http"
	"os"

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

	projectId := d.Get("id").(string)

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
	projectId := d.Get("id").(string)

	r, err := client.ProjectsApi.ProjectsDelete(auth, projectId).Execute()
	if err != nil {
		return manageErrors(err, *r, "Delete Project")
	}

	return nil
}
