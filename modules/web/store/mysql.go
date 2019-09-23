package store

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"github.com/710leo/urlooker/modules/web/g"
)

var Orm *xorm.Engine

func InitMysql() {
	cfg := g.Config

	var err error
	log.Println(g.Config)
	Orm, err = xorm.NewEngine("mysql", cfg.Mysql.Addr)
	if err != nil {
		log.Fatalln("fail to connect mysql", err)
	}
	Orm.SetMaxIdleConns(cfg.Mysql.Idle)
	Orm.SetMaxOpenConns(cfg.Mysql.Max)

}
