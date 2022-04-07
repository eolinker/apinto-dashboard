package apinto

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	"io"
	"io/ioutil"
	"net/http"
)

type response struct {
	data []byte
	code int
}

func writeResult(w http.ResponseWriter, status int, data []byte) {
	w.WriteHeader(status)
	if len(data) == 0 {
		w.Write([]byte("{}"))
		return
	}
	w.Write(data)
}

func readBody(r io.ReadCloser) ([]byte, error) {
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
