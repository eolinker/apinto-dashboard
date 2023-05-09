package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	"strings"
)

type clConfigRedis struct {
	apintoDriverName string
}

func CreateRedis(apintoDriverName string) cluster.ICLConfigDriver {
	return &clConfigRedis{apintoDriverName: apintoDriverName}
}

func (c *clConfigRedis) CheckInput(config []byte) error {
	redisConf := new(cluster_dto.RedisConfigInput)
	err := json.Unmarshal(config, redisConf)
	if err != nil {
		return err
	}
	if strings.TrimSpace(redisConf.Addrs) == "" {
		return errors.New("addrs can't be nil. ")
	}

	for _, addr := range strings.Split(redisConf.Addrs, ",") {
		if !common.IsMatchIpPort(addr) {
			return fmt.Errorf("addr %s is illegal. ", addr)
		}
	}

	return nil
}

func (c *clConfigRedis) ToApinto(name string, config []byte) interface{} {
	redisConf := new(cluster_dto.RedisConfigInput)
	_ = json.Unmarshal(config, redisConf)
	return &v1.RedisOutput{
		OutputConfig: v1.OutputConfig{
			Name:   name,
			Driver: c.apintoDriverName,
		},
		Addrs:    strings.Split(redisConf.Addrs, ","),
		Username: redisConf.Username,
		Password: redisConf.Password,
	}
}

func (c *clConfigRedis) FormatOut(operator string, config *cluster_entry.ClusterConfig) interface{} {
	redisConf := new(cluster_dto.RedisConfigInput)
	_ = json.Unmarshal(config.Data, redisConf)

	return &cluster_dto.RedisConfigOutput{
		Addrs:      redisConf.Addrs,
		Username:   redisConf.Username,
		Password:   redisConf.Password,
		Enable:     config.IsEnable,
		Operator:   operator,
		CreateTime: common.TimeToStr(config.CreateTime),
		UpdateTime: common.TimeToStr(config.UpdateTime),
	}
}

func (c *clConfigRedis) InitConfig(config []byte) error {
	return nil
}
