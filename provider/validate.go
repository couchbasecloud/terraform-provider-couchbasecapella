package provider

import (
	"fmt"
	"regexp"

	couchbasecapella "github.com/couchbaselabs/couchbase-cloud-go-client"
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
	region := val.(string)
	awsRegionValidation := couchbasecapella.AwsRegions(region).IsValid()
	azureRegionValidation := couchbasecapella.AzureRegions(region).IsValid()
	if !awsRegionValidation && !azureRegionValidation {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidRegion, region))
	}
	return
}

func validateProvider(val interface{}, key string) (warns []string, errs []error) {
	provider := val.(string)
	providerValidation := couchbasecapella.V3Provider(provider).IsValid()
	// Temporary check for gcp provider whilst it is not yet available in Capella
	if provider == "gcp" {
		errs = append(errs, fmt.Errorf("GCP is not yet available in Couchbase Capella"))
		return nil, errs
	}
	if !providerValidation {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidProvider, provider))
	}
	return nil, errs
}

func validateCIDR(val interface{}, key string) (warns []string, errs []error) {
	region := val.(string)
	awsRegionValidation := couchbasecapella.AwsRegions(region).IsValid()
	azureRegionValidation := couchbasecapella.AzureRegions(region).IsValid()
	if !awsRegionValidation && !azureRegionValidation {
		errs = append(errs, fmt.Errorf(HostedClusterInvalidRegion, region))
	}
	return
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
	sizeIsValid := size >= 3 && size < 28
	if !sizeIsValid {
		errs = append(errs, fmt.Errorf(ClusterInvalidSize, size))
	}
	return
}

func validateCompute(val interface{}, key string) (warns []string, errs []error) {
	instance := val.(string)
	awsInstanceValidation := couchbasecapella.AwsInstances(instance).IsValid()
	azureInstanceValidation := couchbasecapella.AzureInstances(instance).IsValid()
	if !awsInstanceValidation && !azureInstanceValidation {
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
