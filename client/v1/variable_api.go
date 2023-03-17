package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
)

type IVariable interface {
	Publish(namespace string, maps map[string]string) error //map[key]value
	Get(namespace string) (map[string]string, error)        //map[变量key]变量value
	GetAll() (map[string]map[string]string, error)          //map[namespace]map[key]value
	GetByKey(namespace, key string) (string, error)         //map[节点地址]变量value
}

type variable struct {
	IClient
	name string //variable
}

func (v *variable) Publish(namespace string, maps map[string]string) error {
	var resErr error
	for _, addr := range v.addrs() {
		url := fmt.Sprintf("%s/%s/%s", addr, v.name, namespace)
		bytes, _ := json.Marshal(maps)
		res, code, err := requestDo(http.MethodPost, url, bytes)
		if err != nil {
			log.Errorf("variable-api-publish err=%s", err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("variable-api-publish err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		return nil
	}
	return resErr
}

func (v *variable) Get(namespace string) (map[string]string, error) {
	var resErr error
	for _, addr := range v.addrs() {
		url := fmt.Sprintf("%s/%s/%s", addr, v.name, namespace)
		res, code, err := requestDo(http.MethodGet, url, nil)
		if err != nil {
			log.Errorf("variable-api-get err=%s", err.Error())
			resErr = err
			continue
		}
		resMap := make(map[string]string)
		if code != http.StatusOK {
			log.Errorf("variable-api-get err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		if err = json.Unmarshal(res, &resMap); err != nil {
			log.Errorf("variable-api-get-jsonUnmarshal err=%s", err.Error())
			resErr = err
			continue
		}
		return resMap, nil
	}
	return nil, resErr
}

func (v *variable) GetAll() (map[string]map[string]string, error) {
	var resErr error
	for _, addr := range v.addrs() {
		url := fmt.Sprintf("%s/%s", addr, v.name)
		res, code, err := requestDo(http.MethodGet, url, nil)
		if err != nil {
			log.Errorf("variable-api-getAll err=%s", err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("variable-api-getAll err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		resMap := make(map[string]map[string]string)
		if err = json.Unmarshal(res, &resMap); err != nil {
			log.Errorf("variable-api-getAll-jsonUnmarshal err=%s", err.Error())
			resErr = err
			continue
		}
		return resMap, nil
	}
	return nil, resErr
}

func (v *variable) GetByKey(namespace, key string) (string, error) {
	var resErr error
	for _, addr := range v.addrs() {
		url := fmt.Sprintf("%s/%s/%s/%s", addr, v.name, namespace, key)
		res, code, err := requestDo(http.MethodGet, url, nil)
		if err != nil {
			log.Errorf("variable-api-getByKey err=%s", err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("variable-api-getByKey err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		return string(res), nil
	}
	return "", resErr
}

func newIVariable(IClient IClient, name string) IVariable {
	return &variable{IClient: IClient, name: name}
}
