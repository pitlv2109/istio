// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sds

import (
	"fmt"

	"istio.io/istio/security/pkg/monitoring"
)

const (
	totalPush                        = "total_pushes"
	totalPushError                   = "total_push_errors"
	totalActiveConn                  = "total_active_connections"
	totalStaleConn                   = "total_stale_connections"
	stalePerConn                     = "stale_per_conn_count"
	pendingPushesPerConn             = "pending_pushes_per_connection"
	pushesPerConn                    = "pushes_per_connection"
	pushErrorsPerConn                = "push_errors_per_connection"
	rootCertExpiryTimestampPerConn   = "pushed_root_cert_expiry_timestamp_per_connection"
	serverCertExpiryTimestampPerConn = "pushed_server_cert_expiry_timestamp_per_connection"
)

var (
	// totalPushCounts records total number of SDS pushes since server starts serving.
	totalPushCounts = monitoring.NewSum(
		"total_pushes",
		"The total number of SDS pushes.",
	)

	// totalPushErrorCounts records total number of failed SDS pushes since server starts serving.
	totalPushErrorCounts = monitoring.NewSum(
		"total_push_errors",
		"The total number of failed SDS pushes.",
	)

	// totalActiveConnCounts records total number of active SDS connections.
	totalActiveConnCounts = monitoring.NewSum(
		"total_active_connections",
		"The total number of active SDS connections.",
	)

	// totalStaleConnCounts records total number of stale SDS connections.
	totalStaleConnCounts = monitoring.NewSum(
		"total_stale_connections",
		"The total number of stale SDS connections.",
	)
)

// newMonitoringMetrics creates a new monitoringMetrics.
func newMonitoringMetrics() *MetricsManager {
	fmt.Println("RETURNING METRICS")
	return NewMetricsManager(totalPushCounts, totalPushErrorCounts, totalActiveConnCounts,
		totalStaleConnCounts)
}

// generateResourcePerConnLabel returns a label that can be used to differentiate metrics by
// resource name and connection ID.
// An example label is
// bookinfo-credential-1-cacert+router~10.60.4.32~istio-ingressgateway-5b4458d75b-5kqhg.istio-system~istio-system.svc.cluster.local-7
func generateResourcePerConnLabel(resourceName, conID string) string {
	return resourceName + "+" + conID
}

func getMetricName(metricName, resourceName, conID string) string {
	return fmt.Sprintf("%s[%s]", metricName, generateResourcePerConnLabel(resourceName, conID))
}
