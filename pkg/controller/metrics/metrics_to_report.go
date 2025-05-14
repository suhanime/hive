package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/openshift/hive/apis/hive/v1"
	_ "github.com/openshift/hive/apis/hive/v1/metricsconfig"
)

var (
	// Maintain a list of all the metrics with their names for validation and observing purposes
	// We have to maintain the name of the metrics separately as a map because of the way metrics are currently declared, their names are not available in the code
	counterVecs = map[*prometheus.CounterOpts]string{
		MetricClustersCreated.CounterOpts:         "hive_cluster_deployments_created_total",
		MetricClustersInstalled.CounterOpts:       "hive_cluster_deployments_installed_total",
		MetricClustersDeleted.CounterOpts:         "hive_cluster_deployments_deleted_total",
		MetricProvisionFailedTerminal.CounterOpts: "hive_cluster_deployments_provision_failed_terminal_total",
		MetricClusterProvisionsTotal.CounterOpts:  "hive_cluster_provision_results_total",
		MetricInstallErrors.CounterOpts:           "hive_install_errors",
	}
	histogramVecs = map[*prometheus.HistogramVec]string{
		MetricClusterHibernationTransitionSeconds: "hive_cluster_deployments_hibernation_transition_seconds",
		MetricClusterReadyTransitionSeconds:       "hive_cluster_deployments_running_transition_seconds",
		MetricStoppingClustersSeconds:             "hive_cluster_deployments_stopping_seconds",
		MetricResumingClustersSeconds:             "hive_cluster_deployments_resuming_seconds",
		MetricWaitingForCOClustersSeconds:         "hive_cluster_deployments_waiting_for_cluster_operators_seconds",
	}
	histogramOpts = map[*prometheus.HistogramOpts]string{
		MetricCompletedInstallJobRestarts.HistogramOpts: "hive_cluster_deployment_completed_install_restart",
		MetricInstallFailureSeconds.HistogramOpts:       "hive_cluster_deployment_install_failure_total_seconds",
		MetricInstallSuccessSeconds.HistogramOpts:       "hive_cluster_deployment_install_success_total_seconds",
	}
	customMetrics = map[*prometheus.Desc]string{
		metricClusterSyncFailingSeconds: "hive_clustersync_failing_seconds",
	}
)
