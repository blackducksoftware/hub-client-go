package hubapi

import (
	"net/url"
	"path"
)

const (
	AuthenticateApi      = "/api/tokens/authenticate"
	CodeLocationsApi     = "/api/codelocations"
	ComponentsApi        = "/api/components"
	CurrentUserApi       = "/api/current-user"
	CurrentUserTokensApi = "/api/current-user/tokens"
	CurrentVersionApi    = "/api/current-version"
	DetectUriApi         = "/api/external-config/detect-uri"
	PolicyRulesApi       = "/api/policy-rules"
	ProjectsApi          = "/api/projects"
	DeveloperScansApi    = "/api/developer-scans"
	SsoStatusApi         = "/api/sso/status"
	UsersApi             = "/api/users"
	VulnerabilitiesApi   = "/api/vulnerabilities"
	ReadinessApi         = "/api/health-checks/readiness"
	LivenessApi          = "/api/health-checks/liveness"
	FullResultsApi       = "/full-result"
	SecurityApi          = "/j_spring_security_check"
	RemediatingApi       = "/remediating"
)

func BuildUrl(urlBase string, api string) string {
	baseUrl, err := url.Parse(urlBase)
	if err != nil {
		return ""
	}

	baseUrl.Path = path.Join(baseUrl.Path, api)

	return baseUrl.String()
}

func AddParameters(urlBase string, params map[string]string) string {
	baseUrl, err := url.Parse(urlBase)
	if err != nil {
		return ""
	}

	isFirst := true
	var parameters string

	for key, value := range params {
		if isFirst {
			parameters = "?"
			isFirst = false
		} else {
			parameters = parameters + "&"
		}
		parameters = parameters + key + "=" + value
	}

	parametersUrl, err := url.Parse(parameters)
	if err != nil {
		return ""
	}

	url := baseUrl.ResolveReference(parametersUrl)
	return url.String()
}
