package common

import (
	"os"
	"sync"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/go-logr/logr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("util")

type ResourceManager struct {
	Client            client.Client
	PatchFactory      *PatchFactory
	Log               logr.Logger
	OperatorNamespace string
}

func IndexOf(l []string, s string) int {
	for i, elem := range l {
		if elem == s {
			return i
		}
	}
	return -1
}

func HasLabel(resource metav1.Object, label string) bool {
	labels := resource.GetLabels()
	if labels == nil {
		return false
	}
	_, ok := labels[label]
	return ok
}

func DeleteLabel(resource metav1.Object, label string) {
	labels := resource.GetLabels()
	if labels == nil {
		return
	}
	delete(labels, label)
	resource.SetLabels(labels)
}

func GetLabel(resource metav1.Object, label string) (value string, ok bool) {
	labels := resource.GetLabels()
	if labels == nil {
		return "", false
	}
	value, ok = labels[label]
	return
}

func SetLabels(resource metav1.Object, labels map[string]string) {
	existingLabels := resource.GetLabels()
	if existingLabels == nil {
		existingLabels = map[string]string{}
	}
	for key, value := range labels {
		existingLabels[key] = value
	}
	resource.SetLabels(existingLabels)
}

func SetLabel(resource metav1.Object, label, value string) {
	SetLabels(resource, map[string]string{label: value})
}

func HasAnnotation(resource metav1.Object, annotation string) bool {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		return false
	}
	_, ok := annotations[annotation]
	return ok
}

func DeleteAnnotation(resource metav1.Object, annotation string) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		return
	}
	delete(annotations, annotation)
	resource.SetAnnotations(annotations)
}

func GetAnnotation(resource metav1.Object, annotation string) (value string, ok bool) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		return "", false
	}
	value, ok = annotations[annotation]
	return
}

func SetAnnotation(resource metav1.Object, annotation, value string) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[annotation] = value
	resource.SetAnnotations(annotations)
}

var initOperatorNamespace sync.Once
var operatorNamespace string

// GetOperatorNamespace initializes and caches this operator's namespace; panics on error
func GetOperatorNamespace() string {
	initOperatorNamespace.Do(func() {
		var err error
		if operatorNamespace, err = k8sutil.GetOperatorNamespace(); err != nil {
			if err == k8sutil.ErrNoNamespace {
				// see if dev is manually specifying this during debugging
				if operatorNamespace = os.Getenv("POD_NAMESPACE"); operatorNamespace != "" {
					return
				}
			}
			panic(err)
		}
	})
	return operatorNamespace
}
