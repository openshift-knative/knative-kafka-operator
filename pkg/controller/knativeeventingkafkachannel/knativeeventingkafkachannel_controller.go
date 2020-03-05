package knativeeventingkafkachannel

import (
	"context"
	go_errors "errors"
	"flag"

	"github.com/operator-framework/operator-sdk/pkg/predicate"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	mf "github.com/jcrossley3/manifestival"
	operatorv1alpha1 "github.com/openshift-knative/knative-kafka-operator/pkg/apis/operator/v1alpha1"
	"github.com/openshift-knative/knative-kafka-operator/version"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	channelFilename  = flag.String("channel-filename", "deploy/resources/channel", "The filename containing the YAML resources to apply")
	channelRecursive = flag.Bool("channel-recursive", false, "If filename is a directory, process all manifests recursively")
	log              = logf.Log.WithName("controller_knativeeventingkafkachannel")
)

// Add creates a new KnativeEventingKafkaChannel Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKnativeEventingKafkaChannel{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("knativeeventingkafkachannel-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KnativeEventingKafkaChannel
	err = c.Watch(&source.Kind{Type: &operatorv1alpha1.KnativeEventingKafkaChannel{}}, &handler.EnqueueRequestForObject{}, predicate.GenerationChangedPredicate{})
	if err != nil {
		return err
	}

	// Watch child deployments for availability
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.KnativeEventingKafkaChannel{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileKnativeEventingKafkaChannel{}

// ReconcileKnativeEventingKafkaChannel reconciles a KnativeEventingKafkaChannel object
type ReconcileKnativeEventingKafkaChannel struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	config mf.Manifest
}

// Create manifestival resources and ReconcileKnativeEventingKafkaChannel, if necessary
func (r *ReconcileKnativeEventingKafkaChannel) InjectClient(c client.Client) error {
	m, err := mf.NewManifest(*channelFilename, *channelRecursive, c)
	if err != nil {
		return err
	}
	r.config = m
	return nil
}

// Reconcile reads that state of the cluster for a KnativeEventingKafkaChannel object and makes changes based on the state read
// and what is in the KnativeEventingKafkaChannel.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKnativeEventingKafkaChannel) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KnativeEventingKafkaChannel")

	// Fetch the KnativeEventingKafkaChannel instance
	instance := &operatorv1alpha1.KnativeEventingKafkaChannel{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			r.config.DeleteAll()
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// stages hook for future work (e.g. deleteObsoleteResources)
	stages := []func(*operatorv1alpha1.KnativeEventingKafkaChannel) error{
		r.initStatus,
		r.install,
		r.checkDeployments,
	}

	for _, stage := range stages {
		if err := stage(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// Initialize status conditions
func (r *ReconcileKnativeEventingKafkaChannel) initStatus(instance *operatorv1alpha1.KnativeEventingKafkaChannel) error {
	if len(instance.Status.Conditions) == 0 {
		instance.Status.InitializeConditions()
		if err := r.updateStatus(instance); err != nil {
			return err
		}
	}
	return nil
}

// Update the status subresource
func (r *ReconcileKnativeEventingKafkaChannel) updateStatus(instance *operatorv1alpha1.KnativeEventingKafkaChannel) error {

	// Account for https://github.com/kubernetes-sigs/controller-runtime/issues/406
	gvk := instance.GroupVersionKind()
	defer instance.SetGroupVersionKind(gvk)

	if err := r.client.Status().Update(context.TODO(), instance); err != nil {
		return err
	}
	return nil
}

// Apply the embedded resources
func (r *ReconcileKnativeEventingKafkaChannel) install(instance *operatorv1alpha1.KnativeEventingKafkaChannel) error {
	// Transform resources as appropriate
	fns := []mf.Transformer{
		mf.InjectOwner(instance),
		mf.InjectNamespace(instance.GetNamespace()),
		// TODO: probably not necessary
		//addSCCforSpecialClusterRoles,
		bootstrapServersTransformer(instance.Spec.BootstrapServers),
	}
	r.config.Transform(fns...)

	if instance.Status.IsDeploying() {
		return nil
	}
	defer r.updateStatus(instance)

	// Apply the resources in the YAML file
	if err := r.config.ApplyAll(); err != nil {
		instance.Status.MarkInstallFailed(err.Error())
		return err
	}

	if err := r.setAsDefaultChannel(instance.Spec.SetAsDefaultChannelProvisioner); err != nil {
		return err
	}

	// Update status
	instance.Status.Version = version.Version
	instance.Status.MarkInstallSucceeded()
	log.Info("Install succeeded", "version", version.Version)
	return nil
}

// Check for all deployments available
// TODO: what about statefulsets?
func (r *ReconcileKnativeEventingKafkaChannel) checkDeployments(instance *operatorv1alpha1.KnativeEventingKafkaChannel) error {
	defer r.updateStatus(instance)
	available := func(d *appsv1.Deployment) bool {
		for _, c := range d.Status.Conditions {
			if c.Type == appsv1.DeploymentAvailable && c.Status == corev1.ConditionTrue {
				return true
			}
		}
		return false
	}
	deployment := &appsv1.Deployment{}
	for _, u := range r.config.Resources {
		if u.GetKind() == "Deployment" {
			key := client.ObjectKey{Namespace: u.GetNamespace(), Name: u.GetName()}
			if err := r.client.Get(context.TODO(), key, deployment); err != nil {
				instance.Status.MarkDeploymentsNotReady()
				if errors.IsNotFound(err) {
					return nil
				}
				return err
			}
			if !available(deployment) {
				instance.Status.MarkDeploymentsNotReady()
				return nil
			}
		}
	}
	instance.Status.MarkDeploymentsAvailable()
	log.Info("All deployments are available")
	return nil
}

func (r *ReconcileKnativeEventingKafkaChannel) setAsDefaultChannel(doSet bool) error {
	key := client.ObjectKey{Namespace: "knative-eventing", Name: "default-ch-webhook"}
	result := &unstructured.Unstructured{}
	result.SetAPIVersion("v1")
	result.SetKind("ConfigMap")
	if err := r.client.Get(context.TODO(), key, result); err != nil {
		log.Error(err, "Unable to set or unset KafkaChannel as the default channel.")
		return err
	}
	configMapData := result.Object["data"]
	configMap, ok := configMapData.(map[string]interface{})
	if !ok {
		return go_errors.New("Unexpected structure of knative-eventing/default-ch-webhook ConfigMap")
	}
	defaultChannelConfigValue := "clusterDefault:\n  apiversion: messaging.knative.dev/v1alpha1\n  kind: "
	if doSet {
		defaultChannelConfigValue += `KafkaChannel`
	} else {
		defaultChannelConfigValue += `InMemoryChannel`
	}
	defaultChannelConfigValue += "\n"
	configMap["default-ch-config"] = defaultChannelConfigValue
	err := r.config.Apply(result)
	return err
}

// bootstrapServersTransformer modifies the configmap for Knative's Kafka channel
func bootstrapServersTransformer(bootstrapServers string) mf.Transformer {
	return func(u *unstructured.Unstructured) error {
		if u.GetKind() == "ConfigMap" && u.GetName() == "config-kafka" {
			unstructured.SetNestedField(u.Object, bootstrapServers, "data", "bootstrapServers")
		}
		return nil
	}
}
