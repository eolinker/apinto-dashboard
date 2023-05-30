package apinto_module

import (
	"encoding/json"
	"net/http"
	"net/url"
)

var (
	_ MiddlewareResponseWriter = (*MiddlewareResponse)(nil)
)

type MiddlewareResponse struct {
	Abort  bool        `json:"abort"`
	Action string      `json:"action"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`

	StatusCode  int            `json:"code"`
	ContentType string         `json:"content_type"`
	Keys        map[string]any `json:"addkey"`
}

func (m *MiddlewareResponse) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 0,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	if v := cookie.String(); v != "" {
		m.AddHeader("Set-Cookie", v)
	}
}

func (m *MiddlewareResponse) SetHeader(key string, value string) {
	if m.Header == nil {
		m.Header = make(http.Header)
	}
	m.Header.Set(key, value)
}
func (m *MiddlewareResponse) AddHeader(key string, value string) {
	if m.Header == nil {
		m.Header = make(http.Header)
	}
	m.Header.Add(key, value)
}
func (m *MiddlewareResponse) JSON(code int, obj any) {
	if obj != nil {

		data, err := json.Marshal(obj)
		if err == nil {
			m.Data(code, "application/json", data)
			return
		}

	}
	m.Data(code, "application/json", []byte("{}"))

}

func (m *MiddlewareResponse) Redirect(code int, location string) {
	m.StatusCode = code
	m.Action = "redirect"
	m.Body = []byte(location)
}

func (m *MiddlewareResponse) SetAbort(abort bool) {
	m.Abort = abort
}

func (m *MiddlewareResponse) Set(key string, v any) {
	if m.Keys == nil {
		m.Keys = make(map[string]any)
	}

	m.Keys[key] = v
}

func (m *MiddlewareResponse) Data(code int, contentType string, data []byte) {
	m.ContentType = contentType
	m.StatusCode = code
	m.Body = data
	m.Action = ""
}

type MiddlewareResponseWriter interface {
	Data(code int, contentType string, data []byte)
	JSON(code int, obj any)
	Set(key string, v any)
	SetHeader(key string, value string)
	AddHeader(key string, value string)
	SetAbort(abort bool)
	Redirect(code int, location string)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
}

func DecodeMiddlewareResponse(data interface{}) (*MiddlewareResponse, error) {
	//var ginCtx *gin.Context
	//ginCtx.Set()
	return DecodeFor[MiddlewareResponse](data)
}
func DecodeFor[T any](data interface{}) (*T, error) {
	v := new(T)
	err := decodeObj(data, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
func decodeObj[T any](src interface{}, out *T) error {
	if o, ok := src.(*T); ok {
		*out = *o
		return nil
	}
	var err error
	var data []byte
	switch v := src.(type) {
	case *string:
		data = []byte(*v)
	case *[]byte:
		data = *v
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}
