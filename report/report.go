package report

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	consoleId        string
	lastReport       time.Time
	disableLiveEvent bool
	nodeService      cluster.IClusterNodeService

	nodeCache cache.IRedisCacheNoKey[string]

	file string
)

func init() {

	bean.Autowired(&nodeService)
	nodeCache = cache.CreateRedisCacheNoKey[string](time.Minute, "nodes-report")

}
func reportLevelBanked() {

	bean.AddInitializingBeanFunc(func() {
		go func() {
			ticker := time.NewTicker(periodTime)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if !needReport() {
						continue
					}

					report(getNodes()...)

				}
			}
		}()

	})

}
func InitReport(dir string, disableLive bool) {

	disableLiveEvent = disableLive
	if !disableLiveEvent {
		go reportLevelBanked()
	}
	abs, _ := filepath.Abs(dir)

	file = filepath.Join(abs, ".console.id")
	data, err := os.ReadFile(file)

	if err != nil && os.IsNotExist(err) {
		consoleId = uuid.New()
		lastReport = time.Now()
		f, err := os.Create(file)
		if err != nil {
			return
		}
		defer f.Close()

		f.WriteString(consoleId)
		fmt.Fprint(f, consoleId, "\n", lastReport.Format(time.RFC3339))

		go func() {
			_, _ = http.Get(fmt.Sprintf("%s/report/deploy/c?id=%s", reportAddr, consoleId))
		}()
		return
	}
	ds := strings.Split(string(data), "\n")
	consoleId = ds[0]
	if len(ds) > 1 {
		lastReport, err = time.Parse(time.RFC3339, ds[1])
		if err != nil {
			lastReport = time.Now()
		}
	} else {
		lastReport = time.Now()
	}
}
func report(nodes ...string) {

	url := fmt.Sprintf("%s/report/live?id=%s", reportAddr, consoleId)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	request.Header.Set("Apinto-Nodes", strings.Join(nodes, "&"))
	_, err = http.DefaultClient.Do(request)
	if err == nil {
		UpdateReport()
	}
}

func UpdateReport() {
	lastReport = time.Now()
	data := fmt.Sprint(consoleId, "\n", lastReport.Format(time.RFC3339))
	os.WriteFile(file, []byte(data), 0666)
}
func needReport() bool {
	return !disableLiveEvent && time.Now().Day() != lastReport.Day()

}
func NeedReport() (string, string, bool) {
	if !needReport() {
		return "", "", false
	}
	nodes := getNodes()
	return consoleId, strings.Join(nodes, "&"), true
}

func getNodes() []string {
	nodeps, err := nodeCache.GetAll(context.Background())

	if err == nil && len(nodeps) == 0 {
		return nil
	}
	if len(nodeps) > 0 {

		return nodeps
	}
	nodes, err := nodeService.QueryAllCluster(context.Background())
	if err != nil {
		return nil
	}
	nodeps = make([]string, 0, len(nodes))
	for _, node := range nodes {
		nodeps = append(nodeps, node.Name)
	}

	_ = nodeCache.SetAll(context.Background(), nodeps)
	return nodeps
}
