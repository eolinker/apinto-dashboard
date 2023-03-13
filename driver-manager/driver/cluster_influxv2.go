package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/store/flux"
	"github.com/eolinker/eosc/log"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"strings"
)

type clConfigInfluxV2 struct {
	apintoDriverName string
}

func CreateInfluxV2(apintoDriverName string) ICLConfigDriver {
	return &clConfigInfluxV2{apintoDriverName: apintoDriverName}
}

func (c *clConfigInfluxV2) CheckInput(config []byte) error {
	influxConf := new(dto.InfluxV2ConfigInput)
	err := json.Unmarshal(config, influxConf)
	if err != nil {
		return err
	}
	if strings.TrimSpace(influxConf.Addr) == "" {
		return errors.New("addr can't be nil. ")
	}

	if strings.TrimSpace(influxConf.Org) == "" {
		return errors.New("org can't be nil. ")
	}

	//if strings.TrimSpace(influxConf.Bucket) == "" {
	//	return errors.New("bucket can't be nil. ")
	//}

	//测试连接
	client := influxdb2.NewClient(influxConf.Addr, influxConf.Token)
	_, err = client.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("Fail to connect InfluxDB. Err: %s ", err)
	}
	queryClient := client.QueryAPI(influxConf.Org)
	_, err = queryClient.Query(context.Background(), `
import "array"
array.from(rows: [{a :1}])
`)
	if err != nil {
		return fmt.Errorf("Fail to connect InfluxDB. Err: %s ", err)
	}

	return nil
}

func (c *clConfigInfluxV2) ToApinto(name string, config []byte) interface{} {
	influxConf := new(dto.InfluxV2ConfigInput)
	_ = json.Unmarshal(config, influxConf)

	return &v1.InfluxV2Output{
		OutputConfig: v1.OutputConfig{
			Name:   name,
			Driver: c.apintoDriverName,
		},
		Addr:   influxConf.Addr,
		Org:    influxConf.Org,
		Bucket: "apinto", //写死apinto bucket
		Token:  influxConf.Token,
		Scopes: []string{"monitor"}, //默认
	}
}

func (c *clConfigInfluxV2) FormatOut(operator string, config *entry.ClusterConfig) interface{} {
	influxConf := new(dto.InfluxV2ConfigInput)
	_ = json.Unmarshal(config.Data, influxConf)

	return &dto.InfluxV2ConfigOutput{
		Addr: influxConf.Addr,
		Org:  influxConf.Org,
		//Bucket:     influxConf.Bucket,
		Token:      influxConf.Token,
		Enable:     config.IsEnable,
		Operator:   operator,
		CreateTime: common.TimeToStr(config.CreateTime),
		UpdateTime: common.TimeToStr(config.UpdateTime),
	}
}

func (c *clConfigInfluxV2) InitConfig(config []byte) error {
	influxConf := new(dto.InfluxV2ConfigInput)
	_ = json.Unmarshal(config, influxConf)
	client := influxdb2.NewClient(influxConf.Addr, influxConf.Token)

	orgID := ""
	//初始化bucket
	bucketsAPI := client.BucketsAPI()
	buckets, err := bucketsAPI.FindBucketsByOrgName(context.Background(), influxConf.Org)
	if err != nil {
		return err
	}

	bucketsConf := flux.GetBucketConfigList()
	//要创建的bucket
	toCreateBuckets := common.SliceToMap(bucketsConf, func(t *flux.BucketConf) string {
		return t.BucketName
	})
	if buckets != nil {
		for _, bucket := range *buckets {
			if orgID == "" && bucket.OrgID != nil {
				orgID = *bucket.OrgID
			}
			if _, has := toCreateBuckets[bucket.Name]; has {
				delete(toCreateBuckets, bucket.Name)
			}
		}
	}
	expire := domain.RetentionRuleTypeExpire
	rule := domain.RetentionRule{
		ShardGroupDurationSeconds: nil,
		Type:                      &expire,
	}
	//创建bucket
	for _, bucketConf := range toCreateBuckets {
		rule.EverySeconds = bucketConf.Retention
		_, err := client.BucketsAPI().CreateBucketWithNameWithID(context.Background(), orgID, bucketConf.BucketName, rule)
		if err != nil {
			return err
		}
		log.Infof("Create bucket %s success. organization: %s", bucketConf.BucketName, influxConf.Org)
	}

	//创建定时脚本
	tasksApi := client.TasksAPI()
	taskFilter := &api.TaskFilter{
		OrgID: orgID,
	}
	existedTasks, err := tasksApi.FindTasks(context.Background(), taskFilter)
	if err != nil {
		return err
	}
	tasksConf := flux.GetTaskConfigList()
	//要创建的bucket
	toCreateTasks := common.SliceToMap(tasksConf, func(t *flux.TaskConf) string {
		return t.TaskName
	})
	toDeleteTaskIDs := make([]string, 0, len(toCreateTasks))

	/*
		将influxDB已存在的定时脚本 与 定时脚本配置的进行对比
		1. 配置和influxDB均有则不创建
		2. 配置有，influxDB没有，则创建
		3. 配置没有，influxDB有,且是apinto开头， 则删除
	*/
	for _, task := range existedTasks {
		if _, has := toCreateTasks[task.Name]; has {
			delete(toCreateTasks, task.Name)
		} else {
			if strings.HasPrefix(task.Name, "apinto") {
				toDeleteTaskIDs = append(toDeleteTaskIDs, task.Id)
			}
		}
	}
	//删除旧的apinto定时脚本
	for _, delId := range toDeleteTaskIDs {
		err = tasksApi.DeleteTaskWithID(context.Background(), delId)
		if err != nil {
			return err
		}
	}

	//创建influxDB中没有的定时脚本
	for _, taskConf := range toCreateTasks {
		newTask := &domain.Task{
			Cron:   &taskConf.Cron,
			Flux:   taskConf.Flux,
			Name:   taskConf.TaskName,
			Offset: &taskConf.Offset,
			OrgID:  orgID,
			//Status:          nil,
		}
		_, err := tasksApi.CreateTask(context.Background(), newTask)
		if err != nil {
			return err
		}
		log.Infof("Create task %s success. organization: %s", taskConf.TaskName, influxConf.Org)
	}

	return nil
}
