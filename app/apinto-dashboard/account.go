package main

import (
	"errors"
	"github.com/eolinker/apinto-dashboard/internal/security"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	ErrNilUserDetails = errors.New("config.yaml needs user_details. ")
	ErrFileType       = errors.New("config.yaml user_details.tyoe: only file valid. ")
	ErrNilAccountCfg  = errors.New("config.yaml user_details.file: account config is nil. ")
)

type Account struct {
	AccountList []*AccountItem `yaml:"account_list"`
}

type AccountItem struct {
	UserName string                 `yaml:"user_name"`
	Password string                 `yaml:"password"`
	Info     map[string]interface{} `yaml:"info"`
}

func InitUserDetails(detailsService *security.UserDetailsService, config *UserDetailsConfig) error {
	if config == nil {
		return ErrNilUserDetails
	}
	//暂时只支持file
	if config.Type != "file" {
		return ErrFileType
	}
	accountCfg, err := ReadAccountConfig(config.File)
	if err != nil {
		return err
	}

	if accountCfg == nil || len(accountCfg.AccountList) == 0 {
		return ErrNilAccountCfg
	}

	for _, account := range accountCfg.AccountList {
		detailsService.Add(security.NewUserDetails(account.UserName, account.Password, account.Info))
	}

	return nil
}

func ReadAccountConfig(file string) (*Account, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := new(Account)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
