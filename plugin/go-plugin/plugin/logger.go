package plugin

import (
	"github.com/eolinker/apinto-dashboard/config"
	"github.com/eolinker/eosc/log"
	"github.com/hashicorp/go-hclog"
	"os"
)

var (
	logger hclog.Logger
	levels = map[log.Level]hclog.Level{}
)

func initLogger() {
	for _, l := range log.AllLevels {
		levels[l] = hclog.LevelFromString(l.String())
	}
	logger = hclog.New(&hclog.LoggerOptions{
		Level:      hclog.LevelFromString(config.GetLogLevel().String()),
		Output:     os.Stderr,
		JSONFormat: true,
	})

	log.Reset(&hLogTransport{
		logger: logger,
		level:  config.GetLogLevel(),
	})
}

type hLogTransport struct {
	logger hclog.Logger
	level  log.Level
}

func (h *hLogTransport) Transport(entry *log.Entry) error {
	h.logger.Log(levels[entry.Level], entry.Message, entry.Data)
	return nil
}

func (h *hLogTransport) Level() log.Level {
	return h.level
}

func (h *hLogTransport) Close() error {
	return nil
}
