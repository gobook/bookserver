package models

import (
	"github.com/go-xorm/xorm"
	"github.com/gobook/bookserver/conf"
)

var (
	orm *xorm.Engine
)

func Init() (err error) {
	orm, err = xorm.NewEngine(conf.DriverName, conf.DataSourceName)
	if err != nil {
		return
	}

	err = orm.Sync2(new(Book), new(Class), new(User), new(Theme), new(History))
	if err != nil {
		return
	}

	err = initClasses()
	if err != nil {
		return
	}

	err = initThemes()
	return
}
