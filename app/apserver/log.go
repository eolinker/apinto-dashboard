package main

import (
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

func init() {
	formatter := &log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	}

	fileWriter := filelog.NewFileWriteByPeriod()

	period, err := filelog.ParsePeriod(GetLogPeriod())
	if err != nil {
		panic(err)
	}

	fileWriter.Set(GetLogDir(), GetLogName(), period, GetLogExpire())
	fileWriter.Open()
	var writer io.Writer = fileWriter
	writer = ToCopyToIoWriter(os.Stdout, fileWriter)

	transport := log.NewTransport(writer, GetLogLevel())
	transport.SetFormatter(formatter)
	log.Reset(transport)
}
