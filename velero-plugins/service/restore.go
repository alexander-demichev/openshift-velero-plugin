package service

import (
	"encoding/json"

	"github.com/heptio/velero/pkg/plugin/velero"
	"github.com/sirupsen/logrus"
	corev1API "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// RestorePlugin is a restore item action plugin for Velero
type RestorePlugin struct {
	Log logrus.FieldLogger
}

// AppliesTo returns a velero.ResourceSelector that applies to services
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"services"},
	}, nil
}

// Execute action for the restore plugin for the service resource
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.Log.Info("[service-restore] Entering Service restore plugin")

	service := corev1API.Service{}
	itemMarshal, _ := json.Marshal(input.Item)
	json.Unmarshal(itemMarshal, &service)
	p.Log.Infof("[service-restore] service: %s", service.Name)

	// only clear ExternalIPs for LoadBalancer services
	if service.Spec.Type == corev1API.ServiceTypeLoadBalancer {
		p.Log.Infof("[service-restore] Clearing externalIPs for LoadBalancer service: %s", service.Name)
		service.Spec.ExternalIPs = nil
	}

	var out map[string]interface{}
	objrec, _ := json.Marshal(service)
	json.Unmarshal(objrec, &out)

	return velero.NewRestoreItemActionExecuteOutput(&unstructured.Unstructured{Object: out}), nil
}
