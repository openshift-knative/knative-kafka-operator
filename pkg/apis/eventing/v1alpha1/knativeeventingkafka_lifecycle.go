package v1alpha1

import (
	"github.com/knative/pkg/apis"
)

var conditions = apis.NewLivingConditionSet(
	DeploymentsAvailable,
	InstallSucceeded,
)

// GetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaStatus) GetConditions() apis.Conditions {
	return s.Conditions
}

// SetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaStatus) SetConditions(c apis.Conditions) {
	s.Conditions = c
}

func (is *KnativeEventingKafkaStatus) IsReady() bool {
	return conditions.Manage(is).IsHappy()
}

func (is *KnativeEventingKafkaStatus) IsInstalled() bool {
	return is.GetCondition(InstallSucceeded).IsTrue()
}

func (is *KnativeEventingKafkaStatus) IsAvailable() bool {
	return is.GetCondition(DeploymentsAvailable).IsTrue()
}

func (is *KnativeEventingKafkaStatus) IsDeploying() bool {
	return is.IsInstalled() && !is.IsAvailable()
}

func (is *KnativeEventingKafkaStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return conditions.Manage(is).GetCondition(t)
}

func (is *KnativeEventingKafkaStatus) InitializeConditions() {
	conditions.Manage(is).InitializeConditions()
}

func (is *KnativeEventingKafkaStatus) MarkInstallFailed(msg string) {
	conditions.Manage(is).MarkFalse(
		InstallSucceeded,
		"Error",
		"Install failed with message: %s", msg)
}

func (is *KnativeEventingKafkaStatus) MarkInstallSucceeded() {
	conditions.Manage(is).MarkTrue(InstallSucceeded)
}

func (is *KnativeEventingKafkaStatus) MarkDeploymentsAvailable() {
	conditions.Manage(is).MarkTrue(DeploymentsAvailable)
}

func (is *KnativeEventingKafkaStatus) MarkDeploymentsNotReady() {
	conditions.Manage(is).MarkFalse(
		DeploymentsAvailable,
		"NotReady",
		"Waiting on deployments")
}

func (is *KnativeEventingKafkaStatus) MarkIgnored(msg string) {
	conditions.Manage(is).MarkFalse(
		InstallSucceeded,
		"Ignored",
		"Install not attempted: %s", msg)
}
