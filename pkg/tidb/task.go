package tidb

import "k8s.io/api/core/v1"

type Task struct {
	Type     TaskType
	Template *v1.PodTemplateSpec
}

type TaskType string

const (
	CreatePDPod       TaskType = "CreatePDPod"
	CreatePDService            = "CreatePDService"
	CreateTiKVPod              = "CreateTiKVService"
	CreateTiKVService          = "CreateTiKVService"
	CreateTiDBPod              = "CreateTiDBPod"
	CreateTiDBService          = "CreateTiDBService"
)
