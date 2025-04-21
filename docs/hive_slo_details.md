<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Cluster Creation Latency](#cluster-creation-latency)
    - [SLI description](#sli-description)
    - [SLI Rationale](#sli-rationale)
    - [Implementation details](#implementation-details)
    - [SLO Rationale](#slo-rationale)
    - [Alerts](#alerts)
- [Cluster Destruction Latency](#cluster-destruction-latency)
    - [SLI description and Implementation details](#sli-description-and-implementation-details)
    - [SLO Rationale](#slo-rationale-1)
    - [Alerts](#alerts-1)
- [API Availability](#api-availability)
    - [SLI description](#sli-description-1)
    - [SLI Rationale and Implementation details](#sli-rationale-and-implementation-details)
    - [Alerts](#alerts-2)
- [API Correctness](#api-correctness)
    - [SLI description](#sli-description-2)
    - [SLI Rationale and Implementation details](#sli-rationale-and-implementation-details-1)
    - [Alerts](#alerts-3)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

NOTE: The SLIs, SLOs and alerts covered in this document relate only to the Hive instances deployed by OSD.

# Cluster Creation Latency

### SLI description
Time taken from cluster deployment creation to the creation of install job should take less than 5 minutes.

### SLI Rationale
Cluster creation consists of 2 components:
- Time taken for Hive to kickstart an install job, and
- Time taken for the installer pod to complete cluster provisioning.

While we do report metrics around cluster provisioning, it is not a function of Hive to guarantee its success.
Hence, only the time taken for install job to start is considered for this SLI.

### Implementation details
When a cluster deployment is created (for a cluster that needs to be provisioned), Hive creates an install pod that runs the installer to provision a cluster.
The duration between ClusterDeployment's CreationTimeStamp to that of the creation of the install pod is observed via the histogram `hive_cluster_deployment_install_job_delay_seconds`.
Hive also reconciles existing ClusterInstalls and observes this metric.


The time reported here encompasses various processes Hive goes through to ensure a provision can succeed, from deciding the OpenShift release image to use, validating if
configured pre-conditions have met, to creating DNS zone for the cluster if necessary.

### SLO Rationale
Observing the trends of this histogram over 30 days at a time have helped us set the goal of ensuring the cluster creation latency takes less than 5 minutes 98% of the time.
It must be acknowledged that sometimes the SLO isn't met when the DNS related delay sees an increase.

### Alerts
Presently, a `High` level alert for `InstallJobDelayHigh` is configured to fire if the average over the past 6 hours is more than 300 seconds consistently for 10 minutes.

# Cluster Destruction Latency

### SLI description and Implementation details
When a ClusterDeployment is deleted, a [deprovision](https://github.com/openshift/hive/blob/master/docs/using-hive.md#cluster-deprovisioning) job will spawn which repeatedly tries to teardown all known cloud resources matching the cluster's infra ID tag, until nothing is left.
Hive observes the time taken for the same via the histogram `hive_cluster_deployment_deprovision_underway_seconds`.
This metric is reported via a custom collector and includes the cluster deployment name as a label. For now, we report it always, but we have plans to make this metric optional.

### SLO Rationale 
Presently there is no SLO around cluster deprovisioning. In our experience, it usually fails when the credentials are no longer valid, or in situations when manual intervention might be required.

### Alerts
A `medium` level alert for `ClusterDeprovisioningDelay` is configured to fire if a cluster deprovision takes longer than 6 hours.

# API Availability

### SLI description
Requests to Hive cluster's Kube API should be less than 15 per second.

### SLI Rationale and Implementation details
Hive reports a counter `hive_kube_client_requests_total` that is incremented for each kube client request. If there are too many requests, 
it can indicate a some sort of hotloop that can bring down Hive's API.

The metric also reports the method (POST and non POST) of the kube client requests, and this can be helpful in differentiating the urgency of an associated alert. 
For instance, the high rate of increase in POST requests can point to a hotloop of  etcd creations and would require an urgent scale down of Hive, whereas a similar increase in non POST requests isn't as urgent.

### Alerts
There's a `critical` level alert configured for `LocalKubeClientRequestsHigh - POST requests` if more than 15 requests per second are observed.
Similar alert for `LocalKubeClientRequestsHigh - non-POST requests` is set to `medium` level.

Other related `critical` alerts to API availability (like `HiveControllersDown`, `HiveClusterSyncDown`, `HiveOperatorDown` and `HiveDeploymentFailed`) are configured around whether Hive pods are up and running, and if Hive is deployed properly.

# API Correctness

### SLI description
Hive controllers shouldn't report a sudden increase of reconciliation errors.

### SLI Rationale and Implementation details
Controller Runtime reports the counter `controller_runtime_reconcile_errors_total` for the total number of reconciliation errors per controller.

### Alerts
A `critical` level alert for `ControllerErrorsHigh` is configured to fire if hive controllers report the rate of increase in reconciliation errors more than 1 over 15 minutes.