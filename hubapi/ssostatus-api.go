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

const ContentTypeBdAdminV4 = "application/vnd.blackducksoftware.admin-4+json"

type bdJsonAdminDetailV4 struct{}

func (bdJsonAdminDetailV4) GetMimeType() string {
	return ContentTypeBdAdminV4
}

type SsoStatus struct {
	bdJsonAdminDetailV4
	SsoEnabled bool `json:"ssoEnabled"`
	Meta       Meta `json:"_meta"`
}
