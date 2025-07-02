package metrics

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// empty interface for the ease of setting values for supportedMetrics
	empty interface{}
)

type minimumDurationMetrics struct {
	// supportedMetrics lists all the metrics that hive currently logs, which support metricsToReport.minimumDuration customization
	// It is a map with empty interface set as value, to keep it lightweight and with O(1) look up time.
	supportedMetrics map[string]interface{}

	// metricsWithMinimumDuration will be the map of metrics that have minimumDuration customization added in metricsConfig
	metricsWithMinimumDuration map[string]metav1.Duration
}

func newMinimumDurationMetrics() *minimumDurationMetrics {
	return &minimumDurationMetrics{
		supportedMetrics: map[string]interface{}{
			//histogramVecs
			"hive_cluster_deployments_hibernation_transition_seconds":        empty,
			"hive_cluster_deployments_running_transition_seconds":            empty,
			"hive_cluster_deployments_stopping_seconds":                      empty,
			"hive_cluster_deployments_resuming_seconds":                      empty,
			"hive_cluster_deployments_waiting_for_cluster_operators_seconds": empty,
			//customMetrics
			"hive_clustersync_failing_seconds": empty,
		},
		metricsWithMinimumDuration: make(map[string]metav1.Duration),
	}
}

// isMetricSupported should be used to validate if the metric name is valid and supported minimumDuration customization
func (md *minimumDurationMetrics) isMetricSupported(name string) bool {
	if _, ok := md.supportedMetrics[name]; ok {
		return true
	}
	return false
}
