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
	Group          string            `toml:"group"`
	Name           string            `toml:"name"`
	Region         string            `toml:"region"`
	Ip             string            `toml:"ip"`
	NodeId         string            `toml:"nodeId"`
	Env            AppEnv            `toml:"env"`
	Custom         map[string]string `toml:"custom"`
	FullAppName    string
	HttpPort       int      `toml:"httpPort"`
	RpcPort        int      `toml:"rpcPort"`
	KafkaHookAddrs []string `toml:"kafkaHookAddrs"`
}

//mysql配置
type MysqlConf struct {
	DbName   string `toml:"dbName"`
	Addr     string `toml:"addr"`
	ReadAddr string `toml:"readAddr"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Charset  string `toml:"charset"`
	//sql执行警告阈值，毫秒
	WarnThreshold time.Duration `toml:"warnThreshold"`
}

//redis配置
type RedisConf struct {
	Addr      string   `toml:"addr"`
	Password  string   `toml:"password"`
	DbIndex   int      `toml:"dbIndex"`
	Sentinels []string `toml:"sentinels"` //redis sentinel列表
}

//mongo配置
type MongoConf struct {
	Addr     string `toml:"addr"`
	DbName   string `toml:"dbName"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

//consul配置
type ConsulConf struct {
	Addrs []string `toml:"addrs"`
}

//kafka配置
type KafkaConf struct {
	Addrs []string `toml:"addrs"`
	Topic string   `toml:"topic"`
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
