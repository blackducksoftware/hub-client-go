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

// base type that matches json response from api/snippet-matching
type SnippetMatchResponse struct {
	SnippetMatches SnippetMatch   `json:"snippetMatches,omitempty"`
	Meta           meta           `json:"_meta,omitempty"`
	LogRef         string         `json:"logRef,omitempty"`
	ErrorMessage   string         `json:"errorMessage,omitempty"`
	Arguments      arguments      `json:"arguments,omitempty"`
	ErrorCode      string         `json:"errorCode,omitempty"`
	Errors         []snippetError `json:"errors,omitempty"`
}

type licenseDefinition struct {
	Name               string `json:"name,omitempty"`
	SpdxID             string `json:"spdxId,omitempty"`
	Ownership          string `json:"ownership,omitempty"`
	LicenseDisplayName string `json:"licenseDisplayName,omitempty"`
}

type Regions struct {
	SourceStartLines  []int `json:"sourceStartLines,omitempty"`
	SourceEndLines    []int `json:"sourceEndLines,omitempty"`
	MatchedStartLines []int `json:"matchedStartLines,omitempty"`
	MatchedEndLines   []int `json:"matchedEndLines,omitempty"`
}

func (r Regions) GetSourceStartLine() int32 {
	if len(r.SourceStartLines) == 0 {
		return -1
	}
	return int32(r.SourceStartLines[0])
}

func (r Regions) GetSourceEndLine() int32 {
	if len(r.SourceEndLines) == 0 {
		return -1
	}
	return int32(r.SourceEndLines[0])
}

func (r Regions) GetMatchedStartLine() int32 {
	if len(r.MatchedStartLines) == 0 {
		return -1
	}
	return int32(r.MatchedStartLines[0])
}

func (r Regions) GetMatchedEndLine() int32 {
	if len(r.MatchedEndLines) == 0 {
		return -1
	}
	return int32(r.MatchedEndLines[0])
}

type Snippet struct {
	ProjectName       string            `json:"projectName,omitempty"`
	ReleaseVersion    string            `json:"releaseVersion,omitempty"`
	LicenseDefinition licenseDefinition `json:"licenseDefinition,omitempty"`
	MatchedFilePath   string            `json:"matchedFilePath,omitempty"`
	Regions           Regions           `json:"regions,omitempty"`
}

type SnippetMatch struct {
	Unknown               []Snippet `json:"UNKNOWN,omitempty"`
	Reciprocal            []Snippet `json:"RECIPROCAL,omitempty"`
	ReciprocalAGPL        []Snippet `json:"RECIPROCAL_AGPL,omitempty"`
	WeakReciprocal        []Snippet `json:"WEAK_RECIPROCAL,omitempty"`
	RestrictedProprietary []Snippet `json:"RESTRICTED_PROPRIETARY,omitempty"`
	InternalProprietary   []Snippet `json:"INTERNAL_PROPRIETARY,omitempty"`
	Permissive            []Snippet `json:"PERMISSIVE,omitempty"`
}

type meta struct {
	Href  string `json:"href,omitempty"`
	Links []struct {
		Rel  string `json:"rel,omitempty"`
		Href string `json:"href,omitempty"`
	} `json:"links,omitempty"`
}
type arguments struct {
	NumCharacters string `json:"numCharacters,omitempty"`
	MinCharacters string `json:"minCharacters,omitempty"`
	MaxCharacters string `json:"maxCharacters,omitempty"`
}

type snippetError struct {
	LogRef       string `json:"logRef,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Arguments    struct {
		NumCharacters string `json:"numCharacters,omitempty"`
		MinCharacters string `json:"minCharacters,omitempty"`
		MaxCharacters string `json:"maxCharacters,omitempty"`
	} `json:"arguments,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
}
