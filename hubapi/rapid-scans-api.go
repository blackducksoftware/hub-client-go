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

import "time"

const (
	ContentTypeRapidScanRequest = "application/vnd.blackducksoftware.developer-scan-1-ld-2+json"
	ContentTypeRapidScanResults = "application/vnd.blackducksoftware.scan-6+json"
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
	Source            string     `json:"source"`
	Description       string     `json:"description"`
	Severity          string     `json:"vulnSeverity"`
	OverallScore      float32    `json:"overallScore"`
	ViolatingPolicies []Policy   `json:"violatingPolicies"`
	PublishedDate     *time.Time `json:"publishedDate"`
	VendorFixDate     *time.Time `json:"vendorFixDate,omitempty"`
	Solution          string     `json:"solution"`
	Workaround        string     `json:"workaround,omitempty"`
	CWEIds            []string   `json:"cweIds"`
	Meta              Meta       `json:"_meta"`
}
type ComponentLicense struct {
	Name       string `json:"name"`
	FamilyName string `json:"licenseFamilyName"`
	Meta       Meta   `json:"_meta"`
}
type Risk struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medimum  int `json:"medium"`
	Low      int `json:"low"`
	Unscored int `json:"unscored"`
}
type UpgradeSuggestion struct {
	Version           string `json:"version"`
	VersionName       string `json:"versionName"`
	Origin            string `json:"origin"`
	ExternalId        string `json:"externalId"`
	VulnerabilityRisk Risk   `json:"vulnerabilityRisk"`
}

type TransitiveUpgradeSuggestion struct {
	Component                string            `json:"component"`
	ComponentName            string            `json:"componentName"`
	ExternalId               string            `json:"externalId"`
	VersionName              string            `json:"versionName"`
	OriginExternalNamespace  string            `json:"originExternalNamespace"`
	OriginExternalId         string            `json:"originExternalId"`
	ShortTermUpgradeGuidance UpgradeSuggestion `json:"shortTermUpgradeGuidance"`
	LongTermUpgradeGuidance  UpgradeSuggestion `json:"longTermUpgradeGuidance"`
}

type RapidScanComponent struct {
	Name                           string                        `json:"componentName"`
	Version                        string                        `json:"versionName"`
	Identifier                     string                        `json:"componentIdentifier"`
	ReleaseDate                    string                        `json:"releaseDate"`
	ExternalId                     string                        `json:"externalId"`
	ExternalNamespace              string                        `json:"externalNamespace"`
	PackageUrl                     *string                       `json:"packageUrl"`
	OriginId                       string                        `json:"originId"`
	MatchTypes                     []string                      `json:"matchTypes"`
	ComponentDescription           *string                       `json:"componentDescription"`
	ViolatingPolicies              []Policy                      `json:"violatingPolicies"`
	ComponentViolatingPolicies     []Policy                      `json:"componentViolatingPolicies"`
	Vulnerabilities                []ComponentVulnerability      `json:"allVulnerabilities"`
	Licenses                       []ComponentLicense            `json:"allLicenses"`
	LicenceV2                      ComponentLicenceV2            `json:"licenses"`
	PolicyViolationVulnerabilities []ComponentVulnerability      `json:"policyViolationVulnerabilities"`
	PolicyViolationLicenses        []ComponentLicense            `json:"policyViolationLicenses"`
	PartiallyEvaluatedPolicies     []string                      `json:"partiallyEvaluatedPolicies"`
	NonEvaluatedPolicies           []string                      `json:"nonEvaluatedPolicies"`
	DependencyTrees                [][]string                    `json:"dependencyTrees"`
	ShortTermUpgradeGuidance       UpgradeSuggestion             `json:"shortTermUpgradeGuidance"`
	LongTermUpgradeGuidance        UpgradeSuggestion             `json:"longTermUpgradeGuidance"`
	TransitiveUpgradeGuidance      []TransitiveUpgradeSuggestion `json:"transitiveUpgradeGuidance"`
	Meta                           Meta                          `json:"_meta"`
}

type ComponentLicenceV2 struct {
	LicenseDisplay string               `json:"licenseDisplay"`
	LicenseType    *string              `json:"licenseType"`
	License        string               `json:"license"`
	SpdxId         *string              `json:"spdxId"`
	Licenses       []ComponentLicenceV2 `json:"licenses"`
	Ownership      string               `json:"ownership"`
	FamilyName     string               `json:"licenseFamilyName"`
}
