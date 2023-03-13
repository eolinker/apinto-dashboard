package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

var (
	runConfig *Config
)

type Config struct {
	Org         string `yaml:"org"`
	Bucket      string `yaml:"bucket"`
	Token       string `yaml:"token"`
	InfluxdbUrl string `yaml:"influxdb_url"`

	ApserverAddr    string `yaml:"apserver_addr"`
	ApserverSession string `yaml:"apserver_session"`
	groupUUID       string `yaml:"group_uuid"`

	StartTime              string `yaml:"start_time"`
	EndTime                string `yaml:"end_time"`
	TimeInterval           string `yaml:"time_interval"`
	ConcurrencyPerInterval int    `yaml:"concurrency_per_interval"`
	ClusterID              string `yaml:"cluster_id"`
	Node                   string `yaml:"node"`
}

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	c := new(Config)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
	runConfig = c
}
