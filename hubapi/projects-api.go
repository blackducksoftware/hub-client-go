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

import (
	"time"
)

const (
	ProjectVersionPhasePlanning    = "PLANNING"
	ProjectVersionPhaseDevelopment = "DEVELOPMENT"
	ProjectVersionPhaseReleased    = "RELEASED"
	ProjectVersionPhaseDeprecated  = "DEPRECATED"
	ProjectVersionPhaseArchived    = "ARCHIVED"
)

const (
	ProjectVersionDistributionExternal   = "EXTERNAL"
	ProjectVersionDistributionSaaS       = "SAAS"
	ProjectVersionDistributionInternal   = "INTERNAL"
	ProjectVersionDistributionOpenSource = "OPENSOURCE"
)

type ProjectList struct {
	ItemsListBase
	Items []Project `json:"items"`
}

type Project struct {
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Source                  string `json:"source"`
	ProjectTier             uint32 `json:"projectTier"`
	ProjectLevelAdjustments bool   `json:"projectLevelAdjustments"`
	ProjectOwner            string `json:"projectOwner"`
	Meta                    Meta   `json:"_meta"`
}

type ProjectRequest struct {
	Name                    string                 `json:"name"`
	Description             *string                `json:"description"`
	ProjectTier             *int                   `json:"projectTier"`
	ProjectOwner            *string                `json:"projectOwner"`
	ProjectLevelAdjustments bool                   `json:"projectLevelAdjustments"`
	VersionRequest          *ProjectVersionRequest `json:"versionRequest"`
}

type ProjectVersionList struct {
	ItemsListBase
	Items []ProjectVersion `json:"items"`
}

type ProjectVersion struct {
	VersionName     string     `json:"versionName"`
	Nickname        string     `json:"nickname"`
	ReleaseComments string     `json:"releaseComments"`
	ReleasedOn      *time.Time `json:"releasedOn"`
	Phase           string     `json:"phase"`
	Distribution    string     `json:"distribution"`
	Meta            Meta       `json:"_meta"`
}

type ProjectVersionRequest struct {
	VersionName     string     `json:"versionName"`
	Nickname        string     `json:"nickname"`
	ReleaseComments string     `json:"releaseComments"`
	ReleasedOn      *time.Time `json:"releasedOn"`
	Phase           string     `json:"phase"`
	Distribution    string     `json:"distribution"`
}

// TODO: This is horrible, we need to fix up this API
type ProjectVersionRiskProfile struct {
	Categories       map[string]map[string]int `json:"categories"`
	BomLastUpdatedAt *time.Time                `json:"bomLastUpdatedAt"`
	Meta             Meta                      `json:"_meta"`
}

// V6
type ProjectVersionPolicyStatus struct {
	OverallStatus          string                  `json:"overallStatus"`
	UpdatedAt              *time.Time              `json:"updatedAt"`
	StatusCounts           []StatusCount           `json:"componentVersionStatusCounts"`
	PolicyViolationDetails []PolicyViolationDetail `json:"componentVersionPolicyViolationDetails"`
	Meta                   Meta                    `json:"_meta"`
}

// TODO could the names and values be from an enumeration?
type StatusCount struct {
	Name  string `json:"name"` // [ IN_VIOLATION_OVERRIDDEN, NOT_IN_VIOLATION, IN_VIOLATION ]
	Value int    `json:"value"`
}

type PolicyViolationDetail struct {
	Name           string        `json:"name"` // [ IN_VIOLATION_OVERRIDDEN, NOT_IN_VIOLATION, IN_VIOLATION ]
	SeverityLevels []StatusCount `json:"severityLevels"`
}

// result of bom policy-status link under project's component version
type BomComponentPolicyStatus struct {
	ApprovalStatus string `json:"approvalStatus"`
	Meta           Meta   `json:"_meta"`
}

// result of bom policy-rules link under project's component version
type BomComponentPolicyRulesList struct {
	ItemsListBase
	Items []BomComponentPolicyRule `json:"items"`
}

type BomComponentPolicyRule struct {
	PolicyRule
	PolicyApprovalStatus string `json:"policyApprovalStatus"`
	Comment              string `json:"comment"`
}

func (p *Project) GetProjectVersionsLink() (*ResourceLink, error) {
	return p.Meta.FindLinkByRel("versions")
}

func (p *Project) GetProjectUsersLink() (*ResourceLink, error) {
	return p.Meta.FindLinkByRel("users")
}

func (v *ProjectVersion) GetProjectLink() (*ResourceLink, error) {
	return v.Meta.FindLinkByRel("project")
}

func (v *ProjectVersion) GetCodeLocationsLink() (*ResourceLink, error) {
	return v.Meta.FindLinkByRel("codelocations")
}

func (v *ProjectVersion) GetComponentsLink() (*ResourceLink, error) {
	return v.Meta.FindLinkByRel("components")
}

func (v *ProjectVersion) GetVulnerableComponentsLink() (*ResourceLink, error) {
	return v.Meta.FindLinkByRel("vulnerable-components")
}

func (v *ProjectVersion) GetProjectVersionRiskProfileLink() (*ResourceLink, error) {
	return v.Meta.FindLinkByRel("riskProfile")
}

func (v *ProjectVersion) GetProjectVersionPolicyStatusLink() (*ResourceLink, error) {
	return v.Meta.FindLinkByRel("policy-status")
}
