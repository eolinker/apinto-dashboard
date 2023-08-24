package plugin

import (
	"bytes"
	"net/http"
)

type responseWriter struct {
	buf    *bytes.Buffer
	header http.Header
	status int
}

func newResponseWriter() *responseWriter {
	return &responseWriter{
		header: make(http.Header),
		buf:    bytes.NewBuffer(nil),
	}
}

func (r *responseWriter) Header() http.Header {
	return r.header
}

func (r *responseWriter) Write(data []byte) (int, error) {
	return r.buf.Write(data)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}
