// Copyright 2021 Synopsys, Inc.
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

const (
	ContentTypeRapidScanRequest = "application/vnd.blackducksoftware.developer-scan-1-ld-2+json"
	ContentTypeRapidScanResults = "application/vnd.blackducksoftware.scan-5+json"
)

type RapidScanResult struct {
	Count          int                  `json:"totalCount"`
	Components     []RapidScanComponent `json:"items"`
	AppliedFilters []interface{}        `json:"appliedFilters"`
	Meta           Meta                 `json:"_meta"`
}
type Policy struct {
	Name        string `json:"policyName"`
	Description string `json:"description"`
	Severity    string `json:"policySeverity"`
}
type ComponentVulnerability struct {
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	Severity          string     `json:"vulnSeverity"`
	OverallScore      float32    `json:"overallScore"`
	ViolatingPolicies []Policy   `json:"violatingPolicies"`
	PublishedDate     *time.Time `json:"publishedDate"`
	VendorFixDate     *time.Time `json:"vendorFixDate,omitempty"`
	Solution          string     `json:"solution"`
	Workaround        string     `json:"workaround,omitempty"`
	Meta              Meta        `json:"_meta"`
}
type ComponentLicense struct {
	Name       string `json:"name"`
	FamilyName string `json:"licenseFamilyName"`
	Meta       Meta   `json:"_meta"`
}
type RapidScanComponent struct {
	Name                           string                   `json:"componentName"`
	Version                        string                   `json:"versionName"`
	Identifier                     string                   `json:"componentIdentifier"`
	ExternalId                     string                   `json:"externalId"`
	OriginId                       string                   `json:"originId"`
	ViolatingPolicies              []Policy                 `json:"violatingPolicies"`
	ComponentViolatingPolicies     []Policy                 `json:"componentViolatingPolicies"`
	Vulnerabilities                []ComponentVulnerability `json:"allVulnerabilities"`
	Licenses                       []ComponentLicense       `json:"allLicenses"`
	PolicyViolationVulnerabilities []ComponentVulnerability `json:"policyViolationVulnerabilities"`
	PolicyViolationLicenses        []ComponentLicense       `json:"policyViolationLicenses"`
	PartiallyEvaluatedPolicies     []string                 `json:"partiallyEvaluatedPolicies"`
	NonEvaluatedPolicies           []string                 `json:"nonEvaluatedPolicies"`
	DependencyTree                 [][]string               'json:"dependencyTree"'
	ShortTermUpgradeGuidance       string                   'json:"shortTermUpgradeGuidance"'
	LongTermUpgradeGuidance        string                   'json:"longTermUpgradeGuidance"'
	Meta                           Meta                     `json:"_meta"`
}
