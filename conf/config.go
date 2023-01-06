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
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"reflect"
)

type AppEnv string

const (
	AppEnv_local AppEnv = "local"
	AppEnv_coder AppEnv = "coder"
	AppEnv_dev   AppEnv = "dev"
	AppEnv_beta  AppEnv = "beta"
	AppEnv_prod  AppEnv = "prod"
)

//app运行时的一些参数
type AppRuntimeInfo struct {
	Version string
}

type Config struct {
	App          *AppConf                       `json:"app"         yaml:"app"`
	Log          *LogConf                       `json:"log"         yaml:"log"`
	ConfigCenter *ConfigCenter                  `json:"configCenter"         yaml:"configCenter"`
	Consul       *ConsulConf                    `json:"consul"         yaml:"consul"`
	Kafka        *KafkaConf                     `json:"kafka"         yaml:"kafka"`
	MongoMap     map[string]*MongoConf          `json:"mongo"         yaml:"mongo"`
	MysqlMap     map[string]*MysqlReadWriteConf `json:"mysql"         yaml:"mysql"`
	RedisMap     map[string]*RedisConf          `json:"redis"         yaml:"redis"`
	Tracer       *TracerConf                    `json:"tracer" yaml:"tracer"`
}

//应用配置
type AppConf struct {
	Group         string            `json:"group"         yaml:"group"`
	Name          string            `json:"name"         yaml:"name"`
	Region        string            `json:"region"         yaml:"region"`
	IP            string            `json:"ip"         yaml:"ip"`
	NodeId        string            `json:"nodeId"         yaml:"nodeId"`
	Env           AppEnv            `json:"env"         yaml:"env"`
	FullAppName   string            `json:"fullAppName" `
	HttpPort      int               `json:"httpPort"      yaml:"httpPort"`
	RpcPort       int               `json:"rpcPort"      yaml:"rpcPort"`
	WebSocketPort int               `json:"webSocketPort"      yaml:"webSocketPort"`
	Custom        map[string]string `json:"custom"         yaml:"custom"`
}

//日志配置
type LogConf struct {
	KafkaHookEnable bool     `json:"kafkaHookEnable"      yaml:"kafkaHookEnable"`
	KafkaHookAddrs  []string `json:"kafkaHookAddrs"      yaml:"kafkaHookAddrs"`
	KafkaHookTopic  string   `json:"kafkaHookTopic"      yaml:"kafkaHookTopic"`
	Level           string   `json:"level" yaml:"level"`
	LogoutPath      string   `json:"logoutPath" yaml:"logoutPath"`
	MaxDay          int      `json:"maxDay" yaml:"maxDay"`
}

//配置中心
type ConfigCenter struct {
	ConsulAddrs []string `json:"consulAddrs" yaml:"consulAddrs"`
	ConfKey     string   `json:"confKey" yaml:"confKey"`
	Enable      bool     `json:"enable" yaml:"enable" `
}

//consul配置
type ConsulConf struct {
	Addrs   []string `json:"addrs" yaml:"addrs"`
	ConfKey string   `json:"confKey" yaml:"confKey"`
}

//kafka配置
type KafkaConf struct {
	Addrs       []string `json:"addrs"     yaml:"addrs"`
	Topic       string   `json:"topic" yaml:"topic"`
	Enable      bool     `json:"enable" yaml:"enable" `
	DialTimeout int      `json:"dialTimeout" yaml:"dialTimeout" `
}

//mongo配置
type MongoConf struct {
	Addr     string `json:"addr"     yaml:"addr"`
	DbName   string `json:"dbName"   yaml:"dbName"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Enable   bool   `json:"enable" yaml:"enable" `
}

//mysql配置
type MysqlReadWriteConf struct {
	Master *MysqlConf `json:"master"     yaml:"master"`
	Slave  *MysqlConf `json:"slave"     yaml:"slave"`
	Enable bool       `json:"enable" yaml:"enable" `
}

type MysqlConf struct {
	DSN string `json:"dsn" yaml:"dsn"`
	//sql执行警告阈值，毫秒
	WarnThreshold int  `json:"warnThreshold" yaml:"warnThreshold"`
	DebugInfo     bool `json:"debugInfo" yaml:"debugInfo" `
}

//redis配置
type RedisConf struct {
	Addr     string             `json:"addr"      yaml:"addr"`
	Password string             `json:"password"  yaml:"password"`
	DbIndex  int                `json:"dbIndex"   yaml:"dbIndex"`
	Sentinel *RedisSentinelConf `json:"sentinel" yaml:"sentinel"` //redis sentinel列表
	Enable   bool               `json:"enable" yaml:"enable" `
}

type RedisSentinelConf struct {
	Master string   `json:"master" yaml:"master"`
	Nodes  []string `json:"nodes" yaml:"nodes"`
}

type TracerConf struct {
	Enable     bool   `json:"enable" yaml:"enable"`
	TracerType string `json:"traceType" yaml:"traceType"`
	Addr       string `json:"addr" yaml:"addr"`
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
var conf *Config

var flagConf FlagConf

func parseFlag() {
	//读取flag 配置文件默认路径
	defaultPath, _ := os.Getwd()
	flag.StringVar(&flagConf.configs, "configs", defaultPath, "configs path")
	//flag.StringVar(&flagConf.logout, "logout", defaultPath, "logout path")
	flag.StringVar(&flagConf.version, "version", "1.0.0", "input this app version, eg: -app.version=1.0.0")
	flag.StringVar(&flagConf.env, "env", "dev", "input app runtime environment,eg:dev,beta,prod")
	flag.StringVar(&flagConf.nodeId, "nodeId", "1", "input app node id, must be unique")
	flag.StringVar(&flagConf.ip, "ip", "0.0.0.0", "input app ip, for server discovery")
	flag.Parse()

	fmt.Println("-------config init console-------")
	fmt.Println("--: config path", flagConf.configs)
	//fmt.Println("--: logout path", flagConf.logout)
}

//初始化配置信息
func InitConfig() {
	parseFlag()
	//load default yaml
	if defaultConf, err := loadConfFile("application.yaml"); err != nil {
		panic(err)
	} else {
		conf = defaultConf
	}
	// load env yaml
	if envConf, err := loadConfFile(fmt.Sprintf("application.%s.yaml", flagConf.env)); err != nil {
		panic(err)
	} else {
		// 反射envConf替换conf
		defaultF := reflect.TypeOf(conf).Elem()
		defaultV := reflect.ValueOf(conf).Elem()
		envV := reflect.ValueOf(envConf).Elem()
		for i := 0; i < defaultF.NumField(); i++ {
			name := defaultF.Field(i).Name
			eValue := envV.Field(i).Interface()
			fmt.Println("--: name:", name, "val:", eValue)
			fmt.Println()
			switch name {
			case "App":
				defaultV.Field(i).Set(envV.Field(i))
			default:
				defaultV.Field(i).Set(envV.Field(i))
			}
		}
	}
	//full app name
	conf.App.FullAppName = fmt.Sprintf("%s-%s-%s", conf.App.Group,
		conf.App.Name,
		conf.App.Env)
	// flag config update
	conf.App.Env = AppEnv(flagConf.env)
	conf.App.NodeId = flagConf.nodeId
	//conf.Log.LogoutPath = flagConf.logout
	conf.App.IP = flagConf.ip
	fmt.Printf("--: load all configs success!\n")
}

//加载配置文件
func loadConfFile(fileName string) (*Config, error) {
	fullPath := fmt.Sprintf("%s/%s", flagConf.configs, fileName)
	buf, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	var out Config
	if err = yaml.Unmarshal(buf, &out); err != nil {
		return nil, fmt.Errorf("load config %s fail,err:%v", fileName, err)
	}
	return &out, nil
}

//获得配置信息
func GetConf() *Config {
	return conf
}
