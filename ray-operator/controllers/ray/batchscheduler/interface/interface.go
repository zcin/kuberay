package schedulerinterface

import (
	rayv1 "github.com/ray-project/kuberay/ray-operator/apis/ray/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

// BatchScheduler manages submitting RayCluster pods to a third-party scheduler.
type BatchScheduler interface {
	// Name corresponds to the schedulerName in Kubernetes:
	// https://kubernetes.io/docs/tasks/extend-kubernetes/configure-multiple-schedulers/
	Name() string

	// DoBatchSchedulingOnSubmission handles submitting the RayCluster to the batch scheduler on creation / update
	// For most batch schedulers, this results in the creation of a PodGroup.
	DoBatchSchedulingOnSubmission(app *rayv1.RayCluster) error

	// AddMetadataToPod enriches Pod specs with metadata necessary to tie them to the scheduler.
	// For example, setting labels for queues / priority, and setting schedulerName.
	AddMetadataToPod(app *rayv1.RayCluster, groupName string, pod *corev1.Pod)
}

// BatchSchedulerFactory handles initial setup of the scheduler plugin by registering the
// necessary callbacks with the operator, and the creation of the BatchScheduler itself.
type BatchSchedulerFactory interface {
	// New creates a new BatchScheduler for the scheduler plugin.
	New(config *rest.Config) (BatchScheduler, error)

	// AddToScheme adds the types in this scheduler to the given scheme (runs during init).
	AddToScheme(scheme *runtime.Scheme)

	// ConfigureReconciler configures the RayCluster Reconciler in the process of being built by
	// adding watches for its scheduler-specific custom resource types, and any other needed setup.
	ConfigureReconciler(b *builder.Builder) *builder.Builder
}

type DefaultBatchScheduler struct{}

type DefaultBatchSchedulerFactory struct{}

func GetDefaultPluginName() string {
	return "default"
}

func (d *DefaultBatchScheduler) Name() string {
	return GetDefaultPluginName()
}

func (d *DefaultBatchScheduler) DoBatchSchedulingOnSubmission(app *rayv1.RayCluster) error {
	return nil
}

func (d *DefaultBatchScheduler) AddMetadataToPod(app *rayv1.RayCluster, groupName string, pod *corev1.Pod) {
}

func (df *DefaultBatchSchedulerFactory) New(config *rest.Config) (BatchScheduler, error) {
	return &DefaultBatchScheduler{}, nil
}

func (df *DefaultBatchSchedulerFactory) AddToScheme(scheme *runtime.Scheme) {
}

func (df *DefaultBatchSchedulerFactory) ConfigureReconciler(b *builder.Builder) *builder.Builder {
	return b
}
