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
