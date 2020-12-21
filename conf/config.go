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
	"github.com/BurntSushi/toml"
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

type Config struct {
	AppConf    *AppConf
	MysqlConf  *MysqlConf
	RedisConf  *RedisConf
	MongoConf  *MysqlConf
	ConsulConf *ConsulConf
	KafkaConf  *KafkaConf
}

//应用配置
type AppConf struct {
	Group          string            `json:"group"         toml:"group"`
	Name           string            `json:"name"         toml:"name"`
	Region         string            `json:"region"         toml:"region"`
	Ip             string            `json:"ip"         toml:"ip"`
	NodeId         string            `json:"nodeId"         toml:"nodeId"`
	Env            AppEnv            `json:"env"         toml:"env"`
	Custom         map[string]string `json:"custom"         toml:"custom"`
	FullAppName    string
	HttpPort       int      `json:"httpPort"      toml:"httpPort"`
	RpcPort        int      `json:"rpcPort"      toml:"rpcPort"`
	KafkaHookAddrs []string `json:"kafkaHookAddrs"      toml:"kafkaHookAddrs"`
}

//mysql配置
type MysqlConf struct {
	DbName   string `json:"dbName" toml:"dbName"`
	Addr     string `json:"addr" toml:"addr"`
	ReadAddr string `json:"readAddr" toml:"readAddr"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
	Charset  string `json:"charset" toml:"charset"`
	//sql执行警告阈值，毫秒
	WarnThreshold time.Duration `json:"warnThreshold" toml:"warnThreshold"`
}

//redis配置
type RedisConf struct {
	Addr      string   `json:"addr"      toml:"addr"`
	Password  string   `json:"password"  toml:"password"`
	DbIndex   int      `json:"dbIndex"   toml:"dbIndex"`
	Sentinels []string `json:"sentinels" toml:"sentinels"` //redis sentinel列表
}

//mongo配置
type MongoConf struct {
	Addr     string `json:"addr"     toml:"addr"`
	DbName   string `json:"dbName"   toml:"dbName"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
}

//consul配置
type ConsulConf struct {
	Addrs   []string `json:"addrs" toml:"addrs"`
	ConfKey string   `json:"confKey" toml:"confKey"`
}

//kafka配置
type KafkaConf struct {
	Addrs []string `json:"addrs  toml:"addrs"`
	Topic string   `json:"topic" toml:"topic"`
}

var configPath string
var LogoutPath string
var conf Config

//初始化配置信息
func InitConfig() {
	//读取flag 配置文件默认路径
	defaultConfigPath, _ := os.Getwd()
	defaultConfigPath = filepath.Dir(defaultConfigPath) + "/configs"

	//读取flag 日志输出路径
	defaultLogPath, _ := os.Getwd()
	defaultLogPath = filepath.Dir(defaultLogPath) + "/logout"

	//fmt.Println("默认配置文件路径", "defaultConfigPath")
	flag.StringVar(&configPath, "configs", defaultConfigPath, "configs path")
	flag.StringVar(&LogoutPath, "logout", defaultLogPath, "logout path")
	flag.Parse()

	fmt.Println("配置文件地址", configPath)
	fmt.Println("默认日志路径", LogoutPath)

	//var appConf AppConf
	LoadConfFile("application.toml", &conf.AppConf)
	//var consulConf ConsulConf
	LoadConfFile("consul.toml", &conf.ConsulConf)
	//var kafkaConf KafkaConf
	LoadConfFile("kafka.toml", &conf.KafkaConf)
	LoadConfFile("mongo.toml", &conf.MongoConf)
	LoadConfFile("mysql.toml", &conf.MysqlConf)
	LoadConfFile("redis.toml", &conf.RedisConf)

	//full app name
	conf.AppConf.FullAppName = fmt.Sprintf("%s-%s-%s", conf.AppConf.Group,
		conf.AppConf.Name,
		conf.AppConf.Env)
	fmt.Printf("load all configs success!\n")
}

//加载配置文件
func LoadConfFile(fileName string, out interface{}) {
	if _, err := toml.DecodeFile(path.Join(configPath, fileName), out); err != nil {
		log.Fatalf("load config %s fail,err:%v", fileName, err)
		os.Exit(0)
	} else {
		fmt.Printf("load config:%s, success:%+v\n", fileName, out)
	}
}

//获得配置信息
func GetConf() *Config {
	return &conf
}
