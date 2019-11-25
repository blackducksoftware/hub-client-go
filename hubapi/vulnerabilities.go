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

// Data representation of the values returned by /api/vulnerabilities/$vulerability
type Vulnerability struct {
	Source               string `json:"source"`
	Name                 string `json:"name"`
	Title                string `json:"title"`
	Description          string `json:"description"`
	TechnicalDescription string `json:"technicalDescription"`
	PublishedDate        string `json:"publishedDate"`
	UpdatedDate          string `json:"updatedDate"`
	DisclosureDate       string `json:"disclosureDate"`
	Solution             string `json:"solution"`
	Severity             string `json:"severity"`
	CVSS2                CVSS   `json:"cvss2"`
	CVSS3                CVSS   `json:"cvss3"`
	UseCVSS3             bool   `json:"useCvss3"`
	ZeroDay              bool   `json:"zeroDay"`
	UnderReview          bool   `json:"underReview"`
	ParentAdvisory       bool   `json:"parentAdvisory"`
	Credit               string `json:"credit"`
	Workaround           string `json:"workaround"`
	VendorFixDate        string `json:"vendorFixDate"`
	Meta                 Meta   `json:"_meta"`
}

type CVSS struct {
	BaseScore              float32              `json:"baseScore"`
	ImpactSubscore         float32              `json:"impactSubscore"`
	ExploitabilitySubscore float32              `json:"exploitabilitySubscore"`
	AccessVector           string               `json:"accessVector"`
	AccessComplexity       string               `json:"accessComplexity"`
	Authentication         string               `json:"authentication"`
	ConfidentialityImpact  string               `json:"confidentialityImpact"`
	IntegrityImpact        string               `json:"integrityImpact"`
	AvailabilityImpact     string               `json:"availabilityImpact"`
	PrivilegesRequired     string               `json:"privilegesRequired"`
	Scope                  string               `json:"scope"`
	UserInteraction        string               `json:"userInteraction"`
	TemporalMetrics        VulnerabilityMetrics `json:"temporalMetrics"`
	Vector                 string               `json:"vector"`
}

type VulnerabilityMetrics struct {
	Exploitability   string  `json:"exploitability"`
	RemediationLevel string  `json:"remediationLevel"`
	ReportConfidence string  `json:"reportConfidence"`
	Score            float32 `json:"score"`
}
