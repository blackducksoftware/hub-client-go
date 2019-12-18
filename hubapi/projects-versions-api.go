package hubapi

import "time"

// items related to  /api/projects/{projectId}/versions endpoint

const ContentTypeBdProjectDetailV5 = "application/vnd.blackducksoftware.project-detail-5+json"

type bdJsonProjectDetailV5 struct{}

func (bdJsonProjectDetailV5) GetMimeType() string {
	return ContentTypeBdProjectDetailV5
}

type ProjectVersionList struct {
	bdJsonProjectDetailV5
	ItemsListBase
	Items []ProjectVersion `json:"items"`
}

type ProjectVersion struct {
	bdJsonProjectDetailV5
	VersionName          string          `json:"versionName"`
	Nickname             string          `json:"nickname"`
	ReleaseComments      string          `json:"releaseComments"`
	ReleasedOn           *time.Time      `json:"releasedOn"`
	Phase                string          `json:"phase"`
	Distribution         string          `json:"distribution"`
	License              *ComplexLicense `json:"license,omitempty"`
	CreatedAt            *time.Time      `json:"createdAt"`
	CreatedBy            string          `json:"createdBy"`
	CreatedByUser        string          `json:"createdByUser"`
	SettingUpdatedAt     *time.Time      `json:"settingUpdatedAt"`
	SettingUpdatedBy     string          `json:"settingUpdatedBy"`
	SettingUpdatedByUser string          `json:"settingUpdatedByUser"`
	Source               string          `json:"source"`
	Meta                 Meta            `json:"_meta"`
}

type ProjectVersionRequest struct {
	bdJsonProjectDetailV5
	VersionName         string           `json:"versionName"`
	Nickname            string           `json:"nickname,omitempty"`
	ReleaseComments     string           `json:"releaseComments,omitempty"`
	ReleasedOn          *time.Time       `json:"releasedOn,omitempty"`
	Phase               string           `json:"phase"`
	Distribution        string           `json:"distribution"`
	License             []ComplexLicense `json:"license,omitempty"`
	CloneFromReleaseUrl string           `json:"cloneFromReleaseUrl,omitempty"`
}

// TODO: This is horrible, we need to fix up this API
type ProjectVersionRiskProfile struct {
	bdJsonBomV6
	Categories       map[string]map[string]int `json:"categories"`
	BomLastUpdatedAt *time.Time                `json:"bomLastUpdatedAt"`
	Meta             Meta                      `json:"_meta"`
}

type ProjectVersionPolicyStatus struct {
	bdJsonBomV6
	OverallStatus          string                  `json:"overallStatus"`
	UpdatedAt              *time.Time              `json:"updatedAt"`
	StatusCounts           []StatusCount           `json:"componentVersionStatusCounts"`
	PolicyViolationDetails *PolicyViolationDetails `json:"componentVersionPolicyViolationDetails"`
	Meta                   Meta                    `json:"_meta"`
}

// TODO could the names and values be from an enumeration?
type StatusCount struct {
	Name  string `json:"name"` // [ IN_VIOLATION_OVERRIDDEN, NOT_IN_VIOLATION, IN_VIOLATION ]
	Value int    `json:"value"`
}

type PolicyViolationDetails struct {
	Name           string        `json:"name"` // [ IN_VIOLATION_OVERRIDDEN, NOT_IN_VIOLATION, IN_VIOLATION ]
	SeverityLevels []StatusCount `json:"severityLevels"`
}
