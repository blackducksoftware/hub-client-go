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

import "time"

const ContentTypeBdVulnerabilityV4 = "application/vnd.blackducksoftware.vulnerability-4+json"

type bdJsonVulnerabilityV4 struct{}

func (bdJsonVulnerabilityV4) GetMimeType() string {
	return ContentTypeBdVulnerabilityV4
}

// returned by (what I saw so far --tandr)
// /api/component/xxx/version/yyy/vulnerabilities
type VulnerabilitiesList struct {
	bdJsonVulnerabilityV4
	ItemsListBase
	Items []Vulnerability `json:"items"`
}

// Data representation of the values returned by
// /api/vulnerabilities/$vulnerability
type Vulnerability struct {
	bdJsonVulnerabilityV4
	Source                 string     `json:"source"`
	Name                   string     `json:"name"`
	Title                  string     `json:"title"`
	Description            string     `json:"description"`
	TechnicalDescription   string     `json:"technicalDescription"`
	PublishedDate          *time.Time `json:"publishedDate"`
	UpdatedDate            *time.Time `json:"updatedDate"`
	DisclosureDate         *time.Time `json:"disclosureDate"`
	ExploitPublishDate     *time.Time `json:"exploitPublishDate"`
	SolutionDate           *time.Time `json:"solutionDate"`
	Solution               string     `json:"solution"`
	Severity               string     `json:"severity"` //  [LOW, MEDIUM, HIGH, CRITICAL]
	CVSS2                  *CVSS      `json:"cvss2,omitempty"`
	CVSS3                  *CVSS      `json:"cvss3,omitempty"`
	UseCVSS3               bool       `json:"useCvss3"`
	OverallScore           float32    `json:"overallScore"`
	Classifications        []string   `json:"classifications"`
	Workaround             string     `json:"workaround,omitempty"`
	VendorFixDate          *time.Time `json:"vendorFixDate,omitempty"`
	Credit                 string     `json:"credit,omitempty"`
	DiscoveryDate          *time.Time `json:"discoveryDate,omitempty"`
	VendorNotificationDate *time.Time `json:"vendorNotificationDate,omitempty"`
	ZeroDay                bool       `json:"zeroDay,omitempty"`
	UnderReview            bool       `json:"underReview,omitempty"`
	ParentAdvisory         bool       `json:"parentAdvisory,omitempty"`
	BDSATags               []string   `json:"bdsaTags,omitempty"`
	Meta                   Meta       `json:"_meta"`
}

type CVSS struct {
	BaseScore              float32               `json:"baseScore"`
	ImpactSubscore         float32               `json:"impactSubscore"`
	ExploitabilitySubscore float32               `json:"exploitabilitySubscore"`
	AccessVector           string                `json:"accessVector"` // [LOCAL, ADJACENT_NETWORK, NETWORK]
	AccessComplexity       string                `json:"accessComplexity"`
	Authentication         string                `json:"authentication"`
	ConfidentialityImpact  string                `json:"confidentialityImpact"`
	IntegrityImpact        string                `json:"integrityImpact"`
	AvailabilityImpact     string                `json:"availabilityImpact"`
	PrivilegesRequired     string                `json:"privilegesRequired"`
	Scope                  string                `json:"scope"`
	UserInteraction        string                `json:"userInteraction"`
	TemporalMetrics        *VulnerabilityMetrics `json:"temporalMetrics,omitempty"`
	Vector                 string                `json:"vector"`
}

type VulnerabilityMetrics struct {
	Exploitability   string  `json:"exploitability"`
	RemediationLevel string  `json:"remediationLevel"`
	ReportConfidence string  `json:"reportConfidence"`
	Score            float32 `json:"score"`
}

// Common Weakness Enumeration endpoint result -- retrieved from /api/cwes/{cweId}
// Links: "sources"
type CweDetails struct {
	bdJsonVulnerabilityV4
	CweId               string              `json:"id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	ExtendedDescription string              `json:"extendedDescription"`
	CommonConsequences  []CommonConsequence `json:"commonConsequences"`
	Meta                Meta                `json:"_meta"`
}

type CommonConsequence struct {
	Note             string   `json:"note"`
	Scopes           []string `json:"scopes"`
	TechnicalImpacts []string `json:"technicalImpacts"`
}
