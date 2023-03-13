package dto

type UserLoginInput struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Client        int    `json:"client,omitempty"`
	Type          int    `json:"type,omitempty"`
	CheckCodeType int    `json:"check_code_type,omitempty"`
	CheckCode     string `json:"check_code,omitempty"`
	AppType       int    `json:"app_type,omitempty"`
	Ref           string `json:"ref,omitempty"`
	RefUrl        string `json:"ref_url,omitempty"`
}

type UserLoginData struct {
	Jwt  string `json:"jwt,omitempty"`
	RJWT string `json:"rjwt,omitempty"`
	Type int    `json:"type,omitempty"`
}
