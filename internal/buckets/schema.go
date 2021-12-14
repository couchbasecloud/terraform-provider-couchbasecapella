package buckets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	schemaClusterIdKey          = "cluster_id"
	schemaBucketNameKey         = "name"
	schemaMemoryKey             = "memory_quota"
	schemaReplicasKey           = "replicas"
	schemaConflictResolutionKey = "conflict_resolution"
)

// Schema defines the mapping of resource data from Terraform to what is
// expected from our client API. This allows us to use our own type and not
// have to `Get` the fields and cast from a resource data type.
type Schema struct {
	// ClusterID is the id of the cluster
	ClusterID string

	// ConflictResolution defines the conflict resolution for the bucket
	ConflictResolution string

	// MemoryQuotaInMb is the memory quota in megabytes to set for the bucket
	MemoryQuotaInMb int32

	// Name of the bucket
	Name string

	// Replicas is the amount of replicas for the cluster
	Replicas int32
}

// NewSchema returns a schema that contains its field set based on the
// incoming resource data. If any of the required fields are not present, an
// error is returned to indicate how many and which fields are missing from
// the resource data.
func NewSchema(d *schema.ResourceData) (*Schema, error) {
	// ensure we have the right resource fields set
	s, err := formSchema(d)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// formSchema is meant to ensure the correct fields are set in the resource
// data and to return the schema with the appropriate fields set.
// NOTE: this does not validate that the fields themselves are valid
// in the context of a Couchbase Capella bucket operation. i.e. passing in an
// unsupported conflict resolution
func formSchema(d *schema.ResourceData) (*Schema, error) {
	s := new(Schema)
	var missingFields []string

	for _, chk := range []struct {
		field string
		set   func(s *Schema, d *schema.ResourceData) error
	}{
		{
			field: schemaClusterIdKey,
			set:   func(s *Schema, d *schema.ResourceData) error { return s.withClusterId(d) },
		},

		{
			field: schemaMemoryKey,
			set:   func(s *Schema, d *schema.ResourceData) error { return s.withMemory(d) },
		},

		{
			field: schemaReplicasKey,
			set:   func(s *Schema, d *schema.ResourceData) error { return s.withReplicas(d) },
		},

		{
			field: schemaBucketNameKey,
			set:   func(s *Schema, d *schema.ResourceData) error { return s.withBucketName(d) },
		},

		{
			field: schemaConflictResolutionKey,
			set:   func(s *Schema, d *schema.ResourceData) error { return s.withConflictResolution(d) },
		},
	} {
		if err := chk.set(s, d); err != nil {
			missingFields = append(missingFields, chk.field)
		}
	}

	if len(missingFields) > 0 {
		return nil, fmt.Errorf(
			"missing (%d) required fields: %s",
			len(missingFields),
			strings.Join(missingFields, ","),
		)
	}

	return s, nil
}

// withBucketName ensures the bucket name key is inside the resource data
// and sets the schema field if so. Returns ResourceMissingBucketName otherwise
func (s *Schema) withBucketName(d *schema.ResourceData) error {
	v, ok := d.GetOk(schemaBucketNameKey)
	if !ok {
		//  example of handling an error here
		return ResourceMissingBucketName
	}
	s.Name = v.(string)

	return nil
}

// withClusterId ensures the cluster id key is inside the resource data
// and sets the schema field if so. Returns ResourceMissingClusterId otherwise
func (s *Schema) withClusterId(d *schema.ResourceData) error {
	v, ok := d.GetOk(schemaClusterIdKey)
	if !ok {
		//  example of handling an error here
		return ResourceMissingBucketName
	}
	s.ClusterID = v.(string)

	return nil
}

// withReplicas ensures the replicas key is inside the resource data
// and sets the schema field if so. Returns ResourceMissingReplicas otherwise
func (s *Schema) withReplicas(d *schema.ResourceData) error {
	v, ok := d.GetOk(schemaReplicasKey)
	if !ok {
		//  example of handling an error here
		return ResourceMissingReplicas
	}
	s.Replicas = int32(v.(int))

	return nil
}

// withMemory ensures the bucket memory key is inside the resource data
// and sets the schema field if so. Returns ResourceMissingMemory otherwise
func (s *Schema) withMemory(d *schema.ResourceData) error {
	v, ok := d.GetOk(schemaMemoryKey)
	if !ok {
		//  example of handling an error here
		return ResourceMissingMemory
	}
	s.MemoryQuotaInMb = int32(v.(int))

	return nil
}

// withConflictResolution ensures the conflict resolution key is inside the
// resource data and sets the schema field if so. Returns
// ResourceMissingConflictResolution otherwise
func (s *Schema) withConflictResolution(d *schema.ResourceData) error {
	v, ok := d.GetOk(schemaConflictResolutionKey)
	if !ok {
		//  example of handling an error here
		return ResourceMissingConflictResolution
	}
	s.ConflictResolution = v.(string)

	return nil
}
