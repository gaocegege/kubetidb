/*
Copyright 2018 Caicloud Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package control

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

// ServiceControlInterface is an interface that knows how to add or delete Services
// created as an interface to allow testing.
type ServiceControlInterface interface {
	// CreateServices creates new Services according to the spec.
	CreateServices(namespace string, service *v1.Service, object runtime.Object) error
	// CreateServicesWithControllerRef creates new services according to the spec, and sets object as the service's controller.
	CreateServicesWithControllerRef(namespace string, service *v1.Service, object runtime.Object, controllerRef *metav1.OwnerReference) error
	// PatchService patches the service.
	PatchService(namespace, name string, data []byte) error
}

// RealServiceControl is the default implementation of ServiceControlInterface.
type RealServiceControl struct {
	KubeClient clientset.Interface
	Recorder   record.EventRecorder
}

func (r RealServiceControl) PatchService(namespace, name string, data []byte) error {
	_, err := r.KubeClient.CoreV1().Services(namespace).Patch(name, types.StrategicMergePatchType, data)
	return err
}

func (r RealServiceControl) CreateServices(namespace string, service *v1.Service, object runtime.Object) error {
	return r.createServices(namespace, service, object, nil)
}

func (r RealServiceControl) CreateServicesWithControllerRef(namespace string, service *v1.Service, controllerObject runtime.Object, controllerRef *metav1.OwnerReference) error {
	if err := validateControllerRef(controllerRef); err != nil {
		return err
	}
	return r.createServices(namespace, service, controllerObject, controllerRef)
}

func (r RealServiceControl) createServices(namespace string, service *v1.Service, object runtime.Object, controllerRef *metav1.OwnerReference) error {
	if labels.Set(service.Labels).AsSelectorPreValidated().Empty() {
		return fmt.Errorf("unable to create Services, no labels")
	}

	newService, err := r.KubeClient.CoreV1().Services(namespace).Create(service)
	if err != nil {
		r.Recorder.Eventf(object, v1.EventTypeWarning, FailedCreateServiceReason, "Error creating: %v", err)
		return fmt.Errorf("unable to create services: %v", err)
	}

	accessor, err := meta.Accessor(object)
	if err != nil {
		glog.Errorf("parentObject does not have ObjectMeta, %v", err)
		return nil
	}
	glog.V(4).Infof("Controller %v created service %v", accessor.GetName(), newService.Name)
	r.Recorder.Eventf(object, v1.EventTypeNormal, SuccessfulCreateServiceReason, "Created service: %v", newService.Name)

	return nil
}
