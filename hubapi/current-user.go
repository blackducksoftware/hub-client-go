package hubapi

type ApiToken struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Scopes      []string `json:"scopes"`
}

type ApiTokenWithMeta struct {
	ApiToken
	Meta Meta `json:"_meta"`
}

type CreateApiTokenResponse struct {
	ApiTokenWithMeta
	Token string `json:"token"`
}

type ApiTokenList struct {
	ItemsListBase
	Items []ApiTokenWithMeta `json:"items"`
}

type CurrentUserResponse struct {
	UserName         string `json:"userName"`
	ExternalUserName string `json:"externalUserName"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Email            string `json:"email"`
	Type             string `json:"type"`
	Active           bool   `json:"active"`
	User             string `json:"user"`
	Meta             Meta   `json:"_meta"`
}
