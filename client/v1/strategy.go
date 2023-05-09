package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
)

type IStrategy interface {
	Batch(string, []interface{}) error
}

type strategy struct {
	IClient
	name string
}

func (s *strategy) Batch(batchSettingName string, ts []interface{}) error {
	bytes, _ := json.Marshal(ts)
	var resErr error
	for _, node := range s.addrs() {
		url := fmt.Sprintf("%s/setting/%s", node, batchSettingName)
		res, code, err := requestDo(http.MethodPost, url, bytes)
		if err != nil {
			log.Errorf("strategy-api-batch error=%s", err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("strategy-api-batch error=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}
		return nil
	}
	return resErr
}

func newIStrategy(client IClient) IStrategy {
	s := &strategy{IClient: client}
	return s
}
