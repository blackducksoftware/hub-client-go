// Copyright 2024 Synopsys, Inc.
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

type bdJsonAdminDetailV5 struct{}

func (bdJsonAdminDetailV5) GetMimeType() string {
	return ContentTypeBdAdminV5
}

type MfaStatus struct {
	bdJsonAdminDetailV5
	MfaEnabled bool `json:"mfaEnabled"`
	Meta       Meta `json:"_meta"`
}
