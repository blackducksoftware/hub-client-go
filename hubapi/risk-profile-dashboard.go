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

package hubapi

type RiskProfileDashboard struct {
	OverallRiskAggregate           OverallRiskAggregate           `json:"overallRiskAggregate"`
	ProjectRiskProfilePageView     ProjectList                    `json:"projectRiskProfilePageView"`
	PolicyViolationSeverityProfile PolicyViolationSeverityProfile `json:"policyViolationSeverityProfile"`
}

type OverallRiskAggregate struct {
	Categories Categories `json:"categories"`
}

type Categories struct {
	Version       Severity `json:"VERSION"`
	License       Severity `json:"LICENSE"`
	Vulnerability Severity `json:"VULNERABILITY"`
	Operational   Severity `json:"OPERATIONAL"`
	Activity      Severity `json:"ACTIVITY"`
}

type PolicyViolationSeverityProfile struct {
	Severities Severity `json:"severities"`
	Counts     []int    `json:"counts"`
}

type Severity struct {
	MAJOR       int `json:"MAJOR"`
	UNSPECIFIED int `json:"UNSPECIFIED"`
	TRIVIAL     int `json:"TRIVIAL"`
	MINOR       int `json:"MINOR"`
	CRITICAL    int `json:"CRITICAL"`
	OK          int `json:"OK"`
	BLOCKER     int `json:"BLOCKER"`
}
