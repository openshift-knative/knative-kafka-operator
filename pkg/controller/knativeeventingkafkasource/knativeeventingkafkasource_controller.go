package knativeeventingkafkasource

import (
	"context"
	"flag"
	"github.com/operator-framework/operator-sdk/pkg/predicate"

	mf "github.com/jcrossley3/manifestival"
	operatorv1alpha1 "github.com/openshift-knative/knative-kafka-operator/pkg/apis/operator/v1alpha1"
	"github.com/openshift-knative/knative-kafka-operator/version"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	//"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	sourceFilename  = flag.String("source-filename", "deploy/resources/source", "The filename containing the YAML resources to apply")
	sourceRecursive = flag.Bool("source-recursive", false, "If filename is a directory, process all manifests recursively")
	log             = logf.Log.WithName("controller_knativeeventingkafkasource")
)

// Add creates a new KnativeEventingKafkaSource Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKnativeEventingKafkaSource{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("knativeeventingkafkasource-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KnativeEventingKafkaSource
	err = c.Watch(&source.Kind{Type: &operatorv1alpha1.KnativeEventingKafkaSource{}}, &handler.EnqueueRequestForObject{}, predicate.GenerationChangedPredicate{})
	if err != nil {
		return err
	}

	// Watch child deployments for availability
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.KnativeEventingKafkaSource{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileKnativeEventingKafkaSource{}

// ReconcileKnativeEventingKafkaSource reconciles a KnativeEventingKafkaSource object
type ReconcileKnativeEventingKafkaSource struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	config mf.Manifest
}

// Create manifestival resources and KnativeEventingKafkaSource, if necessary
func (r *ReconcileKnativeEventingKafkaSource) InjectClient(c client.Client) error {
	m, err := mf.NewManifest(*sourceFilename, *sourceRecursive, c)
	if err != nil {
		return err
	}
	r.config = m
	return nil
}

// Reconcile reads that state of the cluster for a KnativeEventingKafkaSource object and makes changes based on the state read
// and what is in the KnativeEventingKafkaSource.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKnativeEventingKafkaSource) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KnativeEventingKafkaSource")

	// Fetch the KnativeEventingKafkaSource instance
	instance := &operatorv1alpha1.KnativeEventingKafkaSource{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			r.config.DeleteAll()
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// stages hook for future work (e.g. deleteObsoleteResources)
	stages := []func(*operatorv1alpha1.KnativeEventingKafkaSource) error{
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
func (r *ReconcileKnativeEventingKafkaSource) initStatus(instance *operatorv1alpha1.KnativeEventingKafkaSource) error {
	if len(instance.Status.Conditions) == 0 {
		operatorv1alpha1.InitializeConditions(&instance.Status)
		if err := r.updateStatus(instance); err != nil {
			return err
		}
	}
	return nil
}

// Update the status subresource
func (r *ReconcileKnativeEventingKafkaSource) updateStatus(instance *operatorv1alpha1.KnativeEventingKafkaSource) error {

	// Account for https://github.com/kubernetes-sigs/controller-runtime/issues/406
	gvk := instance.GroupVersionKind()
	defer instance.SetGroupVersionKind(gvk)

	if err := r.client.Status().Update(context.TODO(), instance); err != nil {
		return err
	}
	return nil
}

// Apply the embedded resources
func (r *ReconcileKnativeEventingKafkaSource) install(instance *operatorv1alpha1.KnativeEventingKafkaSource) error {
	// Transform resources as appropriate
	fns := []mf.Transformer{
		mf.InjectOwner(instance),
		mf.InjectNamespace(instance.GetNamespace()),
		// TODO: probably not necessary
		// addSCCforSpecialClusterRoles,
	}
	r.config.Transform(fns...)

	if operatorv1alpha1.IsDeploying(&instance.Status) {
		return nil
	}
	defer r.updateStatus(instance)

	// Apply the resources in the YAML file
	if err := r.config.ApplyAll(); err != nil {
		operatorv1alpha1.MarkInstallFailed(&instance.Status, err.Error())
		return err
	}

	// Update status
	instance.Status.Version = version.Version
	operatorv1alpha1.MarkInstallSucceeded(&instance.Status)
	log.Info("Install succeeded", "version", version.Version)
	return nil
}

// Check for all deployments available
// TODO: what about statefulsets?
func (r *ReconcileKnativeEventingKafkaSource) checkDeployments(instance *operatorv1alpha1.KnativeEventingKafkaSource) error {
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
				operatorv1alpha1.MarkDeploymentsNotReady(&instance.Status)
				if errors.IsNotFound(err) {
					return nil
				}
				return err
			}
			if !available(deployment) {
				operatorv1alpha1.MarkDeploymentsNotReady(&instance.Status)
				return nil
			}
		}
	}
	operatorv1alpha1.MarkDeploymentsAvailable(&instance.Status)
	log.Info("All deployments are available")
	return nil
}
