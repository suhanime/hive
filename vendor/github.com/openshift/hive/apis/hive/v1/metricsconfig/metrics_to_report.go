package metricsconfig

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// MetricsToReport represents metrics that have additional customizations
type MetricsToReport struct {
	// 	MetricNames is a list of metrics for which the following customizations must be added, if they support the customization
	// The name of the metric here must be valid, and it can only be present once in metricsToReport.
	MetricNames []string `json:"metricNames"`
	// MinimumDuration specifies the threshold duration for the metrics supplied in MetricNames, all of which must be duration based.
	// The corresponding metrics will be logged only if the value they report exceed the threshold duration provided.
	// For example, if a user opts-in for current clusters stopping and mentions
	// 1 hour here, only the clusters stopping for more than an hour will be reported.
	// This is a Duration value; see https://pkg.go.dev/time#ParseDuration for accepted formats.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	// +optional
	MinimumDuration *metav1.Duration `json:"minimumDuration,omitempty"`
	// AdditionalClusterDeploymentLabels allows configuration of additional labels to be applied to certain metrics.
	// The keys can be any string value suitable for a metric label (see https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels).
	// The values can be any ClusterDeployment label key (from metadata.labels). When observing an affected metric,
	// hive will label it with the specified metric key, and copy the value from the specified ClusterDeployment label.
	// For example, including {"ocp_major_version": "hive.openshift.io/version-major"} will cause affected metrics to
	// include a label key ocp_major_version with the value from the hive.openshift.io/version-major ClusterDeployment
	// label -- e.g. "4".
	// NOTE: Avoid ClusterDeployment labels whose values are unbounded, such as those representing cluster names or IDs,
	// as these will cause your prometheus database to grow indefinitely.
	// Affected metrics are those whose type implements the metricsWithDynamicLabels interface found in
	// pkg/controller/metrics/metrics_with_dynamic_labels.go
	// +optional
	AdditionalClusterDeploymentLabels *map[string]string `json:"additionalClusterDeploymentLabels,omitempty"`
}
