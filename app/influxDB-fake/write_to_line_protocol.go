package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

const (
	LineProtocalFileDir         = "./export/line_protocol"
	LineProtocalRequestFileName = "test_request_data.txt"
	LineProtocalProxyFileName   = "test_proxy_data.txt"
)

func writeToLineProtocol() error {
	err := os.MkdirAll(LineProtocalFileDir, 0777)
	if err != nil {
		return err
	}

	csvRequestFileName := fmt.Sprintf(LineProtocalRequestFileName, time.Now().Format("2006-01-02T15:04:05"))
	requestFile, err := os.OpenFile(path.Join(LineProtocalFileDir, csvRequestFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	defer requestFile.Close()
	if err != nil {
		return err
	}
	requestWriter := bufio.NewWriter(requestFile)
	defer requestWriter.Flush()

	//csvProxyFileName := fmt.Sprintf(LineProtocalProxyFileName, time.Now().Format("2006-01-02T15:04:05"))
	//proxyFile, err := os.OpenFile(path.Join(LineProtocalFileDir, csvProxyFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	//defer proxyFile.Close()
	//if err != nil {
	//	return err
	//}
	//proxyWriter := bufio.NewWriter(proxyFile)
	//defer proxyWriter.Flush()

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

			//写转发信息,转发次数=正常转发1次+重试次数
			for p := 0; p < proxySituation.retry+1; p++ {
				proxyTiming := rand.Intn(30) + 1
				timing += proxyTiming

				proxyResponse := rand.Intn(3000) + 300
				if p == proxySituation.retry {
					responseTraffic = proxyResponse
				}

				//proxyPoint := fmt.Sprintf("proxy,api=%s,app=%s,cluster=%s,addr=%s,host=%s,index=%d,method=%s,node=%s,path=%s,request_path=%s,status=%s,upstream=%s timing=%d,request=%d,response=%d %d\n",
				//	apiInfo.APIUUID+"@router", app.Id+"@app", clusterID, "172.0.0.2", "172.0.0.1", p+1, "POST", node, apiInfo.RequestPath, apiInfo.RequestPath, proxySituation.respStatusList[p], apiInfo.ServiceName+"@service",
				//	proxyTiming, requestTraffic, proxyResponse, t.Add(time.Duration(p+1+i*100)*time.Nanosecond).UnixNano())
				//if _, err = proxyWriter.WriteString(proxyPoint); err != nil {
				//	log.Fatalln("error writing record to file", err)
				//}
			}

			requestPoint := fmt.Sprintf("request,api=%s,app=%s,cluster=%s,handler=%s,host=%s,method=%s,node=%s,path=%s,upstream=%s timing=%di,request=%di,response=%di,retry=%di,status=%di %d\n",
				apiInfo.APIUUID+"@router", app.Id+"@app", clusterID, "proxy", "172.0.0.1", "POST", node, apiInfo.RequestPath,
				apiInfo.ServiceName+"@service", timing,
				requestTraffic, responseTraffic, proxySituation.retry, proxySituation.reqStatus, t.Add(time.Duration(proxySituation.retry+1+i*100)*time.Nanosecond).UnixNano())
			if _, err = requestWriter.WriteString(requestPoint); err != nil {
				log.Fatalln("error writing record to file", err)
			}

			//fmt.Println(time.Now().Sub(now).String())
		}
	}

	return nil
}

type newStatusSituation struct {
	reqStatus      int
	retry          int
	respStatusList []int
}

func getNewProxySituationList() []newStatusSituation {
	return []newStatusSituation{
		{
			reqStatus:      200,
			retry:          0,
			respStatusList: []int{200},
		},
		{
			reqStatus:      404,
			retry:          0,
			respStatusList: []int{404},
		},
		{
			reqStatus:      504,
			retry:          2,
			respStatusList: []int{504, 504, 504},
		},
		{
			reqStatus:      200,
			retry:          1,
			respStatusList: []int{504, 200},
		},
	}
}
