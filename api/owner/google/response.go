package google

type ResponseCommon struct {
	Message string `json:"message"`
}

type ResponseAuth struct {
	Message     string `json:"message"`
	AuthUrl     string `json:"auth_url,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}
