package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	couchbasecloud "github.com/couchbaselabs/couchbase-cloud-go-client"
)

func resourceCouchbaseCloudProject() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Couchbase Projects.",

		CreateContext: resourceCouchbaseCloudProjectCreate,
		ReadContext:   resourceCouchbaseCloudProjectRead,
		DeleteContext: resourceCouchbaseCloudProjectDelete,

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

func resourceCouchbaseCloudProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)

	createProjectRequest := *couchbasecloud.NewCreateProjectRequest(d.Get("name").(string))

	project, _, err := client.ProjectsApi.ProjectsCreate(auth).CreateProjectRequest(createProjectRequest).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.Id)

	return resourceCouchbaseCloudProjectRead(ctx, d, meta)
}

func resourceCouchbaseCloudProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)
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

func resourceCouchbaseCloudProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecloud.APIClient)
	auth := getAuth(ctx)

	projectId := d.Get("id").(string)

	_, err := client.ProjectsApi.ProjectsDelete(auth, projectId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
