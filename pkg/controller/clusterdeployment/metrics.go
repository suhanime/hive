package clusterdeployment

import (
	"github.com/prometheus/client_golang/prometheus"

	"sigs.k8s.io/controller-runtime/pkg/metrics"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	hivemetrics "github.com/openshift/hive/pkg/controller/metrics"
	controllerutils "github.com/openshift/hive/pkg/controller/utils"
)

var (
	metricInstallJobDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "hive_cluster_deployment_install_job_duration_seconds",
			Help:    "Distribution of the runtime of completed install jobs.",
			Buckets: []float64{1800, 2400, 3000, 3600, 4500, 5400, 7200},
		},
	)
	metricInstallDelaySeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "hive_cluster_deployment_install_job_delay_seconds",
			Help:    "Time between cluster deployment creation and creation of the job to install/provision the cluster.",
			Buckets: []float64{60, 120, 180, 240, 300, 600, 1200, 1800, 2700, 3600},
		},
	)
	metricImageSetDelaySeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "hive_cluster_deployment_imageset_job_delay_seconds",
			Help:    "Time between cluster deployment creation and creation of the job which resolves the installer image to use for a ClusterImageSet.",
			Buckets: []float64{10, 30, 60, 300, 600, 1200, 1800},
		},
	)
	metricDNSDelaySeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "hive_cluster_deployment_dns_delay_seconds",
			Help:    "Time between cluster deployment with spec.manageDNS creation and the DNSZone becoming ready.",
			Buckets: []float64{10, 30, 60, 300, 600, 1200, 1800},
		},
	)
)

func incProvisionFailedTerminal(cd *hivev1.ClusterDeployment) {
	poolNSName := ""
	if poolRef := cd.Spec.ClusterPoolRef; poolRef != nil {
		poolNSName = poolRef.Namespace + "/" + poolRef.PoolName
	}
	stoppedReason := "unknown"
	stoppedCondition := controllerutils.FindCondition(cd.Status.Conditions, hivev1.ProvisionStoppedCondition)
	if stoppedCondition != nil {
		stoppedReason = stoppedCondition.Reason
	}
	fixedLabels := map[string]string{
		"clusterpool_namespacedname": poolNSName,
		"failure_reason":             stoppedReason,
	}
	hivemetrics.MetricProvisionFailedTerminal.Observe(cd, fixedLabels, 1)
}

func init() {
	metrics.Registry.MustRegister(metricInstallJobDuration)
	metrics.Registry.MustRegister(metricInstallDelaySeconds)
	metrics.Registry.MustRegister(metricImageSetDelaySeconds)
	metrics.Registry.MustRegister(metricDNSDelaySeconds)
}
