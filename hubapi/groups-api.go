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

// type GroupRequest struct {
// 	UserName  string `json:"userName"`
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"lastName"`
// 	Email     string `json:"email"`
// 	Password  string `json:"password"`
// 	Active    bool   `json:"active"`
// }

type Group struct {
	Name        string `json:"name"`
	Active      bool   `json:"Active"`
	CreatedFrom string `json:"createdFrom"`
	Default     bool   `json:"default"`
	Meta        Meta   `json:"_meta"`
}

type GroupList struct {
	TotalCount     int      `json:"totalCount"`
	Items          []Group  `json:"items"`
	AppliedFilters []string `json:"appliedFilters"`
	Meta           Meta     `json:"_meta"`
}
