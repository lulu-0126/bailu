// Copyright @ 2021 Lenovo. All rights reserved
// Confidential and Restricted
// This file was generated using one of kubebuilder templated ( https://github.com/kubernetes-sigs/kubebuilder/tree/master/pkg/plugins/golang/v2/scaffolds/internal/templates) distributed under the terms of Apache-2.0 license.

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	ErrorMaintenanceAndPower       = "maintenance and Power occured at the same time"
	ErrorMaintenanceAndPowerStatus = "maintenance status isn't done, power isn't on"
	ErrorCreateForbidden           = "create isn't allowed"
	ErrorDeleteForbidden           = "delete isn't allowed"
)

// log is for logging in this package.
var nodemanagementlog = logf.Log.WithName("nodemanagement-resource")

func (r *NodeManagement) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-cloud-edge-v1-nodemanagement,mutating=true,failurePolicy=fail,sideEffects=None,groups=cloud.edge,resources=nodemanagements,verbs=create;update,versions=v1,name=mnodemanagement.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &NodeManagement{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *NodeManagement) Default() {
	nodemanagementlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-cloud-edge-v1-nodemanagement,mutating=false,failurePolicy=fail,sideEffects=None,groups=cloud.edge,resources=nodemanagements,verbs=create;update,versions=v1,name=vnodemanagement.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &NodeManagement{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *NodeManagement) ValidateCreate() error {
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *NodeManagement) ValidateUpdate(old runtime.Object) error {
	return nil
}

func contains(tgt NodeStatusType, arr []NodeStatusType) bool {
	for _, e := range arr {
		if tgt == e {
			return true
		}
	}
	return false
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *NodeManagement) ValidateDelete() error {
	return nil
}
