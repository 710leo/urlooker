package g

import (
	"fmt"
	"log"
	"sync"

	"github.com/toolkits/file"
	"gopkg.in/yaml.v2"
)

type LogConfig struct {
	Path     string `yaml:"path"`
	Filename string `yaml:"filename"`
	Level    string `yaml:"level"`
}

type MysqlConfig struct {
	Addr string `yaml:"addr"`
	Idle int    `yaml:"idle"`
	Max  int    `yaml:"max"`
}

type HttpConfig struct {
	Listen string `yaml:"listen"`
	Secret string `yaml:"secret"`
}

type RpcConfig struct {
	Listen string `yaml:"listen"`
}

type AlarmConfig struct {
	Enable      bool              `yaml:"enable"`
	Batch       int               `yaml:"batch"`
	Replicas    int               `yaml:"replicas"`
	ConnTimeout int               `yaml:"connTimeout"`
	CallTimeout int               `yaml:"callTimeout"`
	MaxConns    int               `yaml:"maxConns"`
	MaxIdle     int               `yaml:"maxIdle"`
	SleepTime   int               `yaml:"sleepTime"`
	Cluster     map[string]string `yaml:"cluster"`
}

type FalconConfig struct {
	Enable   bool   `yaml:"enable"`
	Addr     string `yaml:"addr"`
	Interval int    `yaml:"interval"`
}

type LdapConfig struct {
	Enabled    bool     `yaml:"enabled"`
	Addr       string   `yaml:"addr"`
	BindDN     string   `yaml:"bindDN"`
	BaseDN     string   `yaml:"baseDN`
	BindPasswd string   `yaml:"bindPasswd"`
	UserField  string   `yaml:"userField"`
	Attributes []string `yaml:attributes`
}

type InternalDnsConfig struct {
	Enable bool   `yaml:"enable"`
	CMD    string `yaml:"cmd"`
}

type GlobalConfig struct {
	Debug            bool                `yaml:"debug"`
	Admins           []string            `yaml:"admins"`
	Salt             string              `yaml:"salt"`
	Register         bool                `yaml:"register"`
	ShowDurationMin  int                 `yaml:"showDurationMin"`  //查看最近几分钟内的报警历史和绘图，默认为30分钟
	KeepDurationHour int                 `yaml:"keepDurationHour"` //保留历史数据时间长度，默认为12小时
	DNS              string              `yaml:"dns"`              //解析域名的dns服务器地址
	Http             *HttpConfig         `yaml:"http"`
	Rpc              *RpcConfig          `yaml:"rpc"`
	Ldap             *LdapConfig         `yaml:"ldap"`
	Log              *LogConfig          `yaml:"log"`
	Mysql            *MysqlConfig        `yaml:"mysql"`
	Alarm            *AlarmConfig        `yaml:"alarm"`
	Falcon           *FalconConfig       `yaml:"falcon"`
	InternalDns      *InternalDnsConfig  `yaml:"internalDns"`
	MonitorMap       map[string][]string `yaml:"monitorMap"`
}

var (
	Config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = yaml.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	configLock.Lock()
	defer configLock.Unlock()
	Config = &c
	log.Println(Config)

	log.Println("load configuration file", cfg, "successfully")
	return nil
}
