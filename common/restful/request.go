package restful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func doRequest(servers []string, method string, pathGen PathGen, args []string, query url.Values, header http.Header, body []byte) (response *http.Response, err error) {
	var bodyReader io.Reader = nil
	path := pathGen.Gen(args...)
	for _, s := range servers {
		var req *http.Request
		rawUrl := fmt.Sprint(s, path)
		if len(body) > 0 {
			bodyReader = bytes.NewReader(body)

		}
		req, err = http.NewRequest(method, rawUrl, bodyReader)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = query.Encode()
		req.Header = header
		response, err = http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		return response, nil
	}
	return
}

func UnmarshalResponse(response *http.Response, out interface{}) error {
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}
