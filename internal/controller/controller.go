package controller

import (
	"context"
	"log"
	"time"

	"heartbeat-operator/api/v1alpha1"
	"heartbeat-operator/internal/config"
	"heartbeat-operator/internal/metrics"
	"heartbeat-operator/internal/prober"
	"heartbeat-operator/internal/ui"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

type ProbeController struct {
	client    *kubernetes.Clientset
	crdClient *CrdClient
	rule      config.ProbeRule
	probe     prober.Prober
	recorder  record.EventRecorder
}

func New(client *kubernetes.Clientset, crdClient *CrdClient, rule config.ProbeRule, p prober.Prober, recorder record.EventRecorder) *ProbeController {
	return &ProbeController{
		client:    client,
		crdClient: crdClient,
		rule:      rule,
		probe:     p,
		recorder:  recorder,
	}
}

func (c *ProbeController) Start(ctx context.Context) {
	log.Printf("[%s] Started probing %s", c.rule.Name, c.rule.CheckTarget)

	if err := c.ensureCR(ctx); err != nil {
		log.Printf("[%s] Failed to ensure CRD (will retry): %v", c.rule.Name, err)
	}

	interval := config.ParseInterval(c.rule.Interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial check
	c.reconcile(ctx)

	for {
		select {
		case <-ticker.C:
			c.reconcile(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (c *ProbeController) ensureCR(ctx context.Context) error {
	if _, err := c.crdClient.Get(ctx, c.rule.Name); err == nil {
		return nil
	}

	// Create if not exists
	dc := &v1alpha1.Probe{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.rule.Name,
			Namespace: c.rule.Namespace,
		},
		Spec: v1alpha1.ProbeSpec{
			CheckType:   c.rule.CheckType,
			CheckTarget: c.rule.CheckTarget,
			Interval:    c.rule.Interval,
		},
	}
	log.Printf("[%s] Creating Probe CR...", c.rule.Name)
	_, err := c.crdClient.Create(ctx, dc)
	return err
}

func (c *ProbeController) reconcile(ctx context.Context) {
	start := time.Now()
	isHealthy := c.probe.Check()
	duration := time.Since(start).Seconds()

	metrics.ProbeDuration.WithLabelValues(c.rule.Name, c.rule.CheckTarget, c.rule.CheckType).Observe(duration)
	metrics.ProbeLastTimestamp.WithLabelValues(c.rule.Name, c.rule.CheckTarget, c.rule.CheckType).Set(float64(time.Now().Unix()))

	if isHealthy {
		metrics.ProbeSuccess.WithLabelValues(c.rule.Name, c.rule.CheckTarget, c.rule.CheckType).Set(1)
	} else {
		metrics.ProbeSuccess.WithLabelValues(c.rule.Name, c.rule.CheckTarget, c.rule.CheckType).Set(0)
	}

	ui.UpdateState(c.rule.Name, c.rule.CheckTarget, c.rule.CheckType, isHealthy)

	cr, err := c.crdClient.Get(ctx, c.rule.Name)
	if err != nil {
		if err := c.ensureCR(ctx); err != nil {
			log.Printf("[%s] CR missing and failed to create: %v", c.rule.Name, err)
			return
		}
		cr, err = c.crdClient.Get(ctx, c.rule.Name)
		if err != nil {
			log.Printf("[%s] Failed to get CR: %v", c.rule.Name, err)
			return
		}
	}

	// Update status logic
	// Only update if changed or if it's been a while?
	// For now, simple update
	now := metav1.Now()
	msg := "Check passed"
	if !isHealthy {
		msg = "Check failed"
	}

	if cr.Status.Healthy != isHealthy || cr.Status.Message != msg {
		// Emit Event
		eventType := corev1.EventTypeNormal
		reason := "ProbeHealthy"
		if !isHealthy {
			eventType = corev1.EventTypeWarning
			reason = "ProbeFailing"
		}
		c.recorder.Eventf(cr, eventType, reason, "Check target %s status changed to healthy=%v", c.rule.CheckTarget, isHealthy)

		cr.Status.Healthy = isHealthy
		cr.Status.Message = msg
		cr.Status.LastProbeTime = &now
		_, err := c.crdClient.UpdateStatus(ctx, cr)
		if err != nil {
			log.Printf("[%s] Failed to update CR status: %v", c.rule.Name, err)
		} else {
			log.Printf("[%s] Updated CR status: healthy=%v", c.rule.Name, isHealthy)
		}
	} else if cr.Status.LastProbeTime == nil || time.Since(cr.Status.LastProbeTime.Time) > time.Minute {
		cr.Status.LastProbeTime = &now
		if _, err := c.crdClient.UpdateStatus(ctx, cr); err != nil {
			log.Printf("[%s] Failed to update CR timestamp: %v", c.rule.Name, err)
		}
	}
}
