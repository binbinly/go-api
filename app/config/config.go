package config

import (
	"dj-api/tools"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var C *Config

func SetUp() error {
	var confPath string
	if tools.IsDev() {
		confPath = "./conf/dev.yaml"
	} else {
		confPath = tools.GetRootDir() + "/" + "conf/" + tools.GetEnv() + ".yaml"
	}
	fmt.Printf("conf path:%v\n", confPath)
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}

	C = &Config{}
	err = yaml.Unmarshal(data, C)
	if err != nil {
		return err
	}
	return nil
}

type Config struct {
	Api        ApiConfig        `yaml:"api"`
	Third      ThirdConfig      `yaml:"third"`
	Nsq        NaqConfig        `yaml:"nsq"`
	Mysql      MysqlConfig      `yaml:"mysql"`
	Redis      RedisConfig      `yaml:"redis"`
	Registry   RegisterConfig   `yaml:"registry"`
	Log        LogConfig        `yaml:"log"`
	Limit      LimitConfig      `yaml:"limit"`
	Trace      TraceConfig      `yaml:"trace"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Sentry     SentryConfig     `yaml:"sentry"`
	Slack      SlackConfig      `yaml:"slack"`
}

//api配置
type ApiConfig struct {
	Port    int    `yaml:"port"`
	RunMode string `yaml:"run_mode"`
}

//第三方接口服务配置
type ThirdConfig struct {
	Url    string `yaml:"url"`
	JwtKey string `yaml:"jwt_key"`
}

type NaqConfig struct {
	ConsumerHost []string `yaml:"consumer_host"`
	ProdHost     string   `yaml:"prod_host"`
	Topic        string   `yaml:"topic"`
	Channel      string   `yaml:"channel"`
}

type MysqlConfig struct {
	IdleConn int    `yaml:"idle_conn"`
	MaxConn  int    `yaml:"max_conn"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	Db       string `yaml:"db"`
	Prefix   string `yaml:"prefix"`
}

type RedisConfig struct {
	Db       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
	MinConn  int    `yaml:"min_conn"`
	Host     string `yaml:"host"`
	Auth     string `yaml:"auth"`
}

type RegisterConfig struct {
	Port        int    `yaml:"port"`
	ServiceName string `yaml:"service_name"`
	Name        string `yaml:"name"`
	Host        string `yaml:"host"`
}

type LogConfig struct {
	Console bool   `yaml:"console"` //是否打印到控制台
	Level   uint8  `yaml:"level"`
	Dir     string `yaml:"dir"`
}

//限流器
type LimitConfig struct {
	Enable bool `yaml:"enable"`
	Qps    int  `yaml:"qps"`
}

//跟踪器
type TraceConfig struct {
	Enable     bool   `yaml:"enable"`
	ReportPort string `yaml:"report_port"`
	SampleType string `yaml:"sample_type"`
	SampleRate int    `yaml:"sample_rate"`
}

type PrometheusConfig struct {
	Enable bool   `yaml:"enable"`
	Host   string `yaml:"host"`
}

type SentryConfig struct {
	Enable bool   `yaml:"enable"`
	Dsn    string `yaml:"dsn"`
}

type SlackConfig struct {
	Enable  bool   `yaml:"enable"`
	HookUrl string `yaml:"hook_url"`
}
