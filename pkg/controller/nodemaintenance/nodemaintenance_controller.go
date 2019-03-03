package nodemaintenance

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/kubectl/drain"
	kubevirtv1alpha3 "kubevirt.io/node-maintenance-operator/pkg/apis/kubevirt/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nodemaintenance")

// Add creates a new NodeMaintenance Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	r, err := newReconciler(mgr)
	if err != nil {
		return err
	}
	return add(mgr, r)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) (reconcile.Reconciler, error) {
	r := &ReconcileNodeMaintenance{client: mgr.GetClient(), scheme: mgr.GetScheme()}
	err := initDrainer(r, mgr.GetConfig())
	return r, err
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nodemaintenance-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NodeMaintenance
	err = c.Watch(&source.Kind{Type: &kubevirtv1alpha3.NodeMaintenance{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

func initDrainer(r *ReconcileNodeMaintenance, config *rest.Config) error {
	//Continue even if there are pods not managed by a ReplicationController, ReplicaSet, Job, DaemonSet or StatefulSet
	//r.drainer.Force = false
	//Ignore DaemonSet-managed pods.
	//r.drainer.IgnoreAllDaemonSets = false
	//Continue even if there are pods using emptyDir (local data that will be deleted when the node is drained).
	//r.drainer.DeleteLocalData = false
	//Period of time in seconds given to each pod to terminate gracefully. If negative, the default value specified in the pod will be used.
	r.drainer.GracePeriodSeconds = -1
	//The length of time to wait before giving up, zero means infinite
	//r.drainer.Timeout = int64(10)
	//Selector (label query) to filter on
	//r.drainer.Selector = "l"
	//Label selector to filter pods on the node
	//r.drainer.PodSelector = "pod-selector"

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	r.drainer.Client = cs
	r.drainer.DryRun = false

	return nil
}

var _ reconcile.Reconciler = &ReconcileNodeMaintenance{}

// ReconcileNodeMaintenance reconciles a NodeMaintenance object
type ReconcileNodeMaintenance struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client  client.Client
	scheme  *runtime.Scheme
	drainer *drain.Helper
}

// Reconcile reads that state of the cluster for a NodeMaintenance object and makes changes based on the state read
// and what is in the NodeMaintenance.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNodeMaintenance) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NodeMaintenance")

	// Fetch the NodeMaintenance instance
	instance := &kubevirtv1alpha3.NodeMaintenance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("The request object cannot be found.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Info("Error reading the request object, requeuing.")
		return reconcile.Result{}, err
	}

	nodeName := instance.Name

	reqLogger.Info(fmt.Sprintf("Applying Maintenance mode on Node: %s with Reason: %s", nodeName, instance.Spec.Reason))

	node := &corev1.Node{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: nodeName}, node)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Error(err, fmt.Sprintf("Node: %s cannot be found", nodeName))
		return reconcile.Result{}, err
	} else if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to get Node %s: %v\n", nodeName, err))
		return reconcile.Result{}, err
	}

	reqLogger.Info(fmt.Sprintf("Cordon Node: %s", nodeName))

	reqLogger.Info(fmt.Sprintf("Evict all Pods from Node: %s", nodeName))

	return reconcile.Result{}, nil
}
