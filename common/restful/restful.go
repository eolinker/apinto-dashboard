package restful

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type RestArgs = string
type RequestCore[T any] interface {
	Header(name, value string) T
	Query(name, value string) T
}
type Request interface {
	RequestCore[Request]
	Request(args ...RestArgs) (*Response, error)
}
type RequestOneWay[I any] interface {
	RequestCore[RequestOneWay[I]]
	Request(input *I, args ...RestArgs) (*Response, error)
}
type RequestCall[D any] interface {
	RequestCore[RequestCall[D]]
	Request(args ...RestArgs) (*ResponseData[D], error)
}

type RequestRPC[I any, D any] interface {
	RequestCore[RequestRPC[I, D]]
	Request(input *I, args ...RestArgs) (*ResponseData[D], error)
}
type Builder[R any] interface {
	Build() R
	Reset(c Config)
}

type Config struct {
	address []string
	header  http.Header
	query   url.Values
}

func BuildConfig(header http.Header, query url.Values, address ...string) Config {
	for i, _ := range address {
		address[i] = strings.TrimSuffix(address[i], "/")
		if !(strings.HasPrefix(address[i], "http://") || strings.HasPrefix(address[i], "https://")) {
			address[i] = fmt.Sprint("http://", address[i])
		}
	}
	if query == nil {
		query = make(url.Values)
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/json")
	return Config{address: address, header: header, query: query}
}

func Simple(config Config, method string, path string) Builder[Request] {
	return &_SimpleBuilder{
		builderCore: builderCore{
			path:   createPath(path),
			method: method,
			Config: config,
		},
	}
}
func Call[D any](config Config, method, path string) Builder[RequestCall[D]] {
	return &_RequestCallBuilder[D]{
		builderCore: newBuilderCore(config, method, path),
	}
}
func OneWay[I any](config Config, method, path string) Builder[RequestOneWay[I]] {
	return &_OnewayBuilder[I]{
		builderCore: newBuilderCore(config, method, path),
	}
}
func Rpc[I any, O any](config Config, method, path string) Builder[RequestRPC[I, O]] {

	return &_RpcBuilder[I, O]{
		builderCore: newBuilderCore(config, method, path),
	}
}
func cloneValues(vs url.Values) url.Values {

	if vs == nil {
		return nil
	}

	// Find total number of values.
	nv := 0
	for _, vv := range vs {
		nv += len(vv)
	}
	sv := make([]string, nv) // shared backing array for headers' values
	h2 := make(url.Values, len(vs))
	for k, vv := range vs {
		if vv == nil {
			// Preserve nil values. ReverseProxy distinguishes
			// between nil and zero-length header values.
			h2[k] = nil
			continue
		}
		n := copy(sv, vv)
		h2[k] = sv[:n:n]
		sv = sv[n:]
	}
	return h2

}
