package main

import (
	"github.com/eolinker/apinto-dashboard/config"
	plugin_client "github.com/eolinker/apinto-dashboard/plugin/go-plugin/plugin-client"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/log/filelog"
	"io"
	"os"
)

type writes []io.Writer

func ToCopyToIoWriter(ws ...io.Writer) io.Writer {
	return writes(ws)
}
func (ws writes) Write(p []byte) (n int, err error) {
	for _, w := range ws {
		n, err = w.Write(p)
	}
	return
}

func initLog() {
	formatter := &log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	}

	fileWriter := filelog.NewFileWriteByPeriod(filelog.Config{
		Dir:    config.GetLogDir(),
		File:   config.GetLogName(),
		Expire: config.GetLogExpire(),
		Period: filelog.ParsePeriod(config.GetLogPeriod()),
	})

	writer := ToCopyToIoWriter(os.Stdout, fileWriter)

	transport := log.NewTransport(writer, config.GetLogLevel())
	plugin_client.SetLog(config.GetLogLevel().String(), writer)
	transport.SetFormatter(formatter)
	log.Reset(transport)
}
