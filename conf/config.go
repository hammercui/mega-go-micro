/*
@Desc : infra配置
@Version : 1.0.0
@Time : 2020/8/24 13:29
@Author : hammercui
@File : config
@Company: Sdbean
*/
package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

type AppEnv string

const (
	AppEnv_local AppEnv = "local"
	AppEnv_coder AppEnv = "coder"
	AppEnv_beta  AppEnv = "beta"
	AppEnv_prod  AppEnv = "prod"
)

//app运行时的一些参数
type AppRuntimeInfo struct {
	Version string
}

type Config struct {
	App          *AppConf
	Log          *LogConf
	ConfigCenter *ConfigCenter
	Consul       *ConsulConf
	Kafka        *KafkaConf
	RedisMap     *RedisConf
	MongoMap     map[string]*MysqlConf
	MysqlMap     map[string]*MysqlConf
}

//应用配置
type AppConf struct {
	Group         string `json:"group"         yaml:"group"`
	Name          string `json:"name"         yaml:"name"`
	Region        string `json:"region"         yaml:"region"`
	IP            string `json:"ip"         yaml:"ip"`
	NodeId        string `json:"nodeId"         yaml:"nodeId"`
	Env           AppEnv `json:"env"         yaml:"env"`
	FullAppName   string
	HttpPort      int               `json:"httpPort"      yaml:"httpPort"`
	RpcPort       int               `json:"rpcPort"      yaml:"rpcPort"`
	WebSocketPort int               `json:"WebSocketPort"      yaml:"WebSocketPort"`
	Custom        map[string]string `json:"custom"         yaml:"custom"`
}

//日志配置
type LogConf struct {
	KafkaHookAddrs []string `json:"kafkaHookAddrs"      yaml:"kafkaHookAddrs"`
	LogoutPath string `json:"logoutPath" yaml:"logoutPath"`
}

//配置中心
type ConfigCenter struct {
	ConsulAddrs []string `json:"consulAddrs" yaml:"consulAddrs"`
	ConfKey     string   `json:"confKey" yaml:"confKey"`
}

//consul配置
type ConsulConf struct {
	Addrs   []string `json:"addrs" yaml:"addrs"`
	ConfKey string   `json:"confKey" yaml:"confKey"`
}

//kafka配置
type KafkaConf struct {
	Addrs []string `json:"addrs"     yaml:"addrs"`
	Topic string   `json:"topic" yaml:"topic"`
}

//mongo配置
type MongoConf struct {
	Addr     string `json:"addr"     yaml:"addr"`
	DbName   string `json:"dbName"   yaml:"dbName"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

//mysql配置
type MysqlConf struct {
	DSN string `json:"dsn" yaml:"dsn"`
	//sql执行警告阈值，毫秒
	WarnThreshold time.Duration `json:"warnThreshold" yaml:"warnThreshold"`
}

//redis配置
type RedisConf struct {
	Addr      string   `json:"addr"      yaml:"addr"`
	Password  string   `json:"password"  yaml:"password"`
	DbIndex   int      `json:"dbIndex"   yaml:"dbIndex"`
	Sentinels []RedisSentinelConf `json:"sentinel" yaml:"sentinel"` //redis sentinel列表
}

type RedisSentinelConf struct {
	Master string `json:"master" yaml:"master"`
	Nodes []string `json:"nodes" yaml:"nodes"`
}

type AppOpts struct {
	IsConfWatchOn  bool
	IsBrokerOn     bool
	IsRedisOn      bool
	IsMongoOn      bool
	IsSqlOn        bool
	IsSkyWalkingOn bool
	IsKafkaLogsOn  bool
}

//var configPath string
//var LogoutPath string
var conf Config

var flagConf FlagConf

func parseFlag()  {
	//读取flag 配置文件默认路径
	defaultConfigPath, _ := os.Getwd()
	defaultConfigPath = filepath.Dir(defaultConfigPath) + "/configs"
	//读取flag 日志输出路径
	defaultLogPath, _ := os.Getwd()
	defaultLogPath = filepath.Dir(defaultLogPath) + "/logout"
	flag.StringVar(&flagConf.configs, "configs", defaultConfigPath, "configs path")
	flag.StringVar(&flagConf.logout, "logout", defaultLogPath, "logout path")
	flag.StringVar(&flagConf.version, "app.version", "1.0.0", "input this app version, ex: -app.version=1.0.0")
	flag.StringVar(&flagConf.env, "env", "prod", "input app runtime environment,eg:dev,beta,prod")
	flag.Parse()
	fmt.Print("-------console-------")
	fmt.Println("config path", flagConf.configs)
	fmt.Println("logout path", flagConf.logout)
}

//初始化配置信息
func InitConfig() {
	parseFlag()
	//var appConf AppConf
	LoadConfFile("application.yaml", &conf.AppConf)
	//var consulConf ConsulConf
	LoadConfFile("consul.yaml", &conf.ConsulConf)
	LoadConfFile("configCenter.yaml", &conf.ConfigCenter)
	//var kafkaConf KafkaConf
	LoadConfFile("kafka.yaml", &conf.KafkaConf)
	LoadConfFile("mongo.yaml", &conf.MongoConf)
	LoadConfFile("mysql.yaml", &conf.MysqlConf)
	LoadConfFile("redis.yaml", &conf.RedisConf)

	//full app name
	conf.AppConf.FullAppName = fmt.Sprintf("%s-%s-%s", conf.AppConf.Group,
		conf.AppConf.Name,
		conf.AppConf.Env)
	fmt.Printf("load all configs success!\n")
}

func InitConfigWithOpts(opts *AppOpts) {
	//读取flag 配置文件默认路径
	defaultConfigPath, _ := os.Getwd()
	defaultConfigPath = filepath.Dir(defaultConfigPath) + "/configs"

	//读取flag 日志输出路径
	defaultLogPath, _ := os.Getwd()
	defaultLogPath = filepath.Dir(defaultLogPath) + "/logout"

	//fmt.Println("默认配置文件路径", "defaultConfigPath")
	flag.StringVar(&configPath, "configs", defaultConfigPath, "configs path")
	flag.StringVar(&LogoutPath, "logout", defaultLogPath, "logout path")
	flag.StringVar(&AppInfo.Version, "app.version", AppInfo.Version, "input this app version, ex: -app.version=1.0.0")
	flag.Parse()

	fmt.Println("Default configPath:", configPath)
	fmt.Println("Default LogoutPath", LogoutPath)

	//var appConf AppConf
	LoadConfFile("application.yaml", &conf.AppConf)
	//var consulConf ConsulConf
	LoadConfFile("consul.yaml", &conf.ConsulConf)
	if opts.IsConfWatchOn {
		LoadConfFile("configCenter.yaml", &conf.ConfigCenter)
	}
	//var kafkaConf KafkaConf
	if opts.IsBrokerOn {
		LoadConfFile("kafka.yaml", &conf.KafkaConf)
	}
	if opts.IsMongoOn {
		LoadConfFile("mongo.yaml", &conf.MongoConf)
	}
	if opts.IsSqlOn {
		LoadConfFile("mysql.yaml", &conf.MysqlConf)
	}
	if opts.IsRedisOn {
		LoadConfFile("redis.yaml", &conf.RedisConf)
	}
	//full app name
	conf.AppConf.FullAppName = fmt.Sprintf("%s-%s-%s", conf.AppConf.Group,
		conf.AppConf.Name,
		conf.AppConf.Env)
	fmt.Printf("Load all configs success!\n")
}

//加载配置文件
func LoadConfFile(fileName string, out interface{}) {
	if _, err := yaml.DecodeFile(path.Join(configPath, fileName), out); err != nil {
		log.Fatalf("load config %s fail,err:%v", fileName, err)
		os.Exit(0)
	} else {
		bytes, _ := json.Marshal(out)
		fmt.Printf("load config:%s, success! configs:%+v\n", fileName, string(bytes))
	}
}

//获得配置信息
func GetConf() *Config {
	return &conf
}
