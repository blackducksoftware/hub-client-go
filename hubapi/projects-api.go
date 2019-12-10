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

type bdJsonProjectDetailV4 struct{}

func (bdJsonProjectDetailV4) GetMimeType() string {
	return "application/vnd.blackducksoftware.project-detail-4+json"
}

type ProjectList struct {
	bdJsonProjectDetailV4
	ItemsListBase
	Items []Project `json:"items"`
}

type Project struct {
	bdJsonProjectDetailV4
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	Source                  string     `json:"source"`
	ProjectTier             uint32     `json:"projectTier"`
	ProjectLevelAdjustments bool       `json:"projectLevelAdjustments"`
	ProjectOwner            string     `json:"projectOwner"`
	CreatedAt               *time.Time `json:"createdAt,omitempty"`
	CreatedBy               string     `json:"createdBy,omitempty"`
	CreatedByUser           string     `json:"createdByUser,omitempty"`
	UpdatedAt               *time.Time `json:"updatedAt,omitempty"`
	UpdatedBy               string     `json:"updatedBy,omitempty"`
	UpdatedByUser           string     `json:"updatedByUser,omitempty"`
	Meta                    Meta       `json:"_meta"`
}

type ProjectRequest struct {
	bdJsonProjectDetailV4
	Name                    string                 `json:"name"`
	Description             string                 `json:"description"`
	ProjectTier             *int                   `json:"projectTier,omitempty"`
	ProjectOwner            *string                `json:"projectOwner,omitempty"`
	ProjectLevelAdjustments bool                   `json:"projectLevelAdjustments"`
	VersionRequest          *ProjectVersionRequest `json:"versionRequest,omitempty"`
	CloneCategories         []string               `json:"cloneCategories,omitempty"` // [COMPONENT_DATA, VULN_DATA, LICENSE_TERM_FULFILLMENT]
	CustomSignatureEnabled  *bool                  `json:"customSignatureEnabled,omitempty"`
	CustomSignatureDepth    *int                   `json:"customSignatureDepth,omitempty"`
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
