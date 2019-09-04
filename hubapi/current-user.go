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
	TotalCount uint32             `json:"totalCount"`
	Items      []ApiTokenWithMeta `json:"items"`
	Meta       Meta               `json:"_meta"`
}
