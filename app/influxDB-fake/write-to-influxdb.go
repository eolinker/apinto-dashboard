package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type statusSituation struct {
	reqStatus      string
	retry          int
	respStatusList []string
}

func WriteFakeDataToInfluxDB() error {
	org := runConfig.Org
	bucket := runConfig.Bucket
	token := runConfig.Token
	url := runConfig.InfluxdbUrl

	writeAPI, err := newWriteAPI(url, token, org, bucket)
	if err != nil {
		return err
	}

	startTime, err := parseTime(runConfig.StartTime)
	if err != nil {
		return fmt.Errorf("parse start_time fail. err:%s", err)
	}

	endTime, err := parseTime(runConfig.EndTime)
	if err != nil {
		return fmt.Errorf("parse end_time fail. err:%s", err)
	}

	timeInterval, err := parseTimeInterval(runConfig.TimeInterval)
	if err != nil {
		return fmt.Errorf("parse time_interval fail. err:%s", err)
	}

	clusterID := runConfig.ClusterID
	node := runConfig.Node

	concurrencyPerInterval := runConfig.ConcurrencyPerInterval //每次间隔的并发写入次数

	//随机生成ip列表
	//ipCount := 50
	//ipList := getIPs(ipCount)

	//获取api信息列表
	apiList, err := getApiInfosFromFile()
	if err != nil {
		return err
	}
	apiCount := len(apiList)

	//获取app信息列表
	appList, err := getAppItems()
	if err != nil {
		return err
	}
	appCount := len(appList)

	//获取转发场景列表
	proxySitutationList := getNewProxySituationList()
	proxySituaCount := len(proxySitutationList)

	rand.Seed(time.Now().UnixNano())

	for t := startTime; endTime.After(t); t = t.Add(timeInterval) {
		fmt.Println("正在写入，时间：", t.Format("2006-01-02 15:04:05"))
		for i := 0; i < concurrencyPerInterval; i++ {
			//now := time.Now()

			random := rand.Intn(apiCount)
			apiInfo := apiList[random]

			//ip := ipList[rand.Intn(ipCount)] //请求ip
			app := appList[rand.Intn(appCount)]

			//requestID := uuid.New()
			requestTraffic := rand.Intn(3000) + 300

			proxySituation := proxySitutationList[rand.Intn(proxySituaCount)]

			timing := 0
			responseTraffic := 0
			for p := 0; p < proxySituation.retry+1; p++ {
				proxyTiming := rand.Intn(30) + 1
				timing += proxyTiming
				proxyResponse := rand.Intn(3000) + 300
				proxyInfo := write.NewPoint("proxy", map[string]string{
					"node":     node,
					"cluster":  clusterID,
					"app":      app.Id,
					"api":      apiInfo.APIUUID,
					"upstream": apiInfo.ServiceName,
					"method":   "POST",
					"host":     "172.0.0.1",
					"addr":     "172.0.0.2",
					//"path":     apiInfo.RequestPath,
					//"request_id":   requestID,
					//"index": fmt.Sprintf("%d", p+1),
					//"request_ip":   ip,
					//"request_path": apiInfo.RequestPath,
				}, map[string]interface{}{
					"timing":   proxyTiming,
					"request":  requestTraffic,
					"response": proxyResponse,
					"status":   proxySituation.respStatusList[p],
					"index":    p + 1,
				}, time.Now())
				//若是最后一次转发，将转发的response赋予请求的response
				if p == proxySituation.retry {
					responseTraffic = proxyResponse
				}

				if err = writeAPI.WritePoint(context.Background(), proxyInfo); err != nil {
					log.Println(err)
				}
			}

			requestInfo := write.NewPoint("request", map[string]string{
				"node":    node,
				"cluster": clusterID,
				"host":    "172.0.0.1",
				"method":  "POST",
				//"path":    apiInfo.RequestPath,
				//"ip":         ip,
				"api":      apiInfo.APIUUID,
				"upstream": apiInfo.ServiceName,
				"handler":  "proxy",
				//"request_id": requestID,
				"app": app.Id,
			}, map[string]interface{}{
				"timing":   timing,
				"request":  requestTraffic,
				"response": responseTraffic,
				"retry":    proxySituation.retry,
				"status":   proxySituation.reqStatus,
			}, time.Now())

			if err = writeAPI.WritePoint(context.Background(), requestInfo); err != nil {
				log.Println(err)
			}

			//fmt.Println(time.Now().Sub(now).String())
		}
	}

	return nil
}

func newWriteAPI(url, token, org, testBucket string) (api.WriteAPIBlocking, error) {
	client := influxdb2.NewClient(url, token)
	_, err := client.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Fail to connect InfluxDB. Err: %s ", err)
	}
	writeAPI := client.WriteAPIBlocking(org, testBucket)
	return writeAPI, nil
}

func getIPs(num int) []string {
	ips := make([]string, 0, num)
	for i := 1; i <= num; i++ {
		ips = append(ips, fmt.Sprintf("10.1.%d.%d", rand.Intn(256), rand.Intn(256)))
	}
	return ips
}

func getProxySituationList() []statusSituation {
	return []statusSituation{
		{
			reqStatus:      "200",
			retry:          0,
			respStatusList: []string{"200"},
		},
		{
			reqStatus:      "404",
			retry:          0,
			respStatusList: []string{"404"},
		},
		{
			reqStatus:      "504",
			retry:          2,
			respStatusList: []string{"504", "504", "504"},
		},
		{
			reqStatus:      "200",
			retry:          1,
			respStatusList: []string{"504", "200"},
		},
	}
}

func parseTime(timeStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02T15:04:05.000", timeStr, time.Local)
}

func parseTimeInterval(interval string) (time.Duration, error) {
	if strings.HasSuffix(interval, "s") {
		interval = strings.TrimSuffix(interval, "s")
		number, err := strconv.Atoi(interval)
		if err != nil {
			return 0, fmt.Errorf("fail to parse %s, err:%s", interval, err)
		}
		return time.Duration(number) * time.Second, nil

	} else if strings.HasSuffix(interval, "m") {
		interval = strings.TrimSuffix(interval, "m")
		number, err := strconv.Atoi(interval)
		if err != nil {
			return 0, fmt.Errorf("fail to parse %s, err:%s", interval, err)
		}
		return time.Duration(number) * time.Minute, nil

	} else if strings.HasSuffix(interval, "h") {
		interval = strings.TrimSuffix(interval, "h")
		number, err := strconv.Atoi(interval)
		if err != nil {
			return 0, fmt.Errorf("fail to parse %s, err:%s", interval, err)
		}
		return time.Duration(number) * time.Hour, nil

	} else if strings.HasSuffix(interval, "d") {
		interval = strings.TrimSuffix(interval, "d")
		number, err := strconv.Atoi(interval)
		if err != nil {
			return 0, fmt.Errorf("fail to parse %s, err:%s", interval, err)
		}
		return time.Duration(number) * 24 * time.Hour, nil

	}

	return 0, fmt.Errorf("fail to parse %s", interval)
}
