package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/controller"

	api "github.com/gaocegege/kubetidb/pkg/apis/tidb/v1alpha1"
	clientset "github.com/gaocegege/kubetidb/pkg/clientset/versioned"
	tidbscheme "github.com/gaocegege/kubetidb/pkg/clientset/versioned/scheme"
	informers "github.com/gaocegege/kubetidb/pkg/informers/externalversions"
	listers "github.com/gaocegege/kubetidb/pkg/listers/tidb/v1alpha1"
)

const (
	controllerName = "kubetidb"
	// SuccessSynced is used as part of the Event 'reason' when a TiDBCluster is synced
	SuccessSynced = "Synced"

	// MessageResourceSynced is the message used for an Event fired when a TiDBCluster
	// is synced successfully
	MessageResourceSynced = "TiDB synced successfully"
)

// Controller is the type for TiDBCluster controller.
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// tidbClientset is a clientset for our own API group
	tidbClientset clientset.Interface
	tidbLister    listers.TiDBClusterLister
	tidbSynced    cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder

	// A TTLCache of tidb creates/deletes each rc expects to see
	expectations controller.ControllerExpectationsInterface
}

// NewController returns a new tfJob controller.
func NewController(
	kubeclientset kubernetes.Interface,
	tidbClientset clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	tidbInformerFactory informers.SharedInformerFactory) *Controller {

	// obtain references to shared index informers for the tfJob type
	tidbInformer := tidbInformerFactory.Kubetidb().V1alpha1().TiDBClusters()

	// Create event broadcaster
	// Add tfJob-controller types to the default Kubernetes Scheme so Events can be
	// logged for tfJob-controller types.
	tidbscheme.AddToScheme(scheme.Scheme)
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: controllerName})

	controller := &Controller{
		kubeclientset: kubeclientset,
		tidbClientset: tidbClientset,
		tidbLister:    tidbInformer.Lister(),
		expectations:  controller.NewControllerExpectations(),
		tidbSynced:    tidbInformer.Informer().HasSynced,
		workqueue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "tfJobs"),
		recorder:      recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when tfJob resources change
	tidbInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addTiDBCluster,
		UpdateFunc: controller.updateTiDBCluster,
		DeleteFunc: controller.deleteTiDBCluster,
	})

	controller.tidbLister = tidbInformer.Lister()

	return controller
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the TFJob resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	if len(namespace) == 0 || len(name) == 0 {
		return fmt.Errorf("invalid job key %q: either namespace or name is missing", key)
	}

	needsSync := c.expectations.SatisfiedExpectations(key)

	// Get the TFJob resource with this namespace/name
	tidbCluster, err := c.tidbLister.TiDBClusters(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			glog.V(4).Infof("Job has been deleted: %v", key)
			c.expectations.DeleteExpectations(key)
			return nil
		}
		return err
	}

	if needsSync {
		c.syncCluster(tidbCluster)
	}

	return nil
}

func (c *Controller) syncCluster(tidbCluster *api.TiDBCluster) {
	glog.V(4).Infof("Sync TiDBCluster: %v", *tidbCluster)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting tidb controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.tidbSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process tfjob resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// TiDBCluster resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) addTiDBCluster(obj interface{}) {
	c.enqueueTiDBCluster(obj)
}

func (c *Controller) updateTiDBCluster(old, new interface{}) {
	glog.Info("Update TiDB cluster")
	newCluster := new.(*api.TiDBCluster)
	oldCluster := old.(*api.TiDBCluster)
	if newCluster.ResourceVersion == oldCluster.ResourceVersion {
		glog.Infof("ResourceVersion not changed: %s", newCluster.ResourceVersion)
		// Periodic resync will send update events for all known tfJobes.
		// Two different versions of the same tfJob will always have different RVs.
		return
	}
	c.enqueueTiDBCluster(newCluster)
}

func (c *Controller) deleteTiDBCluster(obj interface{}) {
	glog.Errorln("To Be Implemented.")
}

func (c *Controller) enqueueTiDBCluster(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
