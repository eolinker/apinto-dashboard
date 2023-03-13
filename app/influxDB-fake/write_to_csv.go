package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-basic/uuid"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

type Employee struct {
	ID  string
	Age int
}

const (
	CsvFileDir         = "./export/csv"
	CsvRequestFileName = "test_request_data-%s.csv"
	CsvProxyFileName   = "test_proxy_data-%s.csv"
)

func writeToCsv() error {
	err := os.MkdirAll(CsvFileDir, 0777)
	if err != nil {
		return err
	}

	csvRequestFileName := fmt.Sprintf(CsvRequestFileName, time.Now().Format("2006-01-02T15:04:05"))
	requestFile, err := os.OpenFile(path.Join(CsvFileDir, csvRequestFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	defer requestFile.Close()
	if err != nil {
		return err
	}
	requestWriter := csv.NewWriter(requestFile)
	defer requestWriter.Flush()

	csvProxyFileName := fmt.Sprintf(CsvProxyFileName, time.Now().Format("2006-01-02T15:04:05"))
	proxyFile, err := os.OpenFile(path.Join(CsvFileDir, csvProxyFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	defer proxyFile.Close()
	if err != nil {
		return err
	}
	proxyWriter := csv.NewWriter(proxyFile)
	defer proxyWriter.Flush()

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
	ipCount := 50
	ipList := getIPs(ipCount)

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
	proxySitutationList := getProxySituationList()
	proxySituaCount := len(proxySitutationList)

	rand.Seed(time.Now().UnixNano())

	//写请求表的表头
	requestHeader1 := []string{"#group", "FALSE", "FALSE", "FALSE", "FALSE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE"}
	requestHeader2 := []string{"#datatype", "string", "long", "dateTime:RFC3339", "long", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string"}
	requestHeader3 := []string{"#default", "_result", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
	requestHeader4 := []string{"", "result", "table", "_time", "_value", "_field", "_measurement", "api", "app", "cluster", "handler", "host", "ip", "method", "node", "path", "request_id", "status", "upstream"}
	requestTableNo := 0
	if err = requestWriter.Write(requestHeader1); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	if err = requestWriter.Write(requestHeader2); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	if err = requestWriter.Write(requestHeader3); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	if err = requestWriter.Write(requestHeader4); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	//写转发表的表头
	proxyHeader1 := []string{"#group", "FALSE", "FALSE", "FALSE", "FALSE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE", "TRUE"}
	proxyHeader2 := []string{"#datatype", "string", "long", "dateTime:RFC3339", "long", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string", "string"}
	proxyHeader3 := []string{"#default", "_result", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
	proxyHeader4 := []string{"", "result", "table", "_time", "_value", "_field", "_measurement", "addr", "api", "app", "cluster", "host", "index", "method", "node", "path", "request_id", "request_ip", "request_path", "status", "upstream"}
	proxyTableNo := 0
	if err = proxyWriter.Write(proxyHeader1); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	if err = proxyWriter.Write(proxyHeader2); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	if err = proxyWriter.Write(proxyHeader3); err != nil {
		log.Fatalln("error writing record to file", err)
	}
	if err = proxyWriter.Write(proxyHeader4); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	for t := startTime; endTime.After(t); t = t.Add(timeInterval) {
		fmt.Println("正在写入，时间：", t.Format("2006-01-02 15:04:05"))
		for i := 0; i < concurrencyPerInterval; i++ {
			//now := time.Now()

			random := rand.Intn(apiCount)
			apiInfo := apiList[random]

			ip := ipList[rand.Intn(ipCount)] //请求ip
			app := appList[rand.Intn(appCount)]

			requestID := uuid.New()
			requestTraffic := rand.Intn(3000) + 300

			proxySituation := proxySitutationList[rand.Intn(proxySituaCount)]

			timing := 0
			responseTraffic := 0

			//写转发信息
			for p := 0; p < proxySituation.retry+1; p++ {

				proxyTiming := rand.Intn(30) + 1
				timing += proxyTiming
				proxyResponse := rand.Intn(3000) + 300

				proxyFieldTiming := []string{"", "", strconv.Itoa(proxyTableNo), t.Add(time.Duration(p+1+i*1000) * time.Nanosecond).Format(time.RFC3339Nano), strconv.Itoa(proxyTiming),
					"timing", "proxy", "172.0.0.2", apiInfo.APIUUID + "@router", app.Id + "@app", clusterID, "172.0.0.1", fmt.Sprintf("%d", p+1),
					"POST", node, apiInfo.RequestPath, requestID, ip, apiInfo.RequestPath, proxySituation.respStatusList[p], apiInfo.ServiceName + "@service"}
				proxyTableNo++
				proxyFieldRequest := []string{"", "", strconv.Itoa(proxyTableNo), t.Add(time.Duration(p+1+i*1000) * time.Nanosecond).Format(time.RFC3339Nano), strconv.Itoa(requestTraffic),
					"request", "proxy", "172.0.0.2", apiInfo.APIUUID + "@router", app.Id + "@app", clusterID, "172.0.0.1", fmt.Sprintf("%d", p+1),
					"POST", node, apiInfo.RequestPath, requestID, ip, apiInfo.RequestPath, proxySituation.respStatusList[p], apiInfo.ServiceName + "@service"}
				proxyTableNo++
				proxyFieldResponse := []string{"", "", strconv.Itoa(proxyTableNo), t.Add(time.Duration(p+1+i*1000) * time.Nanosecond).Format(time.RFC3339Nano), strconv.Itoa(proxyResponse),
					"response", "proxy", "172.0.0.2", apiInfo.APIUUID + "@router", app.Id + "@app", clusterID, "172.0.0.1", fmt.Sprintf("%d", p+1),
					"POST", node, apiInfo.RequestPath, requestID, ip, apiInfo.RequestPath, proxySituation.respStatusList[p], apiInfo.ServiceName + "@service"}
				proxyTableNo++

				if err = proxyWriter.Write(proxyFieldTiming); err != nil {
					log.Fatalln("error writing record to file", err)
				}
				if err = proxyWriter.Write(proxyFieldRequest); err != nil {
					log.Fatalln("error writing record to file", err)
				}
				if err = proxyWriter.Write(proxyFieldResponse); err != nil {
					log.Fatalln("error writing record to file", err)
				}

			}

			//请求表写入信息
			fieldTiming := []string{"", "", strconv.Itoa(requestTableNo), t.Add(time.Duration(proxySituation.retry+1+i*1000) * time.Nanosecond).Format(time.RFC3339Nano), strconv.Itoa(timing),
				"timing", "request", apiInfo.APIUUID + "@router", app.Id + "@app", clusterID, "proxy", "172.0.0.1", ip, "POST",
				node, apiInfo.RequestPath, requestID, proxySituation.reqStatus, apiInfo.ServiceName + "@service"}
			requestTableNo++
			fieldRequest := []string{"", "", strconv.Itoa(requestTableNo), t.Add(time.Duration(proxySituation.retry+1+i*1000) * time.Nanosecond).Format(time.RFC3339Nano), strconv.Itoa(requestTraffic),
				"request", "request", apiInfo.APIUUID + "@router", app.Id + "@app", clusterID, "proxy", "172.0.0.1", ip, "POST",
				node, apiInfo.RequestPath, requestID, proxySituation.reqStatus, apiInfo.ServiceName + "@service"}
			requestTableNo++
			fieldResponse := []string{"", "", strconv.Itoa(requestTableNo), t.Add(time.Duration(proxySituation.retry+1+i*1000) * time.Nanosecond).Format(time.RFC3339Nano), strconv.Itoa(responseTraffic),
				"response", "request", apiInfo.APIUUID + "@router", app.Id + "@app", clusterID, "proxy", "172.0.0.1", ip, "POST",
				node, apiInfo.RequestPath, requestID, proxySituation.reqStatus, apiInfo.ServiceName + "@service"}
			requestTableNo++
			fieldRetry := []string{"", "", strconv.Itoa(requestTableNo), t.Add(time.Duration(proxySituation.retry+1+i*1000) * time.Nanosecond).Format(time.RFC3339Nano), strconv.Itoa(proxySituation.retry),
				"retry", "request", apiInfo.APIUUID + "@router", app.Id + "@app", clusterID, "proxy", "172.0.0.1", ip, "POST",
				node, apiInfo.RequestPath, requestID, proxySituation.reqStatus, apiInfo.ServiceName + "@service"}
			requestTableNo++

			if err = requestWriter.Write(fieldTiming); err != nil {
				log.Fatalln("error writing record to file", err)
			}
			if err = requestWriter.Write(fieldRequest); err != nil {
				log.Fatalln("error writing record to file", err)
			}
			if err = requestWriter.Write(fieldResponse); err != nil {
				log.Fatalln("error writing record to file", err)
			}
			if err = requestWriter.Write(fieldRetry); err != nil {
				log.Fatalln("error writing record to file", err)
			}
			//fmt.Println(time.Now().Sub(now).String())
		}
	}

	return nil
}
