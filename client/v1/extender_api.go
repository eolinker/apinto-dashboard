package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
)

var _ IExtender = (*extender)(nil)

type IExtender interface {
	List() ([]*ExtenderListItem, error)
	Info(group, project, name string) (*ExtenderInfo, error)
}

type extender struct {
	IClient
}

func newIExtender(IClient IClient) IExtender {
	return &extender{IClient: IClient}
}

func (e *extender) Info(group, project, name string) (*ExtenderInfo, error) {
	var resErr error
	for _, addr := range e.addrs() {
		url := fmt.Sprintf("%s/extender/%s:%s/%s", addr, group, project, name)
		res, code, err := requestDo(http.MethodGet, url, nil)
		if err != nil {
			log.Errorf("extender-api-List err=%s", err.Error())
			resErr = err
			continue
		}
		result := new(ExtenderInfo)
		if code != http.StatusOK {
			log.Errorf("extender-api-Info err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		if err = json.Unmarshal(res, &result); err != nil {
			log.Errorf("extender-api-Info-jsonUnmarshal err=%s", err.Error())
			resErr = err
			continue
		}
		return result, nil
	}
	return nil, resErr
}

func (e *extender) List() ([]*ExtenderListItem, error) {
	var resErr error
	for _, addr := range e.addrs() {
		url := fmt.Sprintf("%s/extender", addr)
		res, code, err := requestDo(http.MethodGet, url, nil)
		if err != nil {
			log.Errorf("extender-api-List err=%s", err.Error())
			resErr = err
			continue
		}
		resultList := make([]*ExtenderListItem, 0)
		if code != http.StatusOK {
			log.Errorf("extender-api-List err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		if err = json.Unmarshal(res, &resultList); err != nil {
			log.Errorf("extender-api-List-jsonUnmarshal err=%s", err.Error())
			resErr = err
			continue
		}
		return resultList, nil
	}
	return nil, resErr
}
