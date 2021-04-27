// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1 "github.com/openshift/hive/apis/hive/v1"
	v1alpha1 "github.com/openshift/hive/apis/hiveinternal/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=hive.openshift.io, Version=v1
	case v1.SchemeGroupVersion.WithResource("checkpoints"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().Checkpoints().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clusterclaims"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().ClusterClaims().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clusterdeployments"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().ClusterDeployments().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clusterdeprovisions"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().ClusterDeprovisions().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clusterimagesets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().ClusterImageSets().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clusterpools"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().ClusterPools().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clusterprovisions"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().ClusterProvisions().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clusterrelocates"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().ClusterRelocates().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("clusterstates"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().ClusterStates().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("dnszones"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().DNSZones().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("fakeclusterinstalls"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().FakeClusterInstalls().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("hiveconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().HiveConfigs().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("machinepools"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().MachinePools().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("machinepoolnameleases"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().MachinePoolNameLeases().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("selectorsyncidentityproviders"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().SelectorSyncIdentityProviders().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("selectorsyncsets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().SelectorSyncSets().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("syncidentityproviders"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().SyncIdentityProviders().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("syncsets"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hive().V1().SyncSets().Informer()}, nil

		// Group=hiveinternal.openshift.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("clustersyncs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hiveinternal().V1alpha1().ClusterSyncs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("clustersyncleases"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Hiveinternal().V1alpha1().ClusterSyncLeases().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
