package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/monitor"
	"github.com/eolinker/apinto-dashboard/modules/monitor/model"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"strings"
)

type monitorInfluxV2 struct {
	sourceType string
}

func CreateMonitorInfluxV2(sourceType string) monitor.IMonitorSourceDriver {
	return &monitorInfluxV2{sourceType: sourceType}
}

func (m monitorInfluxV2) CheckInput(config []byte) ([]byte, error) {
	influxConf := new(model.MonitorInfluxV2Config)
	err := json.Unmarshal(config, influxConf)
	if err != nil {
		return nil, err
	}
	influxConf.Addr = strings.TrimSpace(influxConf.Addr)
	if influxConf.Addr == "" {
		return nil, errors.New("addr can't be nil. ")
	}

	if !common.IsMatchSchemeIpPort(influxConf.Addr) {
		return nil, errors.New("addr is illegal. ")
	}

	influxConf.Org = strings.TrimSpace(influxConf.Org)
	if influxConf.Org == "" {
		return nil, errors.New("org can't be nil. ")
	}

	//influxConf.Bucket = strings.TrimSpace(influxConf.Bucket)
	//if influxConf.Bucket == "" {
	//	return nil, errors.New("bucket can't be nil. ")
	//}

	influxConf.Token = strings.TrimSpace(influxConf.Token)

	//测试连接
	client := influxdb2.NewClient(influxConf.Addr, influxConf.Token)
	_, err = client.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Fail to connect InfluxDB. Err: %s ", err)
	}
	queryClient := client.QueryAPI(influxConf.Org)
	_, err = queryClient.Query(context.Background(), `
import "array"
array.from(rows: [{a :1}])
`)
	if err != nil {
		return nil, err
	}

	return json.Marshal(influxConf)
}

//func (m monitorInfluxV2) FormatOut(config []byte) ([]byte, error) {
//	panic("implement me")
//}
