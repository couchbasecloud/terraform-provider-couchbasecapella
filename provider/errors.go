package provider

// Error provides a custom type for creating named errors in various packages
type Error string

// Error implements the Error interface
func (e Error) Error() string { return string(e) }

const (
	BucketHostedNotSupported        string = "this current release of the terraform provider doesn't support managing buckets in hosted clusters, please log in to the Capella UI where you can update your cluster"
	BucketInvalidName               string = "use letters, numbers, periods (.) or dashes (-). Bucket names cannot exceed 100 characters and must begin with a letter or a number"
	BucketInvalidMemoryQuota        string = "expected a value greater than 100 MiB, got %v MiB"
	BucketInvalidConflictResolution string = "expected a valid value for conflict resolution {lww, seqno}, got %s"

	DatabaseUserHostedNotSupported     string = "this current release of the terraform provider doesn't support managing database users in hosted clusters, please log in to the Capella UI where you can update your cluster"
	DatabaseUserInvalidPassword        string = "password must contain 8+ characters, 1+ lowercase, 1+ uppercase, 1+ symbols, 1+ numbers"
	DatabaseUserInvalidBucketAccess    string = "expected a valid value for bucket access {data_reader, data_writer}, got %s"
	DatabaseUserInvalidAllBucketAccess string = "expected a valid value for all bucket access {data_reader, data_writer}, got %s"

	HostedClusterInvalidProvider               string = "expected a valid value for provider {aws, azure}, got %s"
	HostedClusterInvalidRegion                 string = "expected a valid region for the cloud provider, got %s"
	HostedClusterInvalidCIDR                   string = "expected a valid CIDR address, got %s"
	HostedClusterInvalidSupportPackageTimezone string = "expected a valid value for timzone {ET, GMT, IST, PT}, got %s"
	HostedClusterInvalidSupportPackageType     string = "expected a valid value for support package type {Basic, DeveloperPro, Enterprise}, got %s"
	HostedClusterInvalidCompute                string = "expected a valid value for compute instance, got %s"
	HostedClusterInvalidIOPS                   string = "if storage type is GP3, iops should be a value between 3000 and 16000. If storage type is IO2, iops should be a value between 1000 and 64000"

	VpcClusterUpdateNotSupported         string = "This current release of the terraform provider doesn't support updating vpc clusters, please log in to the Capella UI where you can update your cluster"
	VpcClusterInvalidAwsInstance         string = "expected a valid value Aws instance, got %s"
	VpcClusterInvalidAzureInstance       string = "expected a valid value Azure instance, got %s"
	VpcClusterInvalidAzureVolumeSize     string = "expected a valid value for Azure size, got %s"
	VpcClusterServerDoesNotMatchProvider string = "cluster's server should be the same as the cloud provider"

	ClusterInvalidName             string = "cluster name can include letters, numbers, spaces, periods (.), dashes (-), and underscores (_). Cluster name should be between 2 and 128 characters and must begin with a letter or a number"
	ClusterInvalidSize             string = "expected a value between 2 and 27, got %v"
	ClusterInvalidCouchbaseService string = "expected a valid value for service {data, index, query, search, eventing, analytics}, got %s"
	ClusterInvalidStorageType      string = "expected a valid value for storage type {GP3, IO2}, got %s"
	ClusterProblemAccessing        string = "a problem occurred while accessing the cluster"
	ClusterInvalidStorageSize      string = "expected a value between 50 and 16000, got %v"

	ProjectDeleteClustersStillAssociated string = "Project cannot be deleted whilst there are still clusters associated with the project"
)
