package hubapi

type RiskProfileDashboard struct {
	OverallRiskAggregate           OverallRiskAggregate           `json:"overallRiskAggregate"`
	ProjectRiskProfilePageView     ProjectList                    `json:"projectRiskProfilePageView"`
	PolicyViolationSeverityProfile PolicyViolationSeverityProfile `json:"policyViolationSeverityProfile"`
}

type OverallRiskAggregate struct {
	Categories Categories `json:"categories"`
}

type Categories struct {
	Version       Severity `json:"VERSION"`
	License       Severity `json:"LICENSE"`
	Vulnerability Severity `json:"VULNERABILITY"`
	Operational   Severity `json:"OPERATIONAL"`
	Activity      Severity `json:"ACTIVITY"`
}

type PolicyViolationSeverityProfile struct {
	Severities Severity `json:"severities"`
	Counts     []int    `json:"counts"`
}

type Severity struct {
	MAJOR       string `json:"MAJOR"`
	UNSPECIFIED string `json:"UNSPECIFIED"`
	TRIVIAL     string `json:"TRIVIAL"`
	MINOR       string `json:"MINOR"`
	CRITICAL    string `json:"CRITICAL"`
	OK          string `json:"OK"`
	BLOCKER     string `json:"BLOCKER"`
}
