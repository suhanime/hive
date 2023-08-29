// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v1 "github.com/openshift/hive/apis/hive/v1"
	v1alpha1 "github.com/openshift/hive/apis/hiveinternal/v1alpha1"
	hivev1 "github.com/openshift/hive/pkg/client/applyconfiguration/hive/v1"
	hiveinternalv1alpha1 "github.com/openshift/hive/pkg/client/applyconfiguration/hiveinternal/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=hive.openshift.io, Version=v1
	case v1.SchemeGroupVersion.WithKind("AlibabaCloudClusterDeprovision"):
		return &hivev1.AlibabaCloudClusterDeprovisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ArgoCDConfig"):
		return &hivev1.ArgoCDConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSAssociatedVPC"):
		return &hivev1.AWSAssociatedVPCApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSClusterDeprovision"):
		return &hivev1.AWSClusterDeprovisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSDNSZoneSpec"):
		return &hivev1.AWSDNSZoneSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSDNSZoneStatus"):
		return &hivev1.AWSDNSZoneStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSPrivateLinkConfig"):
		return &hivev1.AWSPrivateLinkConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSPrivateLinkInventory"):
		return &hivev1.AWSPrivateLinkInventoryApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSPrivateLinkSubnet"):
		return &hivev1.AWSPrivateLinkSubnetApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSPrivateLinkVPC"):
		return &hivev1.AWSPrivateLinkVPCApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSResourceTag"):
		return &hivev1.AWSResourceTagApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AWSServiceProviderCredentials"):
		return &hivev1.AWSServiceProviderCredentialsApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AzureClusterDeprovision"):
		return &hivev1.AzureClusterDeprovisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("AzureDNSZoneSpec"):
		return &hivev1.AzureDNSZoneSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("BackupConfig"):
		return &hivev1.BackupConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("BackupReference"):
		return &hivev1.BackupReferenceApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("CertificateBundleSpec"):
		return &hivev1.CertificateBundleSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("CertificateBundleStatus"):
		return &hivev1.CertificateBundleStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Checkpoint"):
		return &hivev1.CheckpointApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("CheckpointSpec"):
		return &hivev1.CheckpointSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterClaim"):
		return &hivev1.ClusterClaimApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterClaimCondition"):
		return &hivev1.ClusterClaimConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterClaimSpec"):
		return &hivev1.ClusterClaimSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterClaimStatus"):
		return &hivev1.ClusterClaimStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeployment"):
		return &hivev1.ClusterDeploymentApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeploymentCondition"):
		return &hivev1.ClusterDeploymentConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeploymentCustomization"):
		return &hivev1.ClusterDeploymentCustomizationApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeploymentCustomizationSpec"):
		return &hivev1.ClusterDeploymentCustomizationSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeploymentCustomizationStatus"):
		return &hivev1.ClusterDeploymentCustomizationStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeploymentSpec"):
		return &hivev1.ClusterDeploymentSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeploymentStatus"):
		return &hivev1.ClusterDeploymentStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeprovision"):
		return &hivev1.ClusterDeprovisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeprovisionCondition"):
		return &hivev1.ClusterDeprovisionConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeprovisionPlatform"):
		return &hivev1.ClusterDeprovisionPlatformApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeprovisionSpec"):
		return &hivev1.ClusterDeprovisionSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterDeprovisionStatus"):
		return &hivev1.ClusterDeprovisionStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterImageSet"):
		return &hivev1.ClusterImageSetApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterImageSetReference"):
		return &hivev1.ClusterImageSetReferenceApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterImageSetSpec"):
		return &hivev1.ClusterImageSetSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterIngress"):
		return &hivev1.ClusterIngressApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterInstallCondition"):
		return &hivev1.ClusterInstallConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterInstallLocalReference"):
		return &hivev1.ClusterInstallLocalReferenceApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterMetadata"):
		return &hivev1.ClusterMetadataApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterOperatorState"):
		return &hivev1.ClusterOperatorStateApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterPlatformMetadata"):
		return &hivev1.ClusterPlatformMetadataApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterPool"):
		return &hivev1.ClusterPoolApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterPoolClaimLifetime"):
		return &hivev1.ClusterPoolClaimLifetimeApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterPoolCondition"):
		return &hivev1.ClusterPoolConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterPoolReference"):
		return &hivev1.ClusterPoolReferenceApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterPoolSpec"):
		return &hivev1.ClusterPoolSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterPoolStatus"):
		return &hivev1.ClusterPoolStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterProvision"):
		return &hivev1.ClusterProvisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterProvisionCondition"):
		return &hivev1.ClusterProvisionConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterProvisionSpec"):
		return &hivev1.ClusterProvisionSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterProvisionStatus"):
		return &hivev1.ClusterProvisionStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterRelocate"):
		return &hivev1.ClusterRelocateApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterRelocateSpec"):
		return &hivev1.ClusterRelocateSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterState"):
		return &hivev1.ClusterStateApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ClusterStateStatus"):
		return &hivev1.ClusterStateStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ControllerConfig"):
		return &hivev1.ControllerConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ControllersConfig"):
		return &hivev1.ControllersConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ControlPlaneAdditionalCertificate"):
		return &hivev1.ControlPlaneAdditionalCertificateApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ControlPlaneConfigSpec"):
		return &hivev1.ControlPlaneConfigSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ControlPlaneServingCertificateSpec"):
		return &hivev1.ControlPlaneServingCertificateSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("DeploymentConfig"):
		return &hivev1.DeploymentConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("DNSZone"):
		return &hivev1.DNSZoneApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("DNSZoneCondition"):
		return &hivev1.DNSZoneConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("DNSZoneSpec"):
		return &hivev1.DNSZoneSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("DNSZoneStatus"):
		return &hivev1.DNSZoneStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("FailedProvisionAWSConfig"):
		return &hivev1.FailedProvisionAWSConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("FailedProvisionConfig"):
		return &hivev1.FailedProvisionConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("FeatureGateSelection"):
		return &hivev1.FeatureGateSelectionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("FeatureGatesEnabled"):
		return &hivev1.FeatureGatesEnabledApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("GCPClusterDeprovision"):
		return &hivev1.GCPClusterDeprovisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("GCPDNSZoneSpec"):
		return &hivev1.GCPDNSZoneSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("GCPDNSZoneStatus"):
		return &hivev1.GCPDNSZoneStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("HibernationConfig"):
		return &hivev1.HibernationConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("HiveConfig"):
		return &hivev1.HiveConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("HiveConfigCondition"):
		return &hivev1.HiveConfigConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("HiveConfigSpec"):
		return &hivev1.HiveConfigSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("HiveConfigStatus"):
		return &hivev1.HiveConfigStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("IBMClusterDeprovision"):
		return &hivev1.IBMClusterDeprovisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("InventoryEntry"):
		return &hivev1.InventoryEntryApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("KubeconfigSecretReference"):
		return &hivev1.KubeconfigSecretReferenceApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MachinePool"):
		return &hivev1.MachinePoolApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MachinePoolAutoscaling"):
		return &hivev1.MachinePoolAutoscalingApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MachinePoolCondition"):
		return &hivev1.MachinePoolConditionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MachinePoolNameLease"):
		return &hivev1.MachinePoolNameLeaseApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MachinePoolPlatform"):
		return &hivev1.MachinePoolPlatformApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MachinePoolSpec"):
		return &hivev1.MachinePoolSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MachinePoolStatus"):
		return &hivev1.MachinePoolStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("MachineSetStatus"):
		return &hivev1.MachineSetStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ManageDNSAWSConfig"):
		return &hivev1.ManageDNSAWSConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ManageDNSAzureConfig"):
		return &hivev1.ManageDNSAzureConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ManageDNSConfig"):
		return &hivev1.ManageDNSConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ManageDNSGCPConfig"):
		return &hivev1.ManageDNSGCPConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("OpenStackClusterDeprovision"):
		return &hivev1.OpenStackClusterDeprovisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("OvirtClusterDeprovision"):
		return &hivev1.OvirtClusterDeprovisionApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PatchEntity"):
		return &hivev1.PatchEntityApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Platform"):
		return &hivev1.PlatformApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("PlatformStatus"):
		return &hivev1.PlatformStatusApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("Provisioning"):
		return &hivev1.ProvisioningApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ReleaseImageVerificationConfigMapReference"):
		return &hivev1.ReleaseImageVerificationConfigMapReferenceApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SecretMapping"):
		return &hivev1.SecretMappingApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SecretReference"):
		return &hivev1.SecretReferenceApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SelectorSyncIdentityProvider"):
		return &hivev1.SelectorSyncIdentityProviderApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SelectorSyncIdentityProviderSpec"):
		return &hivev1.SelectorSyncIdentityProviderSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SelectorSyncSet"):
		return &hivev1.SelectorSyncSetApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SelectorSyncSetSpec"):
		return &hivev1.SelectorSyncSetSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("ServiceProviderCredentials"):
		return &hivev1.ServiceProviderCredentialsApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SpecificControllerConfig"):
		return &hivev1.SpecificControllerConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SyncIdentityProvider"):
		return &hivev1.SyncIdentityProviderApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SyncIdentityProviderCommonSpec"):
		return &hivev1.SyncIdentityProviderCommonSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SyncIdentityProviderSpec"):
		return &hivev1.SyncIdentityProviderSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SyncObjectPatch"):
		return &hivev1.SyncObjectPatchApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SyncSet"):
		return &hivev1.SyncSetApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SyncSetCommonSpec"):
		return &hivev1.SyncSetCommonSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("SyncSetSpec"):
		return &hivev1.SyncSetSpecApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("TaintIdentifier"):
		return &hivev1.TaintIdentifierApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("VeleroBackupConfig"):
		return &hivev1.VeleroBackupConfigApplyConfiguration{}
	case v1.SchemeGroupVersion.WithKind("VSphereClusterDeprovision"):
		return &hivev1.VSphereClusterDeprovisionApplyConfiguration{}

		// Group=hiveinternal.openshift.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSync"):
		return &hiveinternalv1alpha1.ClusterSyncApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSyncCondition"):
		return &hiveinternalv1alpha1.ClusterSyncConditionApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSyncLease"):
		return &hiveinternalv1alpha1.ClusterSyncLeaseApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSyncLeaseSpec"):
		return &hiveinternalv1alpha1.ClusterSyncLeaseSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSyncStatus"):
		return &hiveinternalv1alpha1.ClusterSyncStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FakeClusterInstall"):
		return &hiveinternalv1alpha1.FakeClusterInstallApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FakeClusterInstallSpec"):
		return &hiveinternalv1alpha1.FakeClusterInstallSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FakeClusterInstallStatus"):
		return &hiveinternalv1alpha1.FakeClusterInstallStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SyncResourceReference"):
		return &hiveinternalv1alpha1.SyncResourceReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SyncStatus"):
		return &hiveinternalv1alpha1.SyncStatusApplyConfiguration{}

	}
	return nil
}
