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
	ProjectRiskProfilePageView     RiskProfileProjectList         `json:"projectRiskProfilePageView"`
	PolicyViolationSeverityProfile PolicyViolationSeverityProfile `json:"policyViolationSeverityProfile"`
}

type OverallRiskAggregate struct {
	Categories Categories `json:"categories"`
}

type Categories struct {
	Version       CategotySeverity `json:"VERSION"`
	License       CategotySeverity `json:"LICENSE"`
	Vulnerability CategotySeverity `json:"VULNERABILITY"`
	Operational   CategotySeverity `json:"OPERATIONAL"`
	Activity      CategotySeverity `json:"ACTIVITY"`
}

type PolicyViolationSeverityProfile struct {
	Severities PolicySeverity `json:"severities"`
	Counts     []int          `json:"counts"`
}

type PolicySeverity struct {
	MAJOR       int `json:"MAJOR"`
	UNSPECIFIED int `json:"UNSPECIFIED"`
	TRIVIAL     int `json:"TRIVIAL"`
	MINOR       int `json:"MINOR"`
	CRITICAL    int `json:"CRITICAL"`
	OK          int `json:"OK"`
	BLOCKER     int `json:"BLOCKER"`
}

type CategotySeverity struct {
	CRITICAL int `json:"CRITICAL"`
	HIGH     int `json:"HIGH"`
	MEDIUM   int `json:"MEDIUM"`
	LOW      int `json:"LOW"`
	UNKNOWN  int `json:"UNKNOWN"`
	OK       int `json:"OK"`
}

type RiskProfileProjectList struct {
	TotalCount int                  `json:"totalCount"`
	Items      []RiskProfileProject `json:"items"`
}

type RiskProfileProject struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	RiskProfile RiskProfile `json:"riskProfile"`
}

type RiskProfile struct {
	Categories Categories `json:"categories"`
}
