package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
)

type IProfession[T any, K any] interface {
	List() ([]*K, error)
	//Get(name string) (*T, error)
	Create(t T) error
	Update(name string, t T) error
	Patch(string, map[string]interface{}) error
	Delete(name string) error
}

type profession[T any, K any] struct {
	IClient
	professionName string
}

func (p *profession[T, K]) Create(t T) error {
	var resErr error
	for _, addr := range p.addrs() {
		url := fmt.Sprintf("%s/api/%s", addr, p.professionName)
		bytes, _ := json.Marshal(t)
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

func newIProfession[T, K any](client IClient, professionName string) IProfession[T, K] {
	return &profession[T, K]{IClient: client, professionName: professionName}
}

func (p *profession[T, K]) Patch(name string, maps map[string]interface{}) error {

	bytes, _ := json.Marshal(maps)
	var resErr error
	for _, node := range p.addrs() {
		url := fmt.Sprintf("%s/api/%s/%s", node, p.professionName, name)
		res, code, err := requestDo(http.MethodPatch, url, bytes)
		if err != nil {
			log.Errorf("profession-api-patch error=%s", err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("profession-api-patch error=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		return nil
	}
	return resErr
}

func (p *profession[T, K]) List() ([]*K, error) {
	var resErr error
	for _, addr := range p.addrs() {
		url := fmt.Sprintf("%s/api/%s", addr, p.professionName)
		res, code, err := requestDo(http.MethodGet, url, nil)

		if err != nil {
			log.Errorf("profession-api-list url=%s err=%s", url, err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("profession-api-list url=%s,err=%s", url, string(res))
			resErr = errors.New(string(res))
			continue
		}

		infos := make([]*K, 0)

		if err = json.Unmarshal(res, &infos); err != nil {
			log.Errorf("profession-json.Unmarshal err=%s", err.Error())
			resErr = err
			continue
		}

		return infos, nil
	}
	return nil, resErr
}

func (p *profession[T, K]) Update(name string, t T) error {
	var resErr error
	for _, addr := range p.addrs() {
		url := fmt.Sprintf("%s/api/%s/%s", addr, p.professionName, name)
		bytes, _ := json.Marshal(t)
		res, code, err := requestDo(http.MethodPut, url, bytes)
		if err != nil {
			log.Errorf("profession-api-update err=%s", err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("profession-api-update err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		return nil
	}
	return resErr
}

func (p *profession[T, K]) Get(name string) (*T, error) {
	var resErr error
	for _, addr := range p.addrs() {
		url := fmt.Sprintf("%s/api/%s/%s", addr, p.professionName, name)
		res, code, err := requestDo(http.MethodGet, url, nil)

		if err != nil {
			log.Errorf("profession-api-get error=%s", err)
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("profession-api-get error=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}

		info := new(T)
		if err = json.Unmarshal(res, info); err != nil {
			log.Errorf("profession-json.Unmarshal error=%s", err)
			resErr = err
			continue
		}

		return info, nil
	}
	return nil, resErr
}

func (p *profession[T, K]) Delete(name string) error {
	var resErr error
	for _, addr := range p.addrs() {
		url := fmt.Sprintf("%s/api/%s/%s", addr, p.professionName, name)
		res, code, err := requestDo(http.MethodDelete, url, nil)
		if err != nil {
			log.Errorf("profession-api-delete err=%s", err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("profession-api-delete err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		return nil
	}
	return resErr
}
