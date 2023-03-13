package dto

type RedisConfigInput struct {
	Addrs    string `json:"addrs"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RedisConfigOutput struct {
	Addrs      string `json:"addrs"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Enable     bool   `json:"enable"`
	Operator   string `json:"operator"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type InfluxV2ConfigInput struct {
	Addr string `json:"addr"`
	Org  string `json:"org"`
	//Bucket string `json:"bucket"`
	Token string `json:"token"`
}

type InfluxV2ConfigOutput struct {
	Addr string `json:"addr"`
	Org  string `json:"org"`
	//Bucket     string `json:"bucket"`
	Token      string `json:"token"`
	Enable     bool   `json:"enable"`
	Operator   string `json:"operator"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
