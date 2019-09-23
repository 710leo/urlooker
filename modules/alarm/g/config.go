package g

import (
	"fmt"
	"log"
	"sync"

	"github.com/toolkits/file"
	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	Debug  bool          `yaml:"debug"`
	Remain int           `yaml:"remain"` //
	Rpc    *RpcConfig    `yaml:"rpc"`
	Web    *WebConfig    `yaml:"web"`
	Alarm  *AlarmConfig  `yaml:"alarm"`
	Queue  *QueueConfig  `yaml:"queue"`
	Mysql  *MysqlConfig  `yaml:"mysql"`
	Worker *WorkerConfig `yaml:"worker"`
	Smtp   *SmtpConfig   `yaml:"smtp"`
}

type MysqlConfig struct {
	Addr string `yaml:"addr"`
	Idle int    `yaml:"idle"`
	Max  int    `yaml:"max"`
}

type RpcConfig struct {
	Listen string `yaml:"listen"`
}

type RedisConfig struct {
	Dsn          string `yaml:"dsn"`
	MaxIdle      int    `yaml:"maxIdle"`
	ConnTimeout  int    `yaml:"connTimeout"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

type AlarmConfig struct {
	Enabled      bool         `yaml:"enabled"`
	MinInterval  int64        `yaml:"minInterval"`
	QueuePattern string       `yaml:"queuePattern"`
	Redis        *RedisConfig `yaml:"redis"`
}

type WebConfig struct {
	Addrs    []string `yaml:"addrs"`
	Timeout  int      `yaml:"timeout"`
	Interval int      `yaml:"interval"`
}

type QueueConfig struct {
	Mail string `yaml:"mail"`
	Sms  string `yaml:"sms"`
}

type WorkerConfig struct {
	Sms  int `yaml:"sms"`
	Mail int `yaml:"mail"`
}

type SmtpConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
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
		return fmt.Errorf("configuration file %s is not exists", cfg)
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

	log.Println("load configuration file", cfg, "successfully")
	return nil
}
