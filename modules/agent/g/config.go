package g

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/toolkits/file"
	"gopkg.in/yaml.v2"
)

type WebConf struct {
	Addrs    []string `yaml:"addrs"`
	Interval int      `yaml:"interval"`
	Timeout  int      `yaml:"timeout"`
}

type GlobalConfig struct {
	Debug    bool     `yaml:"debug"`
	Hostname string   `yaml:"hostname"`
	Worker   int      `yaml:"worker"`
	Web      *WebConf `yaml:"web"`
}

var (
	Config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Hostname() (string, error) {
	hostname := Config.Hostname
	if hostname != "" {
		return hostname, nil
	}

	return os.Hostname()
}

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
	fmt.Printf("configContent:\n%s", configContent)
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
