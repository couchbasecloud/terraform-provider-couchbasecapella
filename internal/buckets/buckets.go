package buckets

import (
	"context"
	"errors"
	"os"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.uber.org/zap"
)

const (
	Resource = "couchbasecapella_bucket"
)

// Gateway acts as a bridge between to terraform resources and our API client.
// The gateway is responsible for reading schemas, validating its contents,
// communicating with the API, and handling any errors returned from the API.
// The gateway is also responsible for returning delightful error/warning
// diagnostics to the user
type Gateway struct {
	// NOTE: not sure if logging in a terraform provider is "allowed"
	logger *zap.Logger
}

// NewGateway returns an instantiated instance of a gateway. A gateway provides
// the ability to bridge the Terraform client with our API Client. It has
// the following dependencies:
//
// logger - for structured logging
func NewGateway(logger *zap.Logger) (*Gateway, error) {
	if logger == nil {
		return nil, errors.New("unable to initialize a gateway due to the missing logger dependency")
	}

	return &Gateway{
		logger: logger,
	}, nil
}

// ResourceCouchbaseCapellaBucketCreate is responsible for creating the
// Couchbase Capella bucket using the Terraform resource data.
// TODO: More comments on error cases, what happens, why, etc
func (g *Gateway) ResourceCouchbaseCapellaBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, client, err := getAPIClient(ctx, meta)
	if err != nil {
		g.logger.Error("unable to get client and auth context")
	}

	s, err := NewSchema(d)
	if err != nil {
		g.logger.Error("unable to form schema from resource data", zap.Error(err))
		diag.FromErr(err)
	}

	// TODO: validate domain specific things here such as memory being too high
	//  and low. make sure conflict resolution is correct etc.
	//  I'm not sure whether this is best practice for a provider but if we do
	//  do domain logic checks
	if err := g.validateSchemaResourceCreation(*s); err != nil {
		g.logger.Error("unable to validate schema", zap.Error(err))
		diag.FromErr(err)
	}

	conflictResolution := couchbasecapella.ConflictResolution(s.ConflictResolution)
	couchbaseBucketSpec := couchbasecapella.NewCouchbaseBucketSpec(s.Name, s.MemoryQuotaInMb)
	couchbaseBucketSpec.SetReplicas(s.Replicas)
	couchbaseBucketSpec.SetConflictResolution(conflictResolution)

	_, r, err := client.ClustersApi.ClustersCreateBucket(ctx, s.ClusterID).CouchbaseBucketSpec(*couchbaseBucketSpec).Execute()
	if err != nil {
		if r != nil {
			// TODO:  HANDLE ERRORS GRACEFULLY HERE AND REPORT AN ERROR THAT
			//  MAKES SENSE TO THE USER SO THEY KNOW HOW TO FIX THE PROBLEM
			// diag.FromErr((handleResponse(r))
		}
		g.logger.Error("unable to create bucket", zap.Error(err))
		return diag.FromErr(err)
	}

	d.SetId(s.Name)

	return g.ResourceCouchbaseCapellaBucketRead(ctx, d, meta)
}

// ResourceCouchbaseCapellaBucketRead is responsible for reading a Couchbase
// Capella bucket using the Terraform resource data.
// TODO: More comments on error cases, what happens, why, etc
func (g *Gateway) ResourceCouchbaseCapellaBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*couchbasecapella.APIClient)
	auth := getAuth(ctx)

	s, err := NewSchema(d)
	if err != nil {
		g.logger.Error("unable to form schema from resource data", zap.Error(err))
		diag.FromErr(err)
	}

	// NOTE: Don't we have a singular bucket api call???
	buckets, r, err := client.ClustersApi.ClustersListBuckets(auth, s.ClusterID).Execute()
	if err != nil {
		if r != nil {
			// TODO:  HANDLE ERRORS GRACEFULLY HERE AND REPORT AN ERROR THAT
			//  MAKES SENSE TO THE USER SO THEY KNOW HOW TO FIX THE PROBLEM
			// diag.FromErr((handleResponse(r))
		}
		g.logger.Error("unable to list buckets", zap.Error(err))
		return diag.FromErr(err)
	}
	for _, bucket := range buckets {
		if bucket.Name == s.Name {
			return nil
		}
	}

	return diag.FromErr(BucketResourceNotFound)
}

// JUST AN EXAMPLE, MINIMAL CHECKS FOR THE SAKE OF JUST SHOWING AN EXAMPLE
func (g *Gateway) validateSchemaResourceCreation(s Schema) error {
	// TODO: do checks here to validate schema
	// NOTE: i forget what the exact number is
	if s.MemoryQuotaInMb < 256 {
		return ResourceMemoryTooLow
	}

	return nil
}

// JUST AN EXAMPLE, MINIMAL CHECKS FOR THE SAKE OF JUST SHOWING AN EXAMPLE
func (g *Gateway) validateSchemaResourceUpdate(s Schema) error {
	// TODO: do checks here to validate schema
	// NOTE: i forget what the exact number is
	if s.MemoryQuotaInMb < 256 {
		return ResourceMemoryTooLow
	}

	return nil
}

// JUST AN EXAMPLE, MINIMAL CHECKS FOR THE SAKE OF JUST SHOWING AN EXAMPLE
func (g *Gateway) validateSchemaResourceDelete(s Schema) error {
	// TODO: do checks here to validate schema
	// NOTE: i forget what the exact number is
	if s.MemoryQuotaInMb < 256 {
		return ResourceMemoryTooLow
	}

	return nil
}

func getAPIClient(ctx context.Context, meta interface{}) (context.Context, *couchbasecapella.APIClient, error) {
	client, ok := meta.(*couchbasecapella.APIClient)
	if !ok {
		return nil, nil, errors.New("API client not found in meta interface{}")
	}
	// NOTE: this probably should return an error if its not there
	ctx = getAuth(ctx)

	return ctx, client, nil
}

// TODO: Get Auth with both env variable and terraform ones
// NOTE: We should probably have throw an error if we can't find the
// credentials
func getAuth(ctx context.Context) context.Context {
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
	return auth
}
