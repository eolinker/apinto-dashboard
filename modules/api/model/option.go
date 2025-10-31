package apimodel

type ApiOptionItem struct {
	Uuid  string `json:"uuid"`
	Title string `json:"title"`
	//Service     string `json:"service"`
	Group       string `json:"group"`
	RequestPath string `json:"request_path"`
}
