package provider

import (
	"context"
	"net/http"
	"os"

	couchbasecloud "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		couchbasecloud.ContextAPIKeys,
		map[string]couchbasecloud.APIKey{
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
		default:
			return diag.FromErr(err)
		}
	}
	return nil
}
