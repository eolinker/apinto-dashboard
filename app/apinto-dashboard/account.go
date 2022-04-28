package main

type Account struct {
	AccountList []*AccountItem `yaml:"account_list"`
}

type AccountItem struct {
	UserName string                 `yaml:"user_name"`
	Password string                 `yaml:"password"`
	Info     map[string]interface{} `yaml:"info"`
}

func InitUserDetails() {

}
