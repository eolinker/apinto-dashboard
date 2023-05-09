package config

const (
	MetricsAPP      Metrics = "{application}"
	MetricsAPI      Metrics = "{api}"
	MetricsService  Metrics = "{service}"
	MetricsStrategy Metrics = "{strategy}"
	MetricsIP       Metrics = "{ip}"
)

var (
	metricsTitles = map[Metrics]string{
		MetricsAPP:      "应用",
		MetricsAPI:      "API",
		MetricsService:  "上游服务",
		MetricsStrategy: "策略",
		MetricsIP:       "IP",
	}
	all = []Metrics{
		MetricsAPP, MetricsAPI, MetricsService, MetricsStrategy, MetricsIP,
	}
)

func StrategyMetrics() []Metrics {
	return all
}

type Metrics string

func (m Metrics) Name() string {
	return string(m)
}
func (m Metrics) Title() string {
	t, has := metricsTitles[m]
	if !has {
		return m.Name()
	}
	return t
}
