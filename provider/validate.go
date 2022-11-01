package provider

import (
	"fmt"
	"regexp"

	couchbasecapella "github.com/couchbasecloud/couchbase-capella-api-go-client"
)

func validateBucketName(val interface{}, key string) (warns []string, errs []error) {
	var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9-.]*$`).MatchString
	var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString
	name := val.(string)
	nameValidate := isStringAlphabetic(name) && len(name) > 0 && len(name) < 100 && isAlphaNumeric(name[0:1])
	if !nameValidate {
		errs = append(errs, fmt.Errorf(BucketInvalidName))
	}
	return
}

func validateMemoryQuota(val interface{}, key string) (warns []string, errs []error) {
	memory := val.(int)
	if memory < 100 {
		errs = append(errs, fmt.Errorf(BucketInvalidMemoryQuota, memory))
	}
	return
}

func validateConflictResolution(val interface{}, key string) (warns []string, errs []error) {
	conflict := val.(string)
	conflictValidation := couchbasecapella.ConflictResolution(conflict).IsValid()
	if !conflictValidation {
		errs = append(errs, fmt.Errorf(BucketInvalidConflictResolution, conflict))

	}
	return
}

func validateDatabaseUserPassword(val interface{}, key string) (warns []string, errs []error) {
	password := val.(string)
	passwordValidate := validatePassword(password)
	if !passwordValidate {
		errs = append(errs, fmt.Errorf(DatabaseUserInvalidPassword))
	}
	return
}

func validateBucketAccess(val interface{}, key string) (warns []string, errs []error) {
	access := val.(string)
	accessValidation := couchbasecapella.BucketRoleTypes(access).IsValid()
	if !accessValidation {
		errs = append(errs, fmt.Errorf(DatabaseUserInvalidBucketAccess, access))
	}
	return
}

func validateAllBucketAccess(val interface{}, key string) (warns []string, errs []error) {
	access := val.(string)
	accessValidation := couchbasecapella.BucketRoleTypes(access).IsValid()
	if !accessValidation {
		errs = append(errs, fmt.Errorf(DatabaseUserInvalidAllBucketAccess, access))

	}
	return
}

func validateClusterName(val interface{}, key string) (warns []string, errs []error) {
	var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9-_. ]*$`).MatchString
	var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString
	name := val.(string)
	nameValidate := isStringAlphabetic(name) && len(name) >= 2 && len(name) < 128 && isAlphaNumeric(name[0:1])
	if !nameValidate {
		errs = append(errs, fmt.Errorf(ClusterInvalidName))
	}
	return
}

func validateRegion(val interface{}, key string) (warns []string, errs []error) {
	var allowedGcpRegionsEnumValues = []string{
		"asia-east1",
		"asia-east2",
		"asia-northeast1",
		"asia-northeast2",
		"asia-northeast3",
		"asia-south1",
		"asia-south2",
		"asia-southeast1",
		"asia-southeast2",
		"australia-southeast1",
		"australia-southeast2",
		"europe-central2",
		"europe-north1",
		"europe-west1",
		"europe-west2",
		"europe-west3",
		"europe-west4",
		"europe-west6",
		"europe-west8",
		"northamerica-northeast1",
		"northamerica-northeast2",
		"southamerica-east1",
		"southamerica-west1",
		"us-east1",
		"us-east4",
		"us-west1",
		"us-west2",
		"us-west3",
		"us-west4",
		"us-central1",
		"us-central2",
	}

	region := val.(string)
	awsRegionValidation := couchbasecapella.AwsRegions(region).IsValid()
	azureRegionValidation := couchbasecapella.AzureRegions(region).IsValid()
	gcpRegionValidation := false
	for _, existing := range allowedGcpRegionsEnumValues {
		if existing == region {
			gcpRegionValidation = true
		}
	}

	if !awsRegionValidation && !azureRegionValidation && !gcpRegionValidation {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidRegion, region))
	}
	return
}

func validateProvider(val interface{}, key string) (warns []string, errs []error) {
	provider := val.(string)
	providerValidation := couchbasecapella.V3Provider(provider).IsValid()
	if !providerValidation {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidProvider, provider))
	}
	return nil, errs
}

func validateTimezone(val interface{}, key string) (warns []string, errs []error) {
	timezone := val.(string)
	timezoneValidation := couchbasecapella.V3SupportPackageTimezones(timezone).IsValid()
	if !timezoneValidation {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidSupportPackageTimezone, timezone))
	}
	return
}

func validateSupportPackageType(val interface{}, key string) (warns []string, errs []error) {
	packageType := val.(string)
	packageTypeValidation := couchbasecapella.V3SupportPackageType(packageType).IsValid()
	if !packageTypeValidation {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidSupportPackageType, packageType))
	}
	return
}

func validateSize(val interface{}, key string) (warns []string, errs []error) {
	size := val.(int)
	sizeIsValid := size >= 2 && size < 28
	if !sizeIsValid {
		errs = append(errs, fmt.Errorf(ClusterInvalidSize, size))
	}
	return
}

func validateCompute(val interface{}, key string) (warns []string, errs []error) {
	var allowedGcpComputeEnumValues = []string{
		"n2-standard-2",
		"n2-standard-4",
		"n2-standard-8",
		"n2-standard-16",
		"n2-standard-32",
		"n2-standard-48",
		"n2-standard-64",
		"n2-standard-80",
		"n2-highmem-2",
		"n2-highmem-4",
		"n2-highmem-8",
		"n2-highmem-16",
		"n2-highmem-32",
		"n2-highmem-48",
		"n2-highmem-64",
		"n2-highmem-80",
		"n2-highcpu-2",
		"n2-highcpu-4",
		"n2-highcpu-8",
		"n2-highcpu-16",
		"n2-highcpu-32",
		"n2-highcpu-48",
		"n2-highcpu-64",
		"n2-highcpu-80",
		"n2-custom-2-4096",
		"n2-custom-4-8192",
		"n2-custom-8-16384",
		"n2-custom-16-32768",
		"n2-custom-32-65536",
		"n2-custom-36-73728",
		"n2-custom-48-98304",
		"n2-custom-72-147456",
	}

	instance := val.(string)
	awsInstanceValidation := couchbasecapella.AwsInstances(instance).IsValid()
	azureInstanceValidation := couchbasecapella.AzureInstances(instance).IsValid()
	gcpInstanceValidation := false
	for _, existing := range allowedGcpComputeEnumValues {
		if existing == instance {
			gcpInstanceValidation = true
		}
	}

	if !awsInstanceValidation && !azureInstanceValidation && !gcpInstanceValidation {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidCompute, instance))
	}
	return
}

func validateService(val interface{}, key string) (warns []string, errs []error) {
	service := val.(string)
	serviceValidation := couchbasecapella.V3CouchbaseServices(service).IsValid()
	if !serviceValidation {
		errs = append(errs, fmt.Errorf(ClusterInvalidCouchbaseService, service))
	}
	return
}

func validateStorageType(val interface{}, key string) (warns []string, errs []error) {
	storageType := val.(string)
	storageTypeValidation := couchbasecapella.V3StorageType(storageType).IsValid()
	if !storageTypeValidation {
		errs = append(errs, fmt.Errorf(ClusterInvalidStorageType, storageType))
	}
	return
}

func validateIops(val interface{}, key string) (warns []string, errs []error) {
	iops := val.(int)
	GP3IopsIsValid := iops >= 3000 && iops <= 16000
	IO2IopsIsValid := iops >= 1000 && iops <= 64000
	if !GP3IopsIsValid && !IO2IopsIsValid {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidIOPS))
	}
	return
}

func validateStorageSize(val interface{}, key string) (warns []string, errs []error) {
	storageSize := val.(int)
	storageSizeIsValid := storageSize >= 50 && storageSize <= 16000
	if !storageSizeIsValid {
		errs = append(errs, fmt.Errorf(ClusterInvalidStorageSize, storageSize))
	}
	return
}

func validateAwsInstance(val interface{}, key string) (warns []string, errs []error) {
	instance := val.(string)
	instanceValidation := couchbasecapella.AwsInstances(instance).IsValid()
	if !instanceValidation {
		errs = append(errs, fmt.Errorf(VpcClusterInvalidAwsInstance, instance))
	}
	return
}

func validateAwsVolumeSize(val interface{}, key string) (warns []string, errs []error) {
	size := val.(int)
	sizeIsValid := size >= 50 && size < 16000
	if !sizeIsValid {
		errs = append(errs, fmt.Errorf(ClusterInvalidStorageSize, size))
	}
	return
}

func validateAzureInstance(val interface{}, key string) (warns []string, errs []error) {
	instance := val.(string)
	instanceValidation := couchbasecapella.AzureInstances(instance).IsValid()
	if !instanceValidation {
		errs = append(errs, fmt.Errorf(VpcClusterInvalidAzureInstance, instance))
	}
	return
}

func validateAzureVolume(val interface{}, key string) (warns []string, errs []error) {
	volume := val.(string)
	volumeValidation := couchbasecapella.AzureVolumeTypes(volume).IsValid()
	if !volumeValidation {
		errs = append(errs, fmt.Errorf(VpcClusterInvalidAzureVolumeSize, volume))
	}
	return
}
