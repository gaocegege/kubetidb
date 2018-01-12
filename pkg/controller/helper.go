package controller

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/controller"

	api "github.com/gaocegege/kubetidb/pkg/apis/tidb/v1alpha1"
	clientset "github.com/gaocegege/kubetidb/pkg/clientset/versioned"
	"github.com/gaocegege/kubetidb/pkg/controller/control"
)

var (
	groupVersionKind = schema.GroupVersionKind{
		Group:   api.GroupName,
		Version: api.GroupVersion,
		Kind:    api.TiDBResourceKind,
	}
)

// HelperInterface is the interface for helper.
type HelperInterface interface {
	// CreateService(tidb *api.TiDB, service *v1.Service) error
	// CreatePod(tidb *api.TiDB, template *v1.PodTemplateSpec) error
	// GetPodsForTiDB(tidb *api.TiDB, typ api.TFReplicaType) ([]*v1.Pod, error)
	// GetServicesForTiDB(tidb *api.TiDB, typ api.TFReplicaType) ([]*v1.Service, error)
}

// Helper is the type to manage internal resources in Kubernetes.
type Helper struct {
	tidbClientset clientset.Interface

	podLister  kubelisters.PodLister
	podControl controller.PodControlInterface

	serviceLister  kubelisters.ServiceLister
	serviceControl control.ServiceControlInterface
}

// NewHelper creates a new Helper.
func NewHelper(tidbClientset clientset.Interface, podLister kubelisters.PodLister, podControl controller.PodControlInterface, serviceLister kubelisters.ServiceLister, serviceControl control.ServiceControlInterface) *Helper {
	return &Helper{
		tidbClientset:  tidbClientset,
		podLister:      podLister,
		podControl:     podControl,
		serviceLister:  serviceLister,
		serviceControl: serviceControl,
	}
}
