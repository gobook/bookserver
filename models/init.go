package models

import (
	"github.com/go-xorm/xorm"
	"github.com/gobook/bookserver/conf"
)

var (
	orm *xorm.Engine
)

func Init() error {
	var err error
	orm, err = xorm.NewEngine(conf.DriverName, conf.DataSourceName)
	if err != nil {
		return err
	}

	err = orm.Sync2(new(Book), new(Author), new(Class))
	if err != nil {
		return err
	}
	return nil
}
