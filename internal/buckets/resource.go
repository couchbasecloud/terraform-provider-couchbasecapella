package buckets

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// NewCouchbaseCapellaBucketResource returns a resource with the correct schema
// mapping and handler functions done by the gateway
func NewCouchbaseCapellaBucketResource(gw *Gateway) *schema.Resource {
	// NOTE: I'm not too much of a fan of this but there's not too much we can
	// do to make it better
	return &schema.Resource{
		Description:   "Manage Couchbase Buckets.",
		CreateContext: gw.ResourceCouchbaseCapellaBucketCreate,
		ReadContext:   gw.ResourceCouchbaseCapellaBucketRead,
		//		UpdateContext: gw.ResourceCouchbaseCapellaBucketUpdate,
		//		DeleteContext: gw.ResourceCouchbaseCapellaBucketDelete,
		Schema: couchbaseCapellaSchemaMap(),
	}
}

func couchbaseCapellaSchemaMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		schemaClusterIdKey: {
			Description: "Cluster's id.",
			Type:        schema.TypeString,
			Required:    true,
		},
		schemaBucketNameKey: {
			Description: "Bucket's name.",
			Type:        schema.TypeString,
			Required:    true,
		},
		schemaMemoryKey: {
			Description: "Bucket Memory quota.",
			Type:        schema.TypeInt,
			Required:    true,
		},
		schemaReplicasKey: {
			Description: "replicas.",
			Type:        schema.TypeInt,
			Required:    true,
		},
		schemaConflictResolutionKey: {
			Description: "replicas.",
			Type:        schema.TypeString,
			Required:    true,
		},
	}
}
