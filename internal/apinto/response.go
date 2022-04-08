package apinto

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc/log"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Data map[string]interface{} `json:"data"`
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
}

func WriteResult(w http.ResponseWriter, status int, data []byte) {
	var err error
	res := &Response{
		Code: status,
	}
	if status != http.StatusOK {
		res.Msg = string(data)
	} else {
		err = json.Unmarshal(data, &res.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	d, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(d)
}

func ReadBody(r io.ReadCloser) ([]byte, error) {
	defer func() {
		if err := r.Close(); err != nil {
			log.Errorf("failed to close body, err: %s", err.Error())
		}
	}()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read body, err: %s", err.Error())
	}
	return data, nil
}
