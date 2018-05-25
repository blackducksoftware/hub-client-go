// Copyright 2018 Synopsys, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hubclient

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var httpResponseTimes *prometheus.HistogramVec

func recordResponseTime(path string, method string, statusCode int, duration time.Duration) {
	milliseconds := float64(duration / time.Millisecond)
	httpResponseTimes.With(prometheus.Labels{"path": path, method: method, "statusCode": fmt.Sprintf("%d", statusCode)}).Observe(milliseconds)
}

func init() {
	httpResponseTimes = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "hub",
		Subsystem: "client_go",
		Name:      "http_handled_status_codes",
		Help:      "response times for http requests in milliseconds",
		Buckets:   prometheus.ExponentialBuckets(1, 2, 20),
	}, []string{"path", "method", "code"})
	prometheus.MustRegister(httpResponseTimes)
}
