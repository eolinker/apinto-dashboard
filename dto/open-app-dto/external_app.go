package open_app_dto

type ExternalAppInfoInput struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Desc string `json:"desc"`
}

type ExternalAppInfoOutput struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Desc string `json:"desc"`
}
