package hubclient

type ExternalExtension struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	InfoURL       string `json:"infoUrl"`
	Authenticated bool   `json:"authenticated"`
	Meta          Meta   `json:"_meta"`
}
