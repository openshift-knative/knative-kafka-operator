package v1alpha1

import (
	"github.com/knative/pkg/apis"
)

var conditions = apis.NewLivingConditionSet(
	DeploymentsAvailable,
	InstallSucceeded,
)

// GetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaChannelStatus) GetConditions() apis.Conditions {
	return s.Conditions
}

// SetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaChannelStatus) SetConditions(c apis.Conditions) {
	s.Conditions = c
}

func (is *KnativeEventingKafkaChannelStatus) IsReady() bool {
	return conditions.Manage(is).IsHappy()
}

func (is *KnativeEventingKafkaChannelStatus) IsInstalled() bool {
	return is.GetCondition(InstallSucceeded).IsTrue()
}

func (is *KnativeEventingKafkaChannelStatus) IsAvailable() bool {
	return is.GetCondition(DeploymentsAvailable).IsTrue()
}

func (is *KnativeEventingKafkaChannelStatus) IsDeploying() bool {
	return is.IsInstalled() && !is.IsAvailable()
}

func (is *KnativeEventingKafkaChannelStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return conditions.Manage(is).GetCondition(t)
}

func (is *KnativeEventingKafkaChannelStatus) InitializeConditions() {
	conditions.Manage(is).InitializeConditions()
}

func (is *KnativeEventingKafkaChannelStatus) MarkInstallFailed(msg string) {
	conditions.Manage(is).MarkFalse(
		InstallSucceeded,
		"Error",
		"Install failed with message: %s", msg)
}

func (is *KnativeEventingKafkaChannelStatus) MarkInstallSucceeded() {
	conditions.Manage(is).MarkTrue(InstallSucceeded)
}

func (is *KnativeEventingKafkaChannelStatus) MarkDeploymentsAvailable() {
	conditions.Manage(is).MarkTrue(DeploymentsAvailable)
}

func (is *KnativeEventingKafkaChannelStatus) MarkDeploymentsNotReady() {
	conditions.Manage(is).MarkFalse(
		DeploymentsAvailable,
		"NotReady",
		"Waiting on deployments")
}
