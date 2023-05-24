package application_dto

type ApplicationAuthInput struct {
	Title          string          `json:"title"`
	Driver         string          `json:"driver"`
	HideCredential bool            `json:"hide_credential"`
	TokenName      string          `json:"token_name"`
	Position       string          `json:"position"`
	ExpireTime     int64           `json:"expire_time"`
	Config         AuthConfigProxy `json:"config"`
}

type ApplicationAuthListOut struct {
	Uuid           string `json:"uuid"`
	Title          string `json:"title"`
	Driver         string `json:"driver"`
	HideCredential bool   `json:"hide_credential"`
	ExpireTime     int64  `json:"expire_time"`
	Operator       string `json:"operator"`
	UpdateTime     string `json:"update_time"`
}

type ApplicationAuthOut struct {
	Title          string          `json:"title"`
	Uuid           string          `json:"uuid"`
	Driver         string          `json:"driver"`
	ExpireTime     int64           `json:"expire_time"`
	Operator       string          `json:"operator"`
	Position       string          `json:"position"`
	TokenName      string          `json:"token_name"`
	UpdateTime     string          `json:"update_time"`
	HideCredential bool            `json:"hide_credential"`
	Config         AuthConfigProxy `json:"config"`
}

type AuthConfigProxy []byte

func (c *AuthConfigProxy) String() string {
	return string(*c)
}

func (c *AuthConfigProxy) MarshalJSON() ([]byte, error) {
	return *c, nil
}

func (c *AuthConfigProxy) UnmarshalJSON(bytes []byte) error {
	*c = bytes
	return nil
}

type AuthDriversItem struct {
	Name   string `json:"name"`
	Render Render `json:"render"`
}

type Render []byte

func (r *Render) MarshalJSON() ([]byte, error) {
	return *r, nil
}
