package extract

import (
	"golang.stackrox.io/kube-linter/internal/k8sutil"
)

// Labels extracts labels from the given object.
func Labels(object k8sutil.Object) map[string]string {
	return object.GetLabels()
}
