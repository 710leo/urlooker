package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/710leo/urlooker/modules/web/api"
	"github.com/710leo/urlooker/modules/web/cron"
	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/http"
	"github.com/710leo/urlooker/modules/web/http/cookie"
	"github.com/710leo/urlooker/modules/web/model"
	"github.com/710leo/urlooker/modules/web/sender"
	"github.com/710leo/urlooker/modules/web/store"

	"github.com/toolkits/file"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	prepare()

	cfg := flag.String("c", "", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)
	handleConfig(*cfg)

	store.InitMysql()
	cron.Init()
	sender.Init()
	cookie.Init()
}

func main() {
	err := model.AdminRegister()
	if err != nil {
		log.Fatalln(err)
	}
	go api.Start()
	http.Start()
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleConfig(configFile string) {
	if configFile == "" {
		configFile = "configs/web.yml"
	}

	if file.IsExist("configs/web.local.yml") {
		configFile = "configs/web.local.yml"
	}

	err := g.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
