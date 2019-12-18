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

const ContentTypeBdBomV6 = "application/vnd.blackducksoftware.bill-of-materials-6+json"

type bdJsonBomV6 struct{}

func (bdJsonBomV6) GetMimeType() string {
	return ContentTypeBdBomV6
}

type BomComponentList struct {
	ItemsListBase
	bdJsonBomV6
	Items []BomComponent `json:"items"`
}

// BomComponent represents "application/vnd.blackducksoftware.bill-of-materials-6+json"
// We need to figure out what to do with content type here
type BomComponent struct {
	bdJsonBomV6
	ComponentName          string               `json:"componentName"`
	ComponentVersionName   string               `json:"componentVersionName,omitempty"`
	Component              string               `json:"component"`
	ComponentVersion       string               `json:"componentVersion,omitempty"`
	ComponentPurpose       string               `json:"componentPurpose"`
	ComponentModified      bool                 `json:"componentModified"`
	ComponentModification  string               `json:"componentModification"`
	ReleasedOn             *time.Time           `json:"releasedOn"`
	ReviewStatus           string               `json:"reviewStatus"`
	ReviewedDetails        *BomReviewDetails    `json:"reviewedDetails,omitempty"`
	PolicyStatus           string               `json:"policyStatus"`
	ApprovalStatus         string               `json:"approvalStatus"`
	Ignored                bool                 `json:"ignored"`
	ManuallyAdjusted       bool                 `json:"manuallyAdjusted"`
	Licenses               []ComplexLicense     `json:"licenses"`
	Usages                 []string             `json:"usages"`
	Origins                []BomComponentOrigin `json:"origins"`
	LicenseRiskProfile     BomRiskProfile       `json:"licenseRiskProfile"`
	VersionRiskProfile     BomRiskProfile       `json:"versionRiskProfile"`
	SecurityRiskProfile    BomRiskProfile       `json:"securityRiskProfile"`
	ActivityRiskProfile    BomRiskProfile       `json:"activityRiskProfile"`
	OperationalRiskProfile BomRiskProfile       `json:"operationalRiskProfile"`
	ActivityData           BomActivityData      `json:"activityData"`
	TotalFileMatchCount    int                  `json:"totalFileMatchCount"`
	MatchTypes             []string             `json:"matchTypes"`
	InAttributionReport    bool                 `json:"inAttributionReport"`
	AttributionStatement   string               `json:"attributionStatement"`
	Meta                   Meta                 `json:"_meta"`
}

type BomVulnerableComponentList struct {
	bdJsonBomV6
	ItemsListBase
	Items []BomVulnerableComponent `json:"items"`
}

type BomVulnerableComponent struct {
	bdJsonBomV6
	ComponentName              string                       `json:"componentName"`
	ComponentVersionName       string                       `json:"componentVersionName"`
	ComponentVersion           string                       `json:"componentVersion"`
	ComponentVersionOriginName string                       `json:"componentVersionOriginName"`
	ComponentVersionOriginID   string                       `json:"componentVersionOriginId"`
	License                    ComplexLicense               `json:"license"`
	Vulnerability              VulnerabilityWithRemediation `json:"vulnerabilityWithRemediation"`
	Meta                       Meta                         `json:"_meta"`
}

type VulnerabilityWithRemediation struct {
	VulnerabilityBase
	RemediationStatus    string     `json:"remediationStatus"`
	RemediationCreatedAt *time.Time `json:"remediationCreatedAt"`
	RemediationUpdatedAt *time.Time `json:"remediationUpdatedAt"`
}

type BomRiskProfile struct {
	Counts []BomRiskProfileItem `json:"counts"`
}

type BomRiskProfileItem struct {
	CountType string `json:"countType"`
	Count     int    `json:"count"`
}

type BomComponentOrigin struct {
	Name                          string `json:"name"`
	ExternalNamespace             string `json:"externalNamespace"`
	ExternalID                    string `json:"externalId"`
	ExternalNamespaceDistribution bool   `json:"externalNamespaceDistribution"`
	Meta                          Meta   `json:"_meta"`
}

type BomActivityData struct {
	ContributorCount int        `json:"contributorCount12Month"`
	CommitCount      int        `json:"commitCount12Month"`
	LastCommitDate   *time.Time `json:"lastCommitDate"`
	Trend            string     `json:"trending"` // [DECREASING, STABLE, INCREASING, UNKNOWN]
	NewerReleases    *int       `json:"newerReleases,omitempty"`
}

type BomReviewDetails struct {
	ReviewedAt    *time.Time       `json:"reviewedAt,omitempty"`
	ReviewedBy    string           `json:"reviewedBy,omitempty"`
	ReviewingUser BomReviewingUser `json:"reviewingUser"`
}

type BomReviewingUser struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	User      string `json:"user"`
}

// result of bom policy-status link under project's component version
type BomComponentPolicyStatus struct {
	bdJsonBomV6
	ApprovalStatus string `json:"approvalStatus"`
	Meta           Meta   `json:"_meta"`
}

// result of bom policy-rules link under project's component version
type BomComponentPolicyRulesList struct {
	bdJsonBomV6
	ItemsListBase
	Items []BomComponentPolicyRule `json:"items"`
}

type BomComponentPolicyRule struct {
	bdJsonBomV6
	PolicyRule
	PolicyApprovalStatus string `json:"policyApprovalStatus"`
	Comment              string `json:"comment"`
}
