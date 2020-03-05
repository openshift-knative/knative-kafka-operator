package v1alpha1

import (
	"github.com/knative/pkg/apis"
)

// GetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaSourceStatus) GetConditions() apis.Conditions {
	return s.Conditions
}

// SetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaSourceStatus) SetConditions(c apis.Conditions) {
	s.Conditions = c
}

func (is *KnativeEventingKafkaSourceStatus) IsReady() bool {
	return conditions.Manage(is).IsHappy()
}

func (is *KnativeEventingKafkaSourceStatus) IsInstalled() bool {
	return is.GetCondition(InstallSucceeded).IsTrue()
}

func (is *KnativeEventingKafkaSourceStatus) IsAvailable() bool {
	return is.GetCondition(DeploymentsAvailable).IsTrue()
}

func (is *KnativeEventingKafkaSourceStatus) IsDeploying() bool {
	return is.IsInstalled() && !is.IsAvailable()
}

func (is *KnativeEventingKafkaSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return conditions.Manage(is).GetCondition(t)
}

func (is *KnativeEventingKafkaSourceStatus) InitializeConditions() {
	conditions.Manage(is).InitializeConditions()
}

func (is *KnativeEventingKafkaSourceStatus) MarkInstallFailed(msg string) {
	conditions.Manage(is).MarkFalse(
		InstallSucceeded,
		"Error",
		"Install failed with message: %s", msg)
}

func (is *KnativeEventingKafkaSourceStatus) MarkInstallSucceeded() {
	conditions.Manage(is).MarkTrue(InstallSucceeded)
}

func (is *KnativeEventingKafkaSourceStatus) MarkDeploymentsAvailable() {
	conditions.Manage(is).MarkTrue(DeploymentsAvailable)
}

func (is *KnativeEventingKafkaSourceStatus) MarkDeploymentsNotReady() {
	conditions.Manage(is).MarkFalse(
		DeploymentsAvailable,
		"NotReady",
		"Waiting on deployments")
}
