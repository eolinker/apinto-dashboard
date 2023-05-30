package dto

type EmailInput struct {
	Uuid     string `json:"uuid"`
	SmtpUrl  string `json:"smtp_url"`
	SmtpPort int    `json:"smtp_port"`
	Protocol string `json:"protocol"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type EmailOutput struct {
	Uuid     string `json:"uuid"`
	SmtpUrl  string `json:"smtp_url"`
	SmtpPort int    `json:"smtp_port"`
	Protocol string `json:"protocol"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
}
