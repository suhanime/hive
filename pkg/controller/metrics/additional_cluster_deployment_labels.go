package metrics

type additionalClusterDeploymentLabels struct {
	// supportedMetrics lists all the metrics that hive currently logs, which support metricsToReport.additionalClusterDeploymentLabels customization
	// It is a map with empty interface set as value, to keep it lightweight and with O(1) look up time.
	supportedMetrics map[string]interface{}

	// metricsWithAdditionalCDLabels will be the map of metrics that have additionalClusterDeploymentLabels customization added in metricsConfig,
	// and the corresponding labels provided
	metricsWithAdditionalCDLabels map[string][]string
}

func newAdditionalClusterDeploymentLabels() *additionalClusterDeploymentLabels {
	return &additionalClusterDeploymentLabels{
		supportedMetrics: map[string]interface{}{
			// counterVecs
			"hive_cluster_deployments_created_total":                   empty,
			"hive_cluster_deployments_installed_total":                 empty,
			"hive_cluster_deployments_deleted_total":                   empty,
			"hive_cluster_deployments_provision_failed_terminal_total": empty,
			"hive_cluster_provision_results_total":                     empty,
			"hive_install_errors":                                      empty,
			//	histogramOpts
			"hive_cluster_deployment_completed_install_restart":     empty,
			"hive_cluster_deployment_install_failure_total_seconds": empty,
			"hive_cluster_deployment_install_success_total_seconds": empty,
			//	customMetrics
			"hive_clustersync_failing_seconds": empty,
		},
		metricsWithAdditionalCDLabels: make(map[string][]string),
	}
}

// isMetricSupported should be used to validate if the metric name is valid and supported additionalClusterDeploymentLabels customization
func (al *additionalClusterDeploymentLabels) isMetricSupported(name string) bool {
	if _, ok := al.supportedMetrics[name]; ok {
		return true
	}
	return false
}
