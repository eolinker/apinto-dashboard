package apinto_module

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

var (
	_ io.ReadCloser = (*RepeatReader)(nil)
)

type RepeatReader struct {
	reader    io.Reader
	closer    io.Closer
	buf       bytes.Buffer
	read      func(p []byte) (n int, err error)
	lastError error
}

func (r *RepeatReader) Read(p []byte) (n int, err error) {
	return r.read(p)
}

func (r *RepeatReader) readOrg(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if n > 0 {
		r.buf.Write(p[:n])
	}
	if err != nil {
		r.read = r.readCopy
		r.lastError = err
		r.reader = bytes.NewReader(r.buf.Bytes())

	}
	return
}
func (r *RepeatReader) readCopy(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if err != nil {
		r.reader = bytes.NewReader(r.buf.Bytes())
		err = r.lastError
	}
	return
}

func (r *RepeatReader) Close() error {
	if r.closer != nil {
		r.closer.Close()
		r.closer = nil
	}

	return nil
}

func SetRepeatReader(ginCtx *gin.Context) {
	if ginCtx.Request.Body == http.NoBody {
		return
	}

	reader := NewRepeatReader(ginCtx.Request.Body)
	ginCtx.Request.Body = reader
	ginCtx.Request.GetBody = func() (io.ReadCloser, error) {
		return reader, nil
	}
}
func NewRepeatReader(readCloser io.ReadCloser) *RepeatReader {

	r := &RepeatReader{reader: readCloser, closer: readCloser, buf: bytes.Buffer{}}
	r.read = r.readOrg
	return r
}
