package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
)

type ICert interface {
	Save(name string, config *CertConfig) error
	Del(name string) error
}

type cert struct {
	IClient
	professionName string
}

func newCert(IClient IClient) *cert {
	return &cert{
		IClient:        IClient,
		professionName: "certificate",
	}
}

func (c *cert) Save(name string, config *CertConfig) error {
	var resErr error
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

func (c *cert) Del(name string) error {
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
