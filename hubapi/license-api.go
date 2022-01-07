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

type ComplexLicense struct {
	Name                 string           `json:"name,omitempty"`
	Ownership            string           `json:"ownership,omitempty"`
	LicenseType          string           `json:"type,omitempty"`
	LicenseDisplay       string           `json:"licenseDisplay"`
	Licenses             []ComplexLicense `json:"licenses"`
	License              string           `json:"license,omitempty"` // License URL
	SpdxId               string           `json:"spdxId,omitempty"`  // The ID of the license in the SPDX project’s database, if available
	LicenseFamilySummary *ResourceLink    `json:"licenseFamilySummary,omitempty"`
}

type LicenseDetails struct {
	bdJsonComponentDetailV5
	Name           string       `json:"name"`
	LicenseFamily  ResourceLink `json:"licenseFamily"`
	Ownership      string       `json:"ownership"`
	Notes          string       `json:"notes"`
	ExpirationDate *time.Time   `json:"expirationDate"`

	CreatedAt *time.Time `json:"createdAt"`
	CreatedBy *User      `json:"createdBy"`

	UpdatedAt *time.Time `json:"updatedAt"`
	UpdatedBy *User      `json:"updatedBy"`

	LicenseSource string `json:"licenseSource"`
	SpdxId        string `json:"spdxId,omitempty"` // The ID of the license in the SPDX project’s database, if available

	LicenseStatus   string     `json:"licenseStatus"`
	StatusUpdatedAt *time.Time `json:"statusUpdatedAt"`
	StatusUpdatedBy *User      `json:"statusUpdatedBy"`
}

type License struct {
	Name          string `json:"name"`
	Ownership     string `json:"ownership"`
	CodeSharing   string `json:"codeSharing"`
	Meta          Meta   `json:"_meta"`
	LicenseSource string `json:"licenseSource,omitempty"`
	LicenseStatus string `json:"licenseStatus,omitempty"`
}
