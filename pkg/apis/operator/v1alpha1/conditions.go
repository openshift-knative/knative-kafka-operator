package v1alpha1

import "github.com/knative/pkg/apis"

func IsInstalled(is apis.ConditionsAccessor) bool {
	return GetCondition(is, InstallSucceeded).IsTrue()
}

func IsAvailable(is apis.ConditionsAccessor) bool {
	return GetCondition(is, DeploymentsAvailable).IsTrue()
}

func IsDeploying(is apis.ConditionsAccessor) bool {
	return IsInstalled(is) && !IsAvailable(is)
}

func GetCondition(is apis.ConditionsAccessor, t apis.ConditionType) *apis.Condition {
	return conditions.Manage(is).GetCondition(t)
}

func InitializeConditions(is apis.ConditionsAccessor) {
	conditions.Manage(is).InitializeConditions()
}

func MarkInstallFailed(is apis.ConditionsAccessor, msg string) {
	conditions.Manage(is).MarkFalse(
		InstallSucceeded,
		"Error",
		"Install failed with message: %s", msg)
}

func MarkInstallSucceeded(is apis.ConditionsAccessor) {
	conditions.Manage(is).MarkTrue(InstallSucceeded)
}

func MarkDeploymentsAvailable(is apis.ConditionsAccessor) {
	conditions.Manage(is).MarkTrue(DeploymentsAvailable)
}

func MarkDeploymentsNotReady(is apis.ConditionsAccessor) {
	conditions.Manage(is).MarkFalse(
		DeploymentsAvailable,
		"NotReady",
		"Waiting on deployments")
}
