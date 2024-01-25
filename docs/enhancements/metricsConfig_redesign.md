---
title: forward_to_azure_monitor_logs
authors:
- "@vparfonov"
  reviewers:
- "@jcantrill"
- "@alanconway"
  approvers:
- "@jcantrill"
- "@alanconway"
  api-approvers:
- "@jcantrill"
- "@alanconway"
  creation-date: 2023-11-17
  last-updated: 2023-11-17
  status: implementable
  tracking-link:
- https://issues.redhat.com/browse/LOG-4606
  see-also:
-
superseded-by:
---


# HiveConfig.spec.metricsConfig redesign

[HIVE-2344](https://issues.redhat.com/browse/HIVE-2344)

- [Summary](#summary)
- [Motivation](#motivation)
  - [Current implementations](#current-implementations)
    - [Duration based threshold](#duration-based-threshold)
    - [Additional label support](#additional-label-support)
  - [Goals](#goals)
  - [Non Goals](#non-goals)
- [Proposal](#proposal)
  - [clusterDeploymentRelatedMetrics](#clusterdeploymentrelatedmetrics)
    - [minimumDuration](#minimumduration)
    - [additionalClusterDeploymentLabels](#additionalclusterdeploymentlabels)
    - [clusterDeploymentLabelSelector](#clusterdeploymentlabelselector)
  - [Implementation Details / Notes](#implementation-details--notes)
  - [Risks and Mitigations](#risks-and-mitigations)

## Summary
As a Hive user interested in customizing metrics, I would like to specify customizations for individual metrics.

## Motivation

### Current implementations
Hive publishes many metrics that can be used for observations and alerting. Presently we offer 2 major customizations for some metrics:

#### Duration based threshold
Certain duration based metrics are used for alerting. While making them as specific as possible to quickly flag an issue with a cluster, it is also necessary to be mindful of the cardinality.
It is possible to configure a threshold duration, so that metric wouldn't be logged until the value to be reported matches or exceeds the configured duration.

**Pros for this design**

Threshold can be configured as per the needs of the alerts, so they're reported when needed, and there's likelihood of these alerts needing to be silenced. Also, the customization has a 1-1 relationship with the metric itself.

**Cons we would like to overcome**

Given how we have designed it, raw/whole metric name is not used in the hiveconfig - instead, a corresponding camelCase key is hardcoded to refer to this metric. We would want to move away from such design, and implement a different - less confusing strategy that can be expandable.

#### Additional label support
Hive can manage a lot of clusters, so we usually avoid labelling any metric with cluster name or namespace to keep the cardinality in check. So, in addition to the labels hive reports for a metric, some admins might want extra labels.
For example, if admins can identify that a metric corresponds to a "managed" cluster, they can apply effective filters for their observability and can observe trends across only managed clusters, and also customize their alerts.

**Pros for this design**

The process of applying the labels to clusterdeployment, configuring those labels in hiveconfig to be reported with the metric, and using them as filters for alerts and grafana queries - is all in the hands of the admins/users and does not require any code change.
It also puts the onus of cardinality on the admins.

**Cons we would like to overcome**

The way the code is written right now - we have a hardcoded hack for cluster-type label as mandatory is no additional labels are configured, and, most importantly, any additional labels configured will be applied to all the metrics which have optional label support, even though they might not be needed for all.
We're essentially worsening the cardinality this way, and it makes adding this support to more metrics tricky. 
There's also the caveat of parsing through fixed and optional labels - ensuring there's no overlap and maintaining the order for certain metrics which do not accept map of label-value while observing, are some complications that currently exist and probably will not be fixed as a part of this enhancement.


### Goals
Allow configuring additional labels, duration based threshold and expression matching support for all metrics that are related to cluster deployment, via HiveConfig.

### Non-Goals
Other metrics - like metrics reported by HiveOperator, or metrics that do not have the related cluster deployment available while reporting (for ex, some syncset and selectorsyncset metrics) will not be affected/changed.
Adding labels only refers to the labels propagated from the Cluster Deployment, any other strings as labels cannot be configured this way.


## Proposal
Allow HiveConfig.Spec.metricsConfig to look like
```
HiveConfig:
  spec:
    metricsConfig:
      clusterDeploymentRelatedMetrics:
      - name: hive_foo_counter
        minimumDuration: 10m
        clusterDeploymentLabelSelector:
          matchLabels:
            hive.openshift.io/aro-snowflake: "true"
          matchExpressions:
            - key: hive.openshift.io/limited-support
              operator: NotIn
              values: 
                - "true"
            - key: hive.openshift.io/limited-support
              operator: NotExists
        additionalClusterDeploymentLabels:
          prom_label_name: hive.openshift.io/cd-label-key
```
### clusterDeploymentRelatedMetrics
We are only considering with metrics that are calculated whenever we parse the cluster deployment, or is tied to a cluster deployment, so we will group them all under hiveConfig.spec.metricsConfig.clusterDeploymentRelatedMetrics. Each metric that needs to be customized would need a separate entry.
Implementation would require a separate class that would encapsulate all the clusterDeploymentRelatedMetrics, and related methods for registering and observing the metrics.

#### minimumDuration
Deprecate the current implementation of hiveConfig.spec.metricsConfig.metricsWithDuration, and change it to be reported as metricsConfig.clusterDeploymentRelatedMetrics[$name].minimumDuration. The implementation of using the duration as a threshold before we report the metric stays the same.

#### additionalClusterDeploymentLabels
Deprecate current implementation of hiveConfig.spec.metricsConfig.additionalClusterDeploymentLabels and change it to be reported as metricsConfig.clusterDeploymentRelatedMetrics[$name].additionalClusterDeploymentLabels. It's implementation will not change.

#### clusterDeploymentLabelSelector
This would be a new feature, of type [LabelSelector](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#LabelSelector), and we'd use the LabelSelector.MatchLabels and/or LabelSelector.MatchExpressions to match the conditions in order to decide if a metric should be reported.
This encapsulates slightly more advanced filter logic over the existing clusterdeployment labels. In the example above, hive_foo_counter metric will only be reported for the cluster deployments that are labelled with aro-snowflake, not in limited support, and has the metric value > 10m.

### Implementation Details / Notes
- In order to implement this change, we will have to deprecate the existing hiveconfig.spec.metricsConfig.metricsWithDuration and hiveconfig.spec.metricsConfig.additionalClusterDeploymentLabels.`
- Current configuration of using a shorthand camelCase key to refer to a corresponding metric while configuring hiveconfig.spec.metricsConfig.metricsWithDuration will be removed, and users will have to use the full metric name in the proposed configuration.
- Additionally, we would also want to remove `cluster_type` as a fixed label for all the metrics that do report it now. This would be a breaking change as the consumers who want `cluster_type` will have to configure it via metricsConfig.clusterDeploymentRelatedMetrics[].additionalClusterDeploymentLabels for all the metrics they need it for.
- There are other fixed labels (like platform and version), which are directed fed from the clusterdeployment labels. Ideally we would want to feed them via hiveconfig, but we might need an exception for the cluster deployment labels that are applied by hive. 
- All the options configured for a metric within clusterDeploymentRelatedMetrics, will work in tandem with each other. For ex, if all possible options are specified for a metric, then that metric will be reported with the additional labels as per additionalClusterDeploymentLabels, and will be reported only if it matches the labels and/or expressions as per clusterDeploymentLabelSelector and if the duration to be reported exceeds the minimumDuration
- Failure modes:
  These situations will cause a panic().
  - clusterDeploymentRelatedMetrics[].name doesn't exist as a metric
  - clusterDeploymentRelatedMetrics[].name is a metric without duration but minimumDuration was specified
  - additionalClusterDeploymentLabels conflict with existing fixed label for that metric

### Risks and Mitigations
The biggest risk is the sheer number of metrics that are going to be affected.

Many metrics related to cluster deployment which do not currently have support for additionalClusterDeploymentLabels will now have that support. This might create an avenue for worsening the cardinality, however, that would be the responsibility of Hive admin.

It can get confusing for consumers to know the fixed labels each metric has.Update the [hive_metrics](https://github.com/openshift/hive/blob/master/docs/hive_metrics.md) doc to list out all labels and keep it up-to-date.

Also, the more customizations there are, the longer hiveConfig will be. Currently, we view this as an acceptable change.
