package provider

import (
	"context"
	"os"

	couchbasecloud "github.com/couchbaselabs/couchbase-cloud-go-client"
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
