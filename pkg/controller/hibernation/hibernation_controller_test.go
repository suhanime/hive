package hibernation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	batchv1 "k8s.io/api/batch/v1"
	certsv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakekubeclient "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configv1 "github.com/openshift/api/config/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	hiveintv1alpha1 "github.com/openshift/hive/apis/hiveinternal/v1alpha1"
	"github.com/openshift/hive/pkg/constants"
	"github.com/openshift/hive/pkg/controller/hibernation/mock"
	"github.com/openshift/hive/pkg/remoteclient"
	remoteclientmock "github.com/openshift/hive/pkg/remoteclient/mock"
	testcd "github.com/openshift/hive/pkg/test/clusterdeployment"
	testcs "github.com/openshift/hive/pkg/test/clustersync"
	testgeneric "github.com/openshift/hive/pkg/test/generic"
)

const (
	namespace = "test-namespace"
	cdName    = "test-cluster-deployment"
)

func TestReconcile(t *testing.T) {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	batchv1.AddToScheme(scheme)
	configv1.AddToScheme(scheme)
	hivev1.AddToScheme(scheme)
	hiveintv1alpha1.AddToScheme(scheme)
	machineapi.AddToScheme(scheme)

	cdBuilder := testcd.FullBuilder(namespace, cdName, scheme).Options(
		testcd.Installed(),
		testcd.WithClusterVersion("4.4.9"),
	)
	o := clusterDeploymentOptions{}
	csBuilder := testcs.FullBuilder(namespace, cdName, scheme).Options(
		testcs.WithFirstSuccessTime(time.Now().Add(-10 * time.Hour)),
	)

	tests := []struct {
		name               string
		cd                 *hivev1.ClusterDeployment
		cs                 *hiveintv1alpha1.ClusterSync
		setupActuator      func(actuator *mock.MockHibernationActuator)
		setupCSRHelper     func(helper *mock.MockcsrHelper)
		setupRemote        func(builder *remoteclientmock.MockBuilder)
		validate           func(t *testing.T, cd *hivev1.ClusterDeployment)
		expectError        bool
		expectRequeueAfter time.Duration
	}{
		{
			name: "cluster deleted",
			cd:   cdBuilder.GenericOptions(testgeneric.Deleted()).Options(o.shouldHibernate).Build(),
			cs:   csBuilder.Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				assert.Nil(t, cond, "not expecting hibernating condition")
				assert.Nil(t, runCond, "not expecting running condition")
			},
		},
		{
			name: "hibernation and running condition initialized",
			cd:   cdBuilder.Options(o.notInstalled, o.shouldHibernate).Build(),
			cs:   csBuilder.Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionUnknown, cond.Status)
				assert.Equal(t, hivev1.InitializedConditionReason, cond.Reason)
				assert.Equal(t, corev1.ConditionUnknown, runCond.Status)
				assert.Equal(t, hivev1.InitializedConditionReason, runCond.Reason)
			},
		},
		{
			name: "do not hibernate unsupported versions",
			cd:   cdBuilder.Options(testcd.WithClusterVersion("4.3.11")).Build(),
			cs:   csBuilder.Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.UnsupportedHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionTrue, runCond.Status)
				assert.Equal(t, hivev1.RunningReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.RunningReadyReason, cd.Status.PowerState)
			},
		},
		{
			name: "start hibernating, older version",
			cd:   cdBuilder.Options(o.shouldHibernate, testcd.WithClusterVersion("4.3.11")).Build(),
			cs:   csBuilder.Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.UnsupportedHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionTrue, runCond.Status)
				assert.Equal(t, hivev1.RunningReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.RunningReadyReason, cd.Status.PowerState)
			},
		},
		{
			name:        "clustersync not yet created",
			cd:          cdBuilder.Options(o.shouldHibernate).Build(),
			expectError: true,
		},
		{
			name: "start hibernating, no syncsets",
			cd:   cdBuilder.Options(o.shouldHibernate).Build(),
			// The clustersync controller creates a ClusterSync even when there are no syncsets
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StopMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.StoppingHibernationReason, cond.Reason)
				assert.Equal(t, hivev1.StoppingHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
			},
		},
		{
			name: "start hibernating, syncsets not applied",
			cd:   cdBuilder.Options(o.shouldHibernate, testcd.InstalledTimestamp(time.Now())).Build(),
			cs:   csBuilder.Options(testcs.WithNoFirstSuccessTime()).Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.SyncSetsNotAppliedReason, cond.Reason)
				assert.Equal(t, hivev1.SyncSetsNotAppliedReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionTrue, runCond.Status)
				assert.Equal(t, hivev1.RunningReadyReason, runCond.Reason)
			},
			expectError:        false,
			expectRequeueAfter: time.Duration(time.Minute * 10),
		},
		{
			name: "clear SyncSetsNotApplied",
			cd: cdBuilder.Options(
				testcd.InstalledTimestamp(time.Now()),
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:    hivev1.ClusterHibernatingCondition,
					Status:  corev1.ConditionFalse,
					Reason:  hivev1.SyncSetsNotAppliedReason,
					Message: "Cluster SyncSets have not been applied",
				}),
			).Build(),
			cs: csBuilder.Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, _ := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, hivev1.SyncSetsAppliedReason, cond.Reason)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.SyncSetsAppliedReason, cd.Status.PowerState)
			},
		},
		{
			name: "start hibernating, syncsets not applied but 10 minutes have passed since cd install",
			cd:   cdBuilder.Options(o.shouldHibernate, testcd.InstalledTimestamp(time.Now().Add(-15*time.Minute))).Build(),
			cs:   csBuilder.Options(testcs.WithNoFirstSuccessTime()).Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StopMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.StoppingHibernationReason, cond.Reason)
				assert.Equal(t, hivev1.StoppingHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
			},
		},
		{
			name: "start hibernating",
			cd:   cdBuilder.Options(o.shouldHibernate).Build(),
			cs:   csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StopMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.StoppingHibernationReason, cond.Reason)
				assert.Equal(t, hivev1.StoppingHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
			},
		},
		{
			name: "fail to stop machines",
			cd:   cdBuilder.Options(o.shouldHibernate).Build(),
			cs:   csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StopMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("error"))
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.FailedToStopHibernationReason, cond.Reason)
				assert.Equal(t, hivev1.FailedToStopHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
			},
		},
		{
			name: "stopping, machines have stopped",
			cd:   cdBuilder.Options(o.shouldHibernate, o.stopping).Build(),
			cs:   csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().MachinesStopped(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil, nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionTrue, cond.Status)
				assert.Equal(t, hivev1.HibernatingHibernationReason, cond.Reason)
				assert.Equal(t, hivev1.HibernatingHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
			},
		},
		{
			name: "stopping, machines have not stopped",
			cd:   cdBuilder.Options(o.shouldHibernate, o.stopping).Build(),
			cs:   csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StopMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				actuator.EXPECT().MachinesStopped(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					Return(false, []string{"running-1", "pending-1", "stopping-1"}, nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.StoppingHibernationReason, cond.Reason)
				assert.Equal(t, "Stopping cluster machines. Some machines have not yet stopped: pending-1,running-1,stopping-1", cond.Message)
				assert.Equal(t, hivev1.StoppingHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
			},
			expectRequeueAfter: time.Duration(time.Minute * 1),
		},
		{
			name: "stopping after MachinesFailedToStart",
			cd: cdBuilder.Options(o.shouldHibernate).Build(
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:   hivev1.ClusterHibernatingCondition,
					Status: corev1.ConditionTrue,
					Reason: hivev1.FailedToStartHibernationReason,
				},
				)),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				// Ensure we try to stop machines in this state (bugfix)
				actuator.EXPECT().StopMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.StoppingHibernationReason, cond.Reason)
				assert.Equal(t, hivev1.StoppingHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
			},
		},
		{
			name: "start resuming",
			cd:   cdBuilder.Options(o.hibernating).Build(),
			cs:   csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StartMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.WaitingForMachinesReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.WaitingForMachinesReadyReason, cd.Status.PowerState)
			},
		},
		{
			name: "resuming machines failed to start",
			cd:   cdBuilder.Options(o.hibernating).Build(),
			cs:   csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StartMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("error"))
			},
			expectError: true,
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.FailedToStartHibernationReason, cond.Reason)
				assert.Equal(t, hivev1.FailedToStartHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionTrue, runCond.Status)
				assert.Equal(t, hivev1.RunningReadyReason, runCond.Reason)
			},
		},
		{
			name: "attempt to hibernate after previous failure",
			cd: cdBuilder.Options(o.shouldHibernate).Build(
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:   hivev1.ClusterHibernatingCondition,
					Status: corev1.ConditionFalse,
					Reason: hivev1.FailedToStartHibernationReason,
				}),
				testcd.WithCondition(readyCondition(corev1.ConditionTrue, hivev1.RunningReadyReason, 6*time.Hour))),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StopMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.StoppingHibernationReason, cond.Reason)
				assert.Equal(t, hivev1.StoppingHibernationReason, cd.Status.PowerState)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
			},
		},
		{
			name: "resuming, machines have not started",
			cd: cdBuilder.Options().Build(
				testcd.WithCondition(hibernatingCondition(corev1.ConditionFalse, hivev1.ResumingOrRunningHibernationReason, 6*time.Hour)),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, "unused", 6*time.Hour))),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StartMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				actuator.EXPECT().MachinesRunning(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					Return(false, []string{"stopped-1", "pending-1"}, nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.WaitingForMachinesReadyReason, runCond.Reason)
				assert.Equal(t, "Starting cluster machines. Some machines are not yet running: pending-1,stopped-1 (step 1/4)", runCond.Message)
				assert.Equal(t, hivev1.WaitingForMachinesReadyReason, cd.Status.PowerState)
			},
			expectRequeueAfter: time.Duration(time.Minute * 1),
		},
		{
			name: "resuming unready node",
			cd: cdBuilder.Options().Build(
				testcd.WithCondition(hibernatingCondition(corev1.ConditionFalse, hivev1.ResumingOrRunningHibernationReason, 6*time.Hour)),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, "unused", 6*time.Hour))),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().MachinesRunning(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil, nil)
			},
			setupRemote: func(builder *remoteclientmock.MockBuilder) {
				fakeClient := fake.NewFakeClientWithScheme(scheme, unreadyNode()...)
				fakeKubeClient := fakekubeclient.NewSimpleClientset()
				builder.EXPECT().Build().Times(1).Return(fakeClient, nil)
				builder.EXPECT().BuildKubeClient().Times(1).Return(fakeKubeClient, nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.WaitingForNodesReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.WaitingForNodesReadyReason, cd.Status.PowerState)
			},
			expectRequeueAfter: time.Duration(time.Second * 30),
		},
		{
			name: "resuming pending csrs",
			cd: cdBuilder.Options().Build(
				testcd.WithCondition(hibernatingCondition(corev1.ConditionFalse, hivev1.ResumingOrRunningHibernationReason, 6*time.Hour)),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, "unused", 6*time.Hour))),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().MachinesRunning(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil, nil)
			},
			setupRemote: func(builder *remoteclientmock.MockBuilder) {
				fakeClient := fake.NewFakeClientWithScheme(scheme, unreadyNode()...)
				fakeKubeClient := fakekubeclient.NewSimpleClientset(csrs()...)
				builder.EXPECT().Build().Times(1).Return(fakeClient, nil)
				builder.EXPECT().BuildKubeClient().Times(1).Return(fakeKubeClient, nil)
			},
			setupCSRHelper: func(helper *mock.MockcsrHelper) {
				count := len(csrs())
				helper.EXPECT().IsApproved(gomock.Any()).Times(count).Return(false)
				helper.EXPECT().Parse(gomock.Any()).Times(count).Return(nil, nil)
				helper.EXPECT().Authorize(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(count).Return(nil)
				helper.EXPECT().Approve(gomock.Any(), gomock.Any()).Times(count).Return(nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.WaitingForNodesReadyReason, runCond.Reason)
			},
			expectRequeueAfter: time.Duration(time.Second * 30),
		},
		{
			name: "resuming nodes ready pause for clusteroperators to start and settle",
			cd: cdBuilder.Options().Build(
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:               hivev1.ClusterHibernatingCondition,
					Status:             corev1.ConditionFalse,
					Reason:             hivev1.ResumingOrRunningHibernationReason,
					LastProbeTime:      metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
					LastTransitionTime: metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
				}),
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:               hivev1.ClusterReadyCondition,
					Status:             corev1.ConditionFalse,
					Reason:             hivev1.WaitingForNodesReadyReason,
					LastProbeTime:      metav1.Time{Time: time.Now().Add(-2 * time.Minute)},
					LastTransitionTime: metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
				}),
			),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().MachinesRunning(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil, nil)
			},
			setupRemote: func(builder *remoteclientmock.MockBuilder) {
				objs := []runtime.Object{}
				objs = append(objs, readyNodes()...)
				objs = append(objs, readyClusterOperators()...)
				c := fake.NewFakeClientWithScheme(scheme, objs...)
				builder.EXPECT().Build().Times(1).Return(c, nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.PausingForClusterOperatorsToSettleReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.PausingForClusterOperatorsToSettleReadyReason, cd.Status.PowerState)
			},
			expectRequeueAfter: time.Duration(time.Minute * 2),
		},
		{
			name: "resuming continue to pause for clusteroperators to start and settle",
			cd: cdBuilder.Options().Build(
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:               hivev1.ClusterHibernatingCondition,
					Status:             corev1.ConditionFalse,
					Reason:             hivev1.ResumingOrRunningHibernationReason,
					LastProbeTime:      metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
					LastTransitionTime: metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
				}),
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:               hivev1.ClusterReadyCondition,
					Status:             corev1.ConditionFalse,
					Reason:             hivev1.PausingForClusterOperatorsToSettleReadyReason,
					LastProbeTime:      metav1.Time{Time: time.Now().Add(-2 * time.Minute)},
					LastTransitionTime: metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
				}),
			),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().MachinesRunning(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil, nil)
			},
			setupRemote: func(builder *remoteclientmock.MockBuilder) {
				objs := []runtime.Object{}
				objs = append(objs, readyNodes()...)
				objs = append(objs, degradedClusterOperators()...)
				c := fake.NewFakeClientWithScheme(scheme, objs...)
				builder.EXPECT().Build().Times(1).Return(c, nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.WaitingForClusterOperatorsReadyReason, runCond.Reason)
			},
			expectRequeueAfter: time.Duration(time.Second * 30),
		},
		{
			name: "resuming clusteroperators not ready",
			cd: cdBuilder.Options().Build(
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:               hivev1.ClusterHibernatingCondition,
					Status:             corev1.ConditionFalse,
					Reason:             hivev1.ResumingOrRunningHibernationReason,
					LastProbeTime:      metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
					LastTransitionTime: metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
				}),
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:               hivev1.ClusterReadyCondition,
					Status:             corev1.ConditionFalse,
					Reason:             hivev1.PausingForClusterOperatorsToSettleReadyReason,
					LastProbeTime:      metav1.Time{Time: time.Now().Add(-6 * time.Minute)},
					LastTransitionTime: metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
				}),
			),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().MachinesRunning(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil, nil)
			},
			setupRemote: func(builder *remoteclientmock.MockBuilder) {
				objs := []runtime.Object{}
				objs = append(objs, readyNodes()...)
				objs = append(objs, degradedClusterOperators()...)
				c := fake.NewFakeClientWithScheme(scheme, objs...)
				builder.EXPECT().Build().Times(1).Return(c, nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.WaitingForClusterOperatorsReadyReason, runCond.Reason)
			},
			expectRequeueAfter: time.Duration(time.Second * 30),
		},
		{
			name: "resuming everything ready",
			cd: cdBuilder.Options().Build(testcd.WithCondition(hivev1.ClusterDeploymentCondition{
				Type:               hivev1.ClusterHibernatingCondition,
				Status:             corev1.ConditionFalse,
				Reason:             hivev1.ResumingOrRunningHibernationReason,
				LastProbeTime:      metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
				LastTransitionTime: metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
			}),
				testcd.WithCondition(hivev1.ClusterDeploymentCondition{
					Type:               hivev1.ClusterReadyCondition,
					Status:             corev1.ConditionFalse,
					Reason:             hivev1.PausingForClusterOperatorsToSettleReadyReason,
					LastProbeTime:      metav1.Time{Time: time.Now().Add(-6 * time.Minute)},
					LastTransitionTime: metav1.Time{Time: time.Now().Add(-2 * time.Hour)},
				}),
			),
			cs: csBuilder.Build(),
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().MachinesRunning(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(true, nil, nil)
			},
			setupRemote: func(builder *remoteclientmock.MockBuilder) {
				objs := []runtime.Object{}
				objs = append(objs, readyNodes()...)
				objs = append(objs, readyClusterOperators()...)
				c := fake.NewFakeClientWithScheme(scheme, objs...)
				builder.EXPECT().Build().Times(1).Return(c, nil)
			},
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionTrue, runCond.Status)
				assert.Equal(t, hivev1.RunningReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.RunningReadyReason, cd.Status.PowerState)
			},
		},
		{
			name: "previously unsupported hibernation, now supported",
			cd:   cdBuilder.Options(o.unsupported, testcd.WithHibernateAfter(8*time.Hour)).Build(),
			cs:   csBuilder.Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				assert.Equal(t, "Hibernation capable", cond.Message)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionTrue, runCond.Status)
				assert.Equal(t, hivev1.RunningReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.RunningReadyReason, cd.Status.PowerState)
			},
		},
		{
			name: "hibernate fake cluster",
			cd: cdBuilder.Build(
				o.shouldHibernate,
				testcd.InstalledTimestamp(time.Now().Add(-1*time.Hour)),
				testcd.WithAnnotation(constants.HiveFakeClusterAnnotation, "true")),
			cs: csBuilder.Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, hivev1.HibernatingHibernationReason, cond.Reason)
				assert.Equal(t, corev1.ConditionTrue, cond.Status)
				assert.Equal(t, "Fake cluster is stopped", cond.Message)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionFalse, runCond.Status)
				assert.Equal(t, hivev1.StoppingOrHibernatingReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.HibernatingHibernationReason, cd.Status.PowerState)
			},
		},
		{
			name: "start hibernated fake cluster",
			cd: cdBuilder.Options(o.hibernating,
				testcd.WithPowerState(hivev1.RunningClusterPowerState),
				testcd.WithAnnotation(constants.HiveFakeClusterAnnotation, "true")).Build(),
			cs: csBuilder.Build(),
			validate: func(t *testing.T, cd *hivev1.ClusterDeployment) {
				cond, runCond := getHibernatingAndRunningConditions(cd)
				require.NotNil(t, cond)
				assert.Equal(t, hivev1.ResumingOrRunningHibernationReason, cond.Reason)
				assert.Equal(t, corev1.ConditionFalse, cond.Status)
				require.NotNil(t, runCond)
				assert.Equal(t, corev1.ConditionTrue, runCond.Status)
				assert.Equal(t, hivev1.RunningReadyReason, runCond.Reason)
				assert.Equal(t, hivev1.RunningReadyReason, cd.Status.PowerState)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockActuator := mock.NewMockHibernationActuator(ctrl)
			mockActuator.EXPECT().CanHandle(gomock.Any()).AnyTimes().Return(true)
			if test.setupActuator != nil {
				test.setupActuator(mockActuator)
			}
			mockBuilder := remoteclientmock.NewMockBuilder(ctrl)
			if test.setupRemote != nil {
				test.setupRemote(mockBuilder)
			}
			mockCSRHelper := mock.NewMockcsrHelper(ctrl)
			if test.setupCSRHelper != nil {
				test.setupCSRHelper(mockCSRHelper)
			}
			actuators = []HibernationActuator{mockActuator}
			var c client.Client
			if test.cs != nil {
				c = fake.NewFakeClientWithScheme(scheme, test.cd, test.cs)
			} else {
				c = fake.NewFakeClientWithScheme(scheme, test.cd)
			}

			reconciler := hibernationReconciler{
				Client: c,
				logger: log.WithField("controller", "hibernation"),
				remoteClientBuilder: func(cd *hivev1.ClusterDeployment) remoteclient.Builder {
					return mockBuilder
				},
				csrUtil: mockCSRHelper,
			}
			result, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
				NamespacedName: types.NamespacedName{Namespace: namespace, Name: cdName},
			})

			// Need to do fuzzy requeue after matching
			if test.expectRequeueAfter == 0 {
				assert.Zero(t, result.RequeueAfter)
			} else {
				assert.GreaterOrEqual(t, result.RequeueAfter.Seconds(), (test.expectRequeueAfter - 10*time.Second).Seconds(), "requeue after too small")
				assert.LessOrEqual(t, result.RequeueAfter.Seconds(), (test.expectRequeueAfter + 10*time.Second).Seconds(), "request after too large")
			}

			if test.expectError {
				assert.Error(t, err, "expected error from reconcile")
			} else {
				assert.NoError(t, err, "expected no error from reconcile")
			}
			if test.validate != nil {
				cd := &hivev1.ClusterDeployment{}
				err := c.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: cdName}, cd)
				require.Nil(t, err)
				test.validate(t, cd)
			}
			ctrl.Finish()
		})
	}

}

func TestHibernateAfter(t *testing.T) {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	hivev1.AddToScheme(scheme)
	hiveintv1alpha1.AddToScheme(scheme)

	cdBuilder := testcd.FullBuilder(namespace, cdName, scheme).Options(
		testcd.Installed(),
		testcd.WithClusterVersion("4.4.9"),
	)
	o := clusterDeploymentOptions{}
	csBuilder := testcs.FullBuilder(namespace, cdName, scheme).Options(
		testcs.WithFirstSuccessTime(time.Now().Add(-10 * time.Hour)),
	)

	tests := []struct {
		name          string
		setupActuator func(actuator *mock.MockHibernationActuator)
		cd            *hivev1.ClusterDeployment
		cs            *hiveintv1alpha1.ClusterSync

		expectError             bool
		expectRequeueAfter      time.Duration
		expectedPowerState      hivev1.ClusterPowerState
		expectedConditionReason string
	}{
		{
			name: "cluster due for hibernate no condition", // cluster that has never been hibernated and thus has no running condition
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Hour),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour))),
			cs:                 csBuilder.Build(),
			expectedPowerState: hivev1.HibernatingClusterPowerState,
		},
		{
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Hour),
				testcd.WithClusterVersion("4.3.11"),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour))),
			cs:                      csBuilder.Build(),
			expectedPowerState:      "",
			expectedConditionReason: hivev1.UnsupportedHibernationReason,
		},
		{
			name: "cluster not yet due for hibernate older version", // cluster that has never been hibernated and thus has no running condition
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Hour),
				testcd.WithClusterVersion("4.3.11"),
				testcd.InstalledTimestamp(time.Now().Add(-3*time.Hour))),
			cs: csBuilder.Options(
				testcs.WithFirstSuccessTime(time.Now().Add(-3 * time.Hour)),
			).Build(),
			expectedPowerState:      "",
			expectedConditionReason: hivev1.UnsupportedHibernationReason,
		},
		{
			name: "cluster not yet due for hibernate no running condition", // cluster that has never been hibernated
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(12*time.Hour),
				testcd.WithPowerState(hivev1.RunningClusterPowerState),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour))),
			cs:                 csBuilder.Build(),
			expectRequeueAfter: 2 * time.Hour,
			expectedPowerState: hivev1.RunningClusterPowerState,
		},
		{
			name: "cluster with running condition due for hibernate",
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Hour),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, hivev1.RunningReadyReason, 9*time.Hour)),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour))),
			cs:                 csBuilder.Build(),
			expectedPowerState: hivev1.HibernatingClusterPowerState,
		},
		{
			name: "cluster with running condition not due for hibernate",
			cd: cdBuilder.Build(
				testcd.WithCondition(hibernatingCondition(corev1.ConditionFalse, hivev1.ResumingOrRunningHibernationReason, 6*time.Hour)),
				testcd.WithCondition(readyCondition(corev1.ConditionTrue, hivev1.RunningReadyReason, 6*time.Hour)),
				testcd.WithHibernateAfter(20*time.Hour),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour))),
			cs:                 csBuilder.Build(),
			expectRequeueAfter: 14 * time.Hour,
			expectedPowerState: "",
		},
		{
			name: "cluster waking from hibernate",
			setupActuator: func(actuator *mock.MockHibernationActuator) {
				actuator.EXPECT().StartMachines(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
				actuator.EXPECT().MachinesRunning(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).
					Return(false, []string{"no-running-1", "no-running-2"}, nil)
			},
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Hour),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour)),
				testcd.WithCondition(hibernatingCondition(corev1.ConditionFalse, hivev1.ResumingOrRunningHibernationReason, 8*time.Hour)),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, hivev1.WaitingForNodesReadyReason, 8*time.Hour)),
				o.shouldRun),
			cs:                 csBuilder.Build(),
			expectedPowerState: hivev1.RunningClusterPowerState,
			expectRequeueAfter: stateCheckInterval,
		},
		{
			name: "cluster due for hibernate, no syncsets",
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Minute),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, hivev1.RunningReadyReason, 8*time.Minute)),
				testcd.InstalledTimestamp(time.Now().Add(-8*time.Minute))),
			// The clustersync controller creates an empty ClusterSync even when there are no syncsets.
			cs:                 csBuilder.Build(),
			expectedPowerState: hivev1.HibernatingClusterPowerState,
		},
		{
			name: "cluster due for hibernate but syncsets not applied",
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Minute),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, hivev1.RunningReadyReason, 8*time.Minute)),
				testcd.InstalledTimestamp(time.Now().Add(-8*time.Minute))),
			cs: csBuilder.Options(
				testcs.WithNoFirstSuccessTime(),
			).Build(),
			expectError:        false,
			expectedPowerState: "",
			expectRequeueAfter: time.Duration(time.Minute * 2),
		},
		{
			name: "cluster due for hibernate, syncsets not applied but 10 minutes have passed since cd install",
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Hour),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, hivev1.RunningReadyReason, 9*time.Hour)),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour))),
			cs: csBuilder.Options(
				testcs.WithNoFirstSuccessTime(),
			).Build(),
			expectedPowerState: hivev1.HibernatingClusterPowerState,
		},
		{
			name: "cluster due for hibernate, syncsets successfully applied",
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Hour),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, hivev1.RunningReadyReason, 9*time.Hour)),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour))),
			cs: csBuilder.Options(
				testcs.WithFirstSuccessTime(time.Now()),
			).Build(),
			expectedPowerState: hivev1.HibernatingClusterPowerState,
		},
		{
			name: "fake cluster due for hibernate but syncsets not applied",
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Minute),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, hivev1.RunningReadyReason, 8*time.Minute)),
				testcd.WithAnnotation(constants.HiveFakeClusterAnnotation, "true"),
				testcd.InstalledTimestamp(time.Now().Add(-8*time.Minute))),
			cs: csBuilder.Options(
				testcs.WithNoFirstSuccessTime(),
			).Build(),
			expectError:        false,
			expectedPowerState: "",
			expectRequeueAfter: time.Duration(time.Minute * 2),
		},
		{
			name: "fake cluster due for hibernate, syncsets successfully applied",
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(8*time.Hour),
				testcd.WithCondition(readyCondition(corev1.ConditionFalse, hivev1.RunningReadyReason, 9*time.Hour)),
				testcd.WithAnnotation(constants.HiveFakeClusterAnnotation, "true"),
				testcd.InstalledTimestamp(time.Now().Add(-10*time.Hour))),
			cs: csBuilder.Options(
				testcs.WithFirstSuccessTime(time.Now()),
			).Build(),
			expectedPowerState: hivev1.HibernatingClusterPowerState,
		},
		{
			name: "hibernate fake cluster",
			cd: cdBuilder.Build(
				testcd.WithHibernateAfter(1*time.Hour),
				testcd.InstalledTimestamp(time.Now().Add(-1*time.Hour)),
				testcd.WithAnnotation(constants.HiveFakeClusterAnnotation, "true")),
			cs:                 csBuilder.Build(),
			expectedPowerState: hivev1.HibernatingClusterPowerState,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockActuator := mock.NewMockHibernationActuator(ctrl)
			mockActuator.EXPECT().CanHandle(gomock.Any()).AnyTimes().Return(true)
			if test.setupActuator != nil {
				test.setupActuator(mockActuator)
			}
			mockBuilder := remoteclientmock.NewMockBuilder(ctrl)
			mockCSRHelper := mock.NewMockcsrHelper(ctrl)
			actuators = []HibernationActuator{mockActuator}
			var c client.Client
			if test.cs != nil {
				c = fake.NewFakeClientWithScheme(scheme, test.cd, test.cs)
			} else {
				c = fake.NewFakeClientWithScheme(scheme, test.cd)
			}

			reconciler := hibernationReconciler{
				Client: c,
				logger: log.WithField("controller", "hibernation"),
				remoteClientBuilder: func(cd *hivev1.ClusterDeployment) remoteclient.Builder {
					return mockBuilder
				},
				csrUtil: mockCSRHelper,
			}
			result, err := reconciler.Reconcile(context.TODO(), reconcile.Request{
				NamespacedName: types.NamespacedName{Namespace: namespace, Name: cdName},
			})

			if test.expectError {
				assert.Error(t, err, "expected error from reconcile")
			} else {
				assert.NoError(t, err, "expected no error from reconcile")
			}

			// Need to do fuzzy requeue after matching
			if test.expectRequeueAfter == 0 {
				assert.Zero(t, result.RequeueAfter)
			} else {
				assert.GreaterOrEqual(t, result.RequeueAfter.Seconds(), (test.expectRequeueAfter - 10*time.Second).Seconds(), "requeue after too small")
				assert.LessOrEqual(t, result.RequeueAfter.Seconds(), (test.expectRequeueAfter + 10*time.Second).Seconds(), "request after too large")
			}

			cd := &hivev1.ClusterDeployment{}
			err = c.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: cdName}, cd)
			require.NoError(t, err, "error looking up ClusterDeployment")
			assert.Equal(t, test.expectedPowerState, cd.Spec.PowerState, "unexpected PowerState")

		})
	}
}

func hibernatingCondition(status corev1.ConditionStatus, reason string, lastTransitionAgo time.Duration) hivev1.ClusterDeploymentCondition {
	return hivev1.ClusterDeploymentCondition{
		Type:               hivev1.ClusterHibernatingCondition,
		Status:             status,
		Message:            "unused",
		Reason:             reason,
		LastTransitionTime: metav1.NewTime(time.Now().Add(-lastTransitionAgo)),
	}
}

func readyCondition(status corev1.ConditionStatus, reason string, lastTransitionAgo time.Duration) hivev1.ClusterDeploymentCondition {
	return hivev1.ClusterDeploymentCondition{
		Type:               hivev1.ClusterReadyCondition,
		Status:             status,
		Message:            "unused",
		Reason:             reason,
		LastTransitionTime: metav1.NewTime(time.Now().Add(-lastTransitionAgo)),
	}
}

type clusterDeploymentOptions struct{}

func (*clusterDeploymentOptions) notInstalled(cd *hivev1.ClusterDeployment) {
	cd.Spec.Installed = false
}
func (*clusterDeploymentOptions) shouldHibernate(cd *hivev1.ClusterDeployment) {
	cd.Spec.PowerState = hivev1.HibernatingClusterPowerState
}
func (*clusterDeploymentOptions) shouldRun(cd *hivev1.ClusterDeployment) {
	cd.Spec.PowerState = hivev1.RunningClusterPowerState
}
func (*clusterDeploymentOptions) stopping(cd *hivev1.ClusterDeployment) {
	cd.Status.Conditions = append(cd.Status.Conditions, hivev1.ClusterDeploymentCondition{
		Type:   hivev1.ClusterHibernatingCondition,
		Reason: hivev1.StoppingHibernationReason,
		Status: corev1.ConditionFalse,
	})
}
func (*clusterDeploymentOptions) hibernating(cd *hivev1.ClusterDeployment) {
	cd.Status.Conditions = append(cd.Status.Conditions, hivev1.ClusterDeploymentCondition{
		Type:   hivev1.ClusterHibernatingCondition,
		Reason: hivev1.HibernatingHibernationReason,
		Status: corev1.ConditionTrue,
	})
}
func (*clusterDeploymentOptions) unsupported(cd *hivev1.ClusterDeployment) {
	cd.Status.Conditions = append(cd.Status.Conditions, hivev1.ClusterDeploymentCondition{
		Type:   hivev1.ClusterHibernatingCondition,
		Status: corev1.ConditionFalse,
		Reason: hivev1.UnsupportedHibernationReason,
	})
}

func getHibernatingAndRunningConditions(cd *hivev1.ClusterDeployment) (*hivev1.ClusterDeploymentCondition, *hivev1.ClusterDeploymentCondition) {
	var hibCond *hivev1.ClusterDeploymentCondition
	var runCond *hivev1.ClusterDeploymentCondition
	for i := range cd.Status.Conditions {
		if cd.Status.Conditions[i].Type == hivev1.ClusterHibernatingCondition {
			hibCond = &cd.Status.Conditions[i]
		} else if cd.Status.Conditions[i].Type == hivev1.ClusterReadyCondition {
			runCond = &cd.Status.Conditions[i]
		}
	}
	return hibCond, runCond
}

func readyNodes() []runtime.Object {
	nodes := make([]runtime.Object, 5)
	for i := 0; i < len(nodes); i++ {
		node := &corev1.Node{}
		node.Name = fmt.Sprintf("node-%d", i)
		node.Status.Conditions = []corev1.NodeCondition{
			{
				Type:   corev1.NodeReady,
				Status: corev1.ConditionTrue,
			},
		}
		nodes[i] = node
	}
	return nodes
}

func unreadyNode() []runtime.Object {
	node := &corev1.Node{}
	node.Name = "unready"
	node.Status.Conditions = []corev1.NodeCondition{
		{
			Type:   corev1.NodeReady,
			Status: corev1.ConditionFalse,
		},
	}
	return append(readyNodes(), node)
}

func readyClusterOperators() []runtime.Object {
	cos := make([]runtime.Object, 5)
	for i := 0; i < len(cos); i++ {
		co := &configv1.ClusterOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("clusteroperator%d", i),
			},
			Status: configv1.ClusterOperatorStatus{
				Conditions: []configv1.ClusterOperatorStatusCondition{
					{
						Type:   "Available",
						Status: configv1.ConditionTrue,
					},
					{
						Type:   "Progressing",
						Status: configv1.ConditionFalse,
					},
					{
						Type:   "Degraded",
						Status: configv1.ConditionFalse,
					},
				},
			},
		}
		cos[i] = co
	}
	return cos
}

func degradedClusterOperators() []runtime.Object {
	cos := make([]runtime.Object, 5)
	for i := 0; i < len(cos); i++ {
		co := &configv1.ClusterOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("clusteroperator%d", i),
			},
			Status: configv1.ClusterOperatorStatus{
				Conditions: []configv1.ClusterOperatorStatusCondition{
					{
						Type:   "Available",
						Status: configv1.ConditionTrue,
					},
					{
						Type:   "Progressing",
						Status: configv1.ConditionFalse,
					},
					{
						Type:   "Degraded",
						Status: configv1.ConditionTrue,
					},
				},
			},
		}
		cos[i] = co
	}
	return cos
}

func csrs() []runtime.Object {
	result := make([]runtime.Object, 5)
	for i := 0; i < len(result); i++ {
		csr := &certsv1.CertificateSigningRequest{}
		csr.Name = fmt.Sprintf("csr-%d", i)
		result[i] = csr
	}
	return result
}

func Test_timeBeforeClusterSyncCheck(t *testing.T) {
	buildCD := func(installed time.Time) *hivev1.ClusterDeployment {
		return testcd.BasicBuilder().Build(testcd.InstalledTimestamp(installed))
	}
	tests := []struct {
		name string
		cd   *hivev1.ClusterDeployment
		want time.Duration
	}{
		{
			// We should never hit this code path in real life, but make sure it behaves sanely anyway
			name: "Not yet installed",
			cd:   testcd.BasicBuilder().Build(),
			want: 2 * time.Minute,
		},
		{
			// This should also never happen, but is substantially the same as "waaay in the past"
			name: "Installed at epoch",
			cd:   buildCD(time.Time{}),
			want: 0,
		},
		{
			name: "Older than max wait time",
			cd:   buildCD(time.Now().Add(-3 * hibernateAfterSyncSetsNotApplied)),
			want: 0,
		},
		{
			name: "Really new",
			cd:   buildCD(time.Now().Add(-3 * time.Second)),
			want: 10 * time.Second,
		},
		{
			name: "Pretty new",
			cd:   buildCD(time.Now().Add(-2 * time.Minute)),
			want: time.Minute,
		},
		{
			name: "Middlin",
			cd:   buildCD(time.Now().Add(-5 * time.Minute)),
			want: 3 * time.Minute,
		},
		{
			name: "Almost expired",
			cd:   buildCD(time.Now().Add((-9 * time.Minute) - (30 * time.Second))),
			want: 30 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeBeforeClusterSyncCheck(tt.cd); got.Round(time.Second) != tt.want {
				t.Errorf("timeBeforeClusterSyncCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
