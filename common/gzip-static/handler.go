package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Respponse struct {
	body   []byte
	header http.Header
}
type gzipHandler struct {
	*Options
	gzPool sync.Pool
	lock   sync.RWMutex
	cache  map[string]*Respponse
}

func newGzipHandler(level int, options ...Option) *gzipHandler {
	handler := &gzipHandler{
		cache:   make(map[string]*Respponse),
		Options: DefaultOptions,
		gzPool: sync.Pool{
			New: func() interface{} {
				gz, err := gzip.NewWriterLevel(io.Discard, level)
				if err != nil {
					panic(err)
				}
				return gz
			},
		},
	}
	for _, setter := range options {
		setter(handler.Options)
	}
	return handler
}

func (g *gzipHandler) Handle(c *gin.Context) {
	if fn := g.DecompressFn; fn != nil && c.Request.Header.Get("Content-Encoding") == "gzip" {
		fn(c)
	}

	if !g.shouldCompress(c.Request) {
		return
	}
	g.lock.RLock()
	data, has := g.cache[c.Request.RequestURI]
	g.lock.RUnlock()
	if has {
		for k := range data.header {
			c.Header(k, data.header.Get(k))
		}
		c.Writer.Write(data.body)
		c.Abort()
		return
	}
	g.lock.Lock()
	defer g.lock.Unlock()

	data, has = g.cache[c.Request.RequestURI]
	if has {
		for k := range data.header {
			c.Header(k, data.header.Get(k))
		}
		c.Writer.Write(data.body)
		c.Abort()
		return
	}

	gz := g.gzPool.Get().(*gzip.Writer)
	defer g.gzPool.Put(gz)
	defer gz.Reset(io.Discard)
	buffer := newBuffWriter(c.Writer)
	gz.Reset(buffer)

	c.Header("Content-Encoding", "gzip")
	c.Header("Vary", "Accept-Encoding")

	c.Writer = &gzipWriter{c.Writer, gz}
	defer func() {
		gz.Close()
		c.Header("Content-Length", fmt.Sprint(buffer.writer.Len()))
		g.cache[c.Request.RequestURI] = &Respponse{
			body:   buffer.writer.Bytes(),
			header: c.Writer.Header(),
		}
	}()
	c.Next()

}

func (g *gzipHandler) shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Accept"), "text/event-stream") {
		return false
	}

	extension := filepath.Ext(req.URL.Path)
	if g.ExcludedExtensions.Contains(extension) {
		return false
	}

	if g.ExcludedPaths.Contains(req.URL.Path) {
		return false
	}
	if g.ExcludedPathesRegexs.Contains(req.URL.Path) {
		return false
	}

	return true
}
