package controller

import (
	"github.com/openshift-knative/knative-kafka-operator/pkg/controller/knativeeventingkafka"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, knativeeventingkafka.Add)
}
