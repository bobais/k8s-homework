package main

import (
	"github.com/bobais/eventstats-watcher/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"time"
)

var debug = utils.IsPositiveStatement(os.Getenv("DEBUG"))

func main() {
	cl, err := utils.GetResolvedKubernetesClient()
	if err != nil {
		log.Fatalf("Failed to initialize kubernetes client %v", err)
	}

	eventWatcher, err := cl.CoreV1().Events(metav1.NamespaceAll).Watch(utils.GetEventListingOptions())
	if err != nil {
		log.Fatalf("Failed to create event watcher %v", err)
	}

	// Runs in goroutine
	utils.StartHttpServer()

	var (
		// Time to keep track when the loop started to skip portion of old events
		startOfLoop      = time.Now()
		eventCounterVect = promauto.NewCounterVec(
			prometheus.CounterOpts{Name: "events_summary_total"},
			[]string{"kind", "namespace", "reason"})
		loopRestartCounter = promauto.NewCounter(prometheus.CounterOpts{Name: "watcher_restarts_total"})
	)

	// Infinite loop to cover channel close when API fails or something.
	for {
		for eventGen := range eventWatcher.ResultChan() {
			event, ok := eventGen.Object.(*v1.Event)
			if !ok {
				log.Printf("unexpected type %s", eventGen.Object.GetObjectKind())
				continue
			}

			if event.LastTimestamp.Time.Before(startOfLoop) {
				utils.PrintfIfDebug(debug, "Skipping old event %s %s %s %s",
					event.Type, event.InvolvedObject.Kind, event.Reason, event.InvolvedObject.Namespace)
				continue
			}

			utils.PrintfIfDebug(debug, "Event %s %s %s %s",
				event.Type, event.InvolvedObject.Kind, event.Reason, event.InvolvedObject.Namespace)

			increaseStatistics(eventCounterVect, event)
		}
		// Update time when we are facing closed channel and need to go again.
		startOfLoop = time.Now()
		loopRestartCounter.Inc()
		log.Print("Loop will continue.")
	}
}

func increaseStatistics(counters *prometheus.CounterVec, event *v1.Event) {
	cnt, err := counters.GetMetricWith(
		prometheus.Labels{
			"kind":      event.InvolvedObject.Kind,
			"namespace": event.InvolvedObject.Namespace,
			"reason":    event.Reason,
		})
	if err != nil {
		log.Printf("Unable to get counter for kind: %s namespace: %s reason: %s",
			event.InvolvedObject.Kind, event.InvolvedObject.Namespace, event.Reason)
	} else {
		utils.PrintfIfDebug(debug, "Recording event %s %s %s %s",
			event.Type, event.InvolvedObject.Kind, event.Reason, event.InvolvedObject.Namespace)
		cnt.Inc()
	}
}
