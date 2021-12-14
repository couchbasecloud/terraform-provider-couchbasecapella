package buckets

import "github.com/JamesWilkinsonCB/terraform-provider-couchbasecapella/internal/errors"

const (
	// missing stuff
	ResourceMissingClusterID          errors.Error = "the resource is missing the cluster id. something something something"
	ResourceMissingBucketName         errors.Error = "the resource is missing the bucket name. something something something"
	ResourceMissingMemory             errors.Error = "the resource is missing the memory quota for the bucket. something something something"
	ResourceMissingReplicas           errors.Error = "the resource is missing the replicas field. something something something"
	ResourceMissingConflictResolution errors.Error = "the resource is missing the conflict resolution. something something something"

	// domain specific stuff
	// NOTE: I'm not sure what the best practice is with TF and handling this
	// domain logic within a provider.
	ResourceMemoryTooLow              errors.Error = "the memory quota is too low. blah blah blah"
	ResourceInvalidReplicaCount       errors.Error = "the replica count is incorrect. blah blah blah"
	ResourceInvalidConflictResolution errors.Error = "the conflict resolution is not supported. blah blah blah"

	BucketResourceNotFound errors.Error = "the bucket resource at the specified name was not found"
)
