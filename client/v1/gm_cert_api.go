package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
)

type IGMCert interface {
	Save(name string, config *GMCertConfig) error
	Del(name string) error
}

type gmCert struct {
	IClient
	professionName string
}

func newGMCert(IClient IClient) *gmCert {
	return &gmCert{
		IClient:        IClient,
		professionName: "certificate",
	}
}

func (c *gmCert) Save(name string, config *GMCertConfig) error {
	var resErr error
	if config.Driver == "" {
		config.Driver = "gm-server"
	}
	for _, addr := range c.addrs() {
		url := fmt.Sprintf("%s/api/%s/%s", addr, c.professionName, name)
		bytes, _ := json.Marshal(config)
		res, code, err := requestDo(http.MethodPost, url, bytes)
		if err != nil {
			log.Errorf("profession-api-create err=%s", err.Error())
			continue
		}
		if code != http.StatusOK {
			log.Errorf("profession-api-create err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		return nil
	}
	return resErr
}

func (c *gmCert) Del(name string) error {
	var resErr error
	for _, addr := range c.addrs() {
		url := fmt.Sprintf("%s/api/%s/%s", addr, c.professionName, name)
		res, code, err := requestDo(http.MethodDelete, url, nil)
		if err != nil {
			log.Errorf("profession-api-del err=%s", err.Error())
			continue
		}
		if code != http.StatusOK {
			log.Errorf("profession-api-del err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		return nil
	}
	return resErr
}
