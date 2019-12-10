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

type bdJsonComponentDetailV4 struct{}

func (bdJsonComponentDetailV4) GetMimeType() string {
	return "application/vnd.blackducksoftware.component-detail-4+json"
}

type bdJsonComponentDetailV5 struct{}

func (bdJsonComponentDetailV5) GetMimeType() string {
	return "application/vnd.blackducksoftware.component-detail-5+json"
}

type ComponentList struct {
	bdJsonComponentDetailV4
	ItemsListBase
	Items []ComponentVariant `json:"items"`
}

type ComponentVariant struct {
	bdJsonComponentDetailV4
	ComponentName string `json:"componentName"`
	VersionName   string `json:"versionName,omitempty"`
	OriginID      string `json:"originId,omitempty"`
	Component     string `json:"component"`
	Version       string `json:"version,omitempty"`
	Variant       string `json:"variant,omitempty"`
}

type ComponentVersionList struct {
	bdJsonComponentDetailV5
	ItemsListBase
	Items []ComponentVersion `json:"items"`
}

type ComponentVersion struct {
	bdJsonComponentDetailV5
	VersionName         string         `json:"versionName,omitempty"`
	ReleasedOn          time.Time      `json:"releasedOn,omitempty"`
	License             ComplexLicense `json:"license"`
	Source              string         `json:"source"`
	Type                string         `json:"type"`
	AdditionalHomepages []string       `json:"additionalHomepages"`
	ApprovalStatus      string         `json:"approvalStatus"`
	Notes               string         `json:"notes,omitempty"`
	Meta                Meta           `json:"_meta"`
}

type ComponentVersionOriginList struct {
	bdJsonComponentDetailV5
	ItemsListBase
	Items []ComponentVersionOrigin `json:"items"`
}

type ComponentVersionOrigin struct {
	ComponentVersion
	Origin   string `json:"originName"`
	OriginID string `json:"originId"`
}

type ComponentVersionVulnerabilityList struct {
	bdJsonVulnerabilityV4
	ItemsListBase
	Items []ComponentVersionVulnerability `json:"items"`
}

type ComponentVersionVulnerability struct {
	VulnerabilityBase
	AccessVector          string `json:"accessVector"`
	AccessComplexity      string `json:"accessComplexity"`
	Authentication        string `json:"authentication"`
	ConfidentialityImpact string `json:"confidentialityImpact"`
	IntegrityImpact       string `json:"integrityImpact"`
	AvailabilityImpact    string `json:"availabilityImpact"`
	Meta                  Meta   `json:"_meta"`
}

// returned by "references" component meta link
type ComponentProjectReferenceList struct {
	bdJsonApplicationJson // TODO get the type, nothing in documentation
	ItemsListBase
	Items []ComponentProjectReference
}

type ComponentProjectReference struct {
	Distribution      string `json:"distribution"`
	Phase             string `json:"phase"`
	ProjectName       string `json:"projectName"`
	ProjectTier       int    `json:"projectTier"`
	ProjectUrl        string `json:"projectUrl"`
	ProjectVersionUrl string `json:"projectVersionUrl"`
	VersionName       string `json:"versionName"`
}

type Component struct {
	bdJsonComponentDetailV4
	Name                string   `json:"name"`
	Description         string   `json:"description,omitempty"`
	ApprovalStatus      string   `json:"approvalStatus"`
	Homepage            string   `json:"url,omitempty"`
	AdditionalHomepages []string `json:"additionalHomepages"`
	PrimaryLanguage     string   `json:"primaryLanguage,omitempty"`
	Source              string   `json:"source"`
	Type                string   `json:"type"`
	Notes               string   `json:"notes,omitempty"`
	Meta                Meta     `json:"_meta"`
}

type ComponentRequest struct {
	bdJsonComponentDetailV4
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Homepage            string   `json:"url,omitempty"`
	AdditionalHomepages []string `json:"additionalHomepages"`
	ApprovalStatus      string   `json:"approvalStatus"` // [UNREVIEWED, IN_REVIEW, REVIEWED, APPROVED, LIMITED_APPROVAL, REJECTED, DEPRECATED]
	Type                string   `json:"type"`
}

type ComponentRemediation struct {
	bdJsonComponentDetailV5
	FixesPreviousVulnerabilities *RemediationInfo `json:"fixesPreviousVulnerabilities,omitempty"`
	NoVulnerabilities            *RemediationInfo `json:"noVulnerabilities,omitempty"`
	LatestAfterCurrent           *RemediationInfo `json:"latestAfterCurrent,omitempty"`
	Meta                         Meta             `json:"_meta"`
}

type RemediationInfo struct {
	Name               string     `json:"name"`
	ComponentVersion   string     `json:"componentVersion"`
	ReleasedOn         *time.Time `json:"releasedOn"`
	VulnerabilityCount int        `json:"vulnerabilityCount"`
}
