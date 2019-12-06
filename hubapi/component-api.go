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

type ComponentList struct {
	ItemsListBase
	Items []ComponentVariant `json:"items"`
}

type ComponentVariant struct {
	ComponentName string `json:"componentName"`
	VersionName   string `json:"versionName,omitempty"`
	OriginID      string `json:"originId,omitempty"`
	Component     string `json:"component"`
	Version       string `json:"version,omitempty"`
	Variant       string `json:"variant,omitempty"`
}

type ComponentVersionList struct {
	ItemsListBase
	Items []ComponentVersion `json:"items"`
}

type ComponentVersion struct {
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
	ItemsListBase
	Items []ComponentVersionOrigin `json:"items"`
}

type ComponentVersionOrigin struct {
	ComponentVersion
	Origin   string `json:"originName"`
	OriginID string `json:"originId"`
}

type ComponentVersionVulnerabilityList struct {
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
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Homepage            string   `json:"url,omitempty"`
	AdditionalHomepages []string `json:"additionalHomepages"`
	ApprovalStatus      string   `json:"approvalStatus"`
	Type                string   `json:"type"`
}

type ComponentRemediation struct {
	FixesPreviousVulnerabilities *RemediationInfo `json:"fixesPreviousVulnerabilities,omitempty"`
	NoVulnerabilities            *RemediationInfo `json:"noVulnerabilities,omitempty"`
	LatestAfterCurrent           *RemediationInfo `json:"latestAfterCurrent,omitempty"`
	Meta                         Meta            `json:"_meta"`
}

type RemediationInfo struct {
	Name               string    `json:"name"`
	ComponentVersion   string    `json:"componentVersion"`
	ReleasedOn         time.Time `json:"releasedOn"`
	VulnerabilityCount int       `json:"vulnerabilityCount"`
}
