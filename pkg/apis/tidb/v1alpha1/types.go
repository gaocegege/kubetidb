package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TiDBCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterSpec   `json:"spec"`
	Status            ClusterStatus `json:"status"`
}

type ClusterSpec struct {
	PDSpec   PDSpec   `json:"pdSpec"`
	TiKVSpec TiKVSpec `json:"tikvSpec"`
	TiDBSpec TiDBSpec `json:"tidbSpec"`
}

type PDSpec struct {
	// Optional. The number of desired replicas. Default 1.
	Replicas *int32 `json:"replicas,omitempty"`
	// Template describes the data a pod should have when created from a template
	Template *v1.PodTemplateSpec `json:"template,omitempty"`
}

type TiKVSpec struct {
	// Optional. The number of desired replicas. Default 1.
	Replicas *int32 `json:"replicas,omitempty"`
	// Template describes the data a pod should have when created from a template
	Template *v1.PodTemplateSpec `json:"template,omitempty"`
}

type TiDBSpec struct {
	// Optional. The number of desired replicas. Default 1.
	Replicas *int32 `json:"replicas,omitempty"`
	// Template describes the data a pod should have when created from a template
	Template *v1.PodTemplateSpec `json:"template,omitempty"`
}

// ClusterStatus define the most recently observed status of the cluster.
type ClusterStatus struct {
	Phase ClusterPhase `json:"phase"`

	// Represents time when the cluster was acknowledged by the cluster controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the cluster was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents the latest available observations of a cluster object's current state.
	Conditions []*ClusterCondition `json:"conditions"`

	// InstanceStatus represents all of the instances' status in cluster.
	InstanceStatus InstanceStatus `json:"instanceStatus"`
}

type ClusterPhase string

const (
	TFJobNone      ClusterPhase = ""
	TFJobUnknown                = "Unknown"
	TFJobPending                = "Pending"
	TFJobRunning                = "Running"
	TFJobSucceeded              = "Succeeded"
	TFJobFailed                 = "Failed"
)

type InstanceStatus map[string]string

// ClusterCondition represents one current condition of a TiDB cluster.
// A condition might not show up if it is not happening.
// For example, if a cluster is not upgrading, the Upgrading condition would not show up.
// If a cluster is upgrading and encountered a problem that prevents the upgrade,
// the Upgrading condition's status will would be False and communicate the problem back.
type ClusterCondition struct {
	// Type of cluster condition.
	Type ClusterConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type ClusterConditionType string

const (
	ClusterConditionAvailable ClusterConditionType = "Available"
	ClusterConditionFailed                         = "Failed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TiDBList is a list of Foo resources
type TiDBClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []TiDBCluster `json:"items"`
}
