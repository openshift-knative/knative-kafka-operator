package v1alpha1

import "github.com/knative/pkg/apis"

const (
	InstallSucceeded     apis.ConditionType = "InstallSucceeded"
	DeploymentsAvailable apis.ConditionType = "DeploymentsAvailable"
)

var conditions = apis.NewLivingConditionSet(
	DeploymentsAvailable,
	InstallSucceeded,
)
