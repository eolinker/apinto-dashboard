package main

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/eosc/log"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	systemConfig *Config
)

type Config struct {
	Port          int            `yaml:"port"`
	UserCenterUrl string         `yaml:"user_center_url"`
	MysqlConfig   DBConfig       `yaml:"mysql"`
	ErrorLog      ErrorLogConfig `yaml:"error_log"`
	RedisConfig   RedisConfig    `yaml:"redis"`
}

type ErrorLogConfig struct {
	LogDir    string `yaml:"dir"`
	FileName  string `yaml:"file_name"`
	LogLevel  string `yaml:"log_level"`
	LogExpire string `yaml:"log_expire"`
	LogPeriod string `yaml:"log_period"`
}

type DBConfig struct {
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Db       string `yaml:"db"`
}

type RedisConfig struct {
	UserName string   `yaml:"user_name"`
	Password string   `yaml:"password"`
	Addr     []string `yaml:"addr"`
}

func init() {
	data, err := os.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}

	c := new(Config)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
	systemConfig = c
}

func GetPort() int {
	return systemConfig.Port
}

func GetUserCenterUrl() string {
	return systemConfig.UserCenterUrl
}

func GetDBUserName() string {
	return systemConfig.MysqlConfig.UserName
}

func GetDBPassword() string {
	return systemConfig.MysqlConfig.Password
}

func GetDBIp() string {
	return systemConfig.MysqlConfig.Ip
}

func GetDBPort() int {
	return systemConfig.MysqlConfig.Port
}

func GetDbName() string {
	return systemConfig.MysqlConfig.Db
}

func getRedisAddr() []string {
	return systemConfig.RedisConfig.Addr
}

func getRedisUserName() string {
	return systemConfig.RedisConfig.UserName
}

func getRedisPwd() string {
	return systemConfig.RedisConfig.Password
}

func GetLogDir() string {
	logDir := systemConfig.ErrorLog.LogDir
	if logDir == "" {
		//默认路径是可执行程序的上一层目录的 work/logs 根据系统自适应
		lastDir, err := common.GetLastAbsPathByExecutable()
		if err != nil {
			panic(err)
		}
		logDir = fmt.Sprintf("%s%swork%slog", lastDir, string(os.PathSeparator), string(os.PathSeparator))
	} else if !strings.HasPrefix(logDir, string(os.PathSeparator)) {
		//若目录配置不为绝对路径, 则路径为 上一层目录路径 + 配置的目录路径
		lastDir, err := common.GetLastAbsPathByExecutable()
		if err != nil {
			panic(err)
		}
		relativePathPrefix := fmt.Sprintf("..%s", string(os.PathSeparator))
		logDir = path.Join(lastDir, strings.TrimPrefix(logDir, relativePathPrefix))
	}
	dirPath, err := filepath.Abs(logDir)
	if err != nil {
		panic(err)
	}
	return dirPath
}

func GetLogName() string {
	if systemConfig.ErrorLog.FileName == "" {
		return "error.log"
	}
	return systemConfig.ErrorLog.FileName
}

func GetLogLevel() log.Level {
	l, err := log.ParseLevel(systemConfig.ErrorLog.LogLevel)
	if err != nil {
		l = log.InfoLevel
	}
	return l
}

func GetLogExpire() time.Duration {
	if strings.HasSuffix(systemConfig.ErrorLog.LogExpire, "h") {
		d, err := time.ParseDuration(systemConfig.ErrorLog.LogExpire)
		if err != nil {
			return 7 * time.Hour
		}
		return d
	}
	if strings.HasSuffix(systemConfig.ErrorLog.LogExpire, "d") {

		d, err := strconv.Atoi(strings.Split(systemConfig.ErrorLog.LogExpire, "d")[0])
		if err != nil {
			return 7 * 24 * time.Hour
		}
		return time.Duration(d) * 24 * time.Hour
	}
	return 7 * 24 * time.Hour
}

func GetLogPeriod() string {
	if systemConfig.ErrorLog.LogPeriod == "" {
		return "day"
	}
	return systemConfig.ErrorLog.LogPeriod
}
