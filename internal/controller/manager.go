package controller

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"heartbeat-operator/api/v1alpha1"
	"heartbeat-operator/internal/metrics"
	"heartbeat-operator/internal/prober"
	"heartbeat-operator/internal/ui"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

type ProbeManager struct {
	client    *kubernetes.Clientset
	crdClient *CrdClient
	recorder  record.EventRecorder

	informer     cache.SharedIndexInformer
	activeProbes map[string]context.CancelFunc
	mu           sync.Mutex
}

func NewManager(client *kubernetes.Clientset, crdClient *CrdClient, recorder record.EventRecorder) *ProbeManager {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return crdClient.List(context.Background(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return crdClient.Watch(context.Background(), options)
		},
	}

	informer := cache.NewSharedIndexInformer(
		lw,
		&v1alpha1.Probe{},
		5*time.Minute,
		cache.Indexers{},
	)

	m := &ProbeManager{
		client:       client,
		crdClient:    crdClient,
		recorder:     recorder,
		informer:     informer,
		activeProbes: make(map[string]context.CancelFunc),
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    m.onAdd,
		UpdateFunc: m.onUpdate,
		DeleteFunc: m.onDelete,
	})

	return m
}

func (m *ProbeManager) Start(ctx context.Context) {
	log.Println("Starting ProbeManager Informer...")
	go m.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), m.informer.HasSynced) {
		log.Println("Failed to sync cache")
		return
	}
	log.Println("Cache synced. Ready to process probes.")
	<-ctx.Done()
	log.Println("Stopping ProbeManager...")
}

func (m *ProbeManager) onAdd(obj interface{}) {
	cr := obj.(*v1alpha1.Probe)
	m.startProbe(cr)
}

func (m *ProbeManager) onUpdate(oldObj, newObj interface{}) {
	newCr := newObj.(*v1alpha1.Probe)
	oldCr := oldObj.(*v1alpha1.Probe)

	if oldCr.Spec.Interval != newCr.Spec.Interval ||
		oldCr.Spec.CheckTarget != newCr.Spec.CheckTarget ||
		oldCr.Spec.CheckType != newCr.Spec.CheckType ||
		oldCr.Spec.Timeout != newCr.Spec.Timeout {
		m.stopProbe(oldCr)
		m.startProbe(newCr)
	}
}

func (m *ProbeManager) onDelete(obj interface{}) {
	cr, ok := obj.(*v1alpha1.Probe)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return
		}
		cr, ok = tombstone.Obj.(*v1alpha1.Probe)
		if !ok {
			return
		}
	}
	m.stopProbe(cr)
}

func (m *ProbeManager) startProbe(cr *v1alpha1.Probe) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s/%s", cr.Namespace, cr.Name)
	if _, ok := m.activeProbes[key]; ok {
		return
	}

	log.Printf("[%s] Starting probe for %s", key, cr.Spec.CheckTarget)

	var p prober.Prober
	timeout, err := time.ParseDuration(cr.Spec.Timeout)
	if err != nil || timeout == 0 {
		timeout = 2 * time.Second
	}

	switch cr.Spec.CheckType {
	case "http":
		p = prober.NewHttpProber(cr.Spec.CheckTarget, timeout)
	case "tcp":
		p = prober.NewTcpProber(cr.Spec.CheckTarget, timeout)
	case "exec":
		p = prober.NewExecProber(cr.Spec.CheckTarget)
	default:
		log.Printf("[%s] Unknown check type %s", key, cr.Spec.CheckType)
		return
	}

	interval, err := time.ParseDuration(cr.Spec.Interval)
	if err != nil || interval == 0 {
		interval = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.activeProbes[key] = cancel

	go m.runTicker(ctx, cr.Name, cr.Namespace, cr.Spec, p, interval)
}

func (m *ProbeManager) stopProbe(cr *v1alpha1.Probe) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s/%s", cr.Namespace, cr.Name)
	if cancel, ok := m.activeProbes[key]; ok {
		log.Printf("[%s] Stopping probe", key)
		cancel()
		delete(m.activeProbes, key)
	}
	
	ui.RemoveState(cr.Name)
}

func (m *ProbeManager) runTicker(ctx context.Context, name, namespace string, spec v1alpha1.ProbeSpec, p prober.Prober, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	m.executeCheck(ctx, name, namespace, spec, p)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.executeCheck(ctx, name, namespace, spec, p)
		}
	}
}

func (m *ProbeManager) executeCheck(ctx context.Context, name, namespace string, spec v1alpha1.ProbeSpec, p prober.Prober) {
	start := time.Now()
	isHealthy := p.Check()
	duration := time.Since(start).Seconds()

	metrics.ProbeDuration.WithLabelValues(name, spec.CheckTarget, spec.CheckType).Observe(duration)
	metrics.ProbeLastTimestamp.WithLabelValues(name, spec.CheckTarget, spec.CheckType).Set(float64(time.Now().Unix()))

	if isHealthy {
		metrics.ProbeSuccess.WithLabelValues(name, spec.CheckTarget, spec.CheckType).Set(1)
	} else {
		metrics.ProbeSuccess.WithLabelValues(name, spec.CheckTarget, spec.CheckType).Set(0)
	}

	ui.UpdateState(name, spec.CheckTarget, spec.CheckType, isHealthy)

	obj, exists, err := m.informer.GetIndexer().GetByKey(fmt.Sprintf("%s/%s", namespace, name))
	if err != nil || !exists {
		return
	}
	cr := obj.(*v1alpha1.Probe).DeepCopy()

	now := metav1.Now()
	msg := "Check passed"
	if !isHealthy {
		msg = "Check failed"
	}

	if cr.Status.Healthy != isHealthy || cr.Status.Message != msg {
		eventType := corev1.EventTypeNormal
		reason := "ProbeHealthy"
		if !isHealthy {
			eventType = corev1.EventTypeWarning
			reason = "ProbeFailing"
		}
		m.recorder.Eventf(cr, eventType, reason, "Check target %s status changed to healthy=%v", spec.CheckTarget, isHealthy)

		cr.Status.Healthy = isHealthy
		cr.Status.Message = msg
		cr.Status.LastProbeTime = &now
		if _, err := m.crdClient.UpdateStatus(context.Background(), cr); err != nil {
			log.Printf("[%s] Failed to update CR status: %v", name, err)
		} else {
			log.Printf("[%s] Updated CR status: healthy=%v", name, isHealthy)
		}
	} else if cr.Status.LastProbeTime == nil || time.Since(cr.Status.LastProbeTime.Time) > time.Minute {
		cr.Status.LastProbeTime = &now
		if _, err := m.crdClient.UpdateStatus(context.Background(), cr); err != nil {
			log.Printf("[%s] Failed to update CR timestamp: %v", name, err)
		}
	}
}
