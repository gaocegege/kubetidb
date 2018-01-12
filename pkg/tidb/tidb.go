package tidb

import (
	"k8s.io/api/core/v1"

	api "github.com/gaocegege/kubetidb/pkg/apis/tidb/v1alpha1"
)

type TiDB struct {
	activePDPods   []*v1.Pod
	pdPDServices   []*v1.Service
	activeTiKVPods []*v1.Pod
	tikvServices   []*v1.Service
	activeTiDBPods []*v1.Pod
	tidbServices   []*v1.Service

	tidb *api.TiDB
}

func New(tidb *api.TiDB, activePDPods []*v1.Pod, pdPDServices []*v1.Service, activeTiKVPods []*v1.Pod, tikvServices []*v1.Service, activeTiDBPods []*v1.Pod, tidbServices []*v1.Service) *TiDB {
	return &TiDB{
		tidb:           tidb,
		activePDPods:   activePDPods,
		pdPDServices:   pdPDServices,
		activeTiKVPods: activeTiKVPods,
		tikvServices:   tikvServices,
		activeTiDBPods: activeTiDBPods,
		tidbServices:   tidbServices,
	}
}

func (tidb *TiDB) Action() []Task {
	tasks := make([]Task, 0)

	expectedPD := int(*tidb.tidb.Spec.Specs[tidb.getPDIndex()])
	tasks = append(tasks, handle(CreatePDPod, , activePDPods, pdPDServices, expectedPD))
}

func (tidb *TiDB) handle(typ TaskType, pods []*v1.Pod, services []*v1.Service, expected int) []Task {
	tasks := make([]Task, 0)
	currentPods := len(activePDPods)
	if expected > currentPods {
		for i := 0; i < expected; i++ {
			if podMissed(pods, i) {
				tasks = append(tasks, Task{
					Type:typ,
					Template: tidb.composePodTemplate(expected)
				})
			}
		}
	}
}

func (tidb *TiDB) composePodTemplate(componentSpec *api.ComponentSpec, typ TaskType, replicas int) *v1.PodTemplateSpec {
	switch typ {
	case CreatePDPod:
		pdComposer := composer.NewPD(componentSpec, typ, replicas)
	}
}

func podMissed(pods []*v1.Pod, index int) bool {
	missed := true
	for _, pod := range pods {
		if int(pod.Labels[indexName]) == index {
			missed = false
			return missed
		}
	}
	return missed
}

func (tidb *TiDB) getPDIndex() int {
	for i, spec := range tidb.tidb.Spec.Specs {
		if spec.Type == api.TypePD {
			return i
		}
	}
	return -1
}
