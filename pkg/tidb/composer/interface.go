package composer

import "k8s.io/api/core/v1"

type Composer interface {
	Compose() *v1.PodTemplateSpec
}
