package composer

import (
	"fmt"

	"k8s.io/api/core/v1"

	api "github.com/gaocegege/kubetidb/pkg/apis/tidb/v1alpha1"
)

const (
	pdClientPort = 2379
	pdPeerPort   = 2380
)

type PD struct {
	tidbName string
	replicas int
	spec     *api.ComponentSpec
}

func NewPD(tidb *api.ComponentSpec, tidbName string, replicas int) *PD {
	return &PD{
		tidbName: tidbName,
		replicas: replicas,
		spec:     spec,
	}
}

func (pd *PD) composePod(index int) *v1.PodTemplateSpec {
	template := pd.spec.Template.DeepCopy()
	template.Labels = composeLabel(index)
	template.Spec.Containers[0].Args = append(template.Spec.Containers[0].Args, pd.composeArgs(index)...)
	return template
}

func (pd *PD) composeArgs(index int) []string {
	args := make([]string, 6)
	args[0] = fmt.Sprintf("--name=pd%d", index)
	args[1] = fmt.Sprintf("--client-urls=http://0.0.0.0:%d", pdClientPort)
	args[2] = fmt.Sprintf("--advertise-client-urls=%s:%d", pd.composeServiceName(index), pdClientPort)
	args[3] = fmt.Sprintf("--peer-urls=http://0.0.0.0:%d", pdPeerPort)
	args[4] = fmt.Sprintf("--advertise-peer-urls=%s:%d", pd.composeServiceName(index), pdPeerPort)
	template := "--initial-cluster=\"%s\""
	args5 := ""
	for i := 0; i < pd.replicas; i++ {
		args5 = args5 + fmt.Sprintf("pd%d=%s:%d", index, pd.composeServiceName(index), pdPeerPort)
	}
	args[5] = fmt.Sprintf(template, args5)
	return args
}

func (pd *PD) composeServiceName(index int) string {
	return fmt.Sprintf("%s-pd-%d", pd.tidbName, index)
}

func composeLabel(index int) map[string]string {
	return map[string]string{
		"gaocegege.com": "true",
		"index":         string(index),
		"type":          "pd",
	}
}
