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

type Component struct {
	Name                string   `json:"name"`
	Description         string   `json:"description,omitempty"`
	ApprovalStatus      string   `json:"approvalStatus"`
	Homepage            string   `json:"url,omitempty"`
	AdditionalHomepages []string `json:"additionalHomepages"`
	PrimaryLanguage     string   `json:"primaryLanguage,omitempty"`
	Source              string   `json:"source"`
	Type                string   `json:"type"`
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
