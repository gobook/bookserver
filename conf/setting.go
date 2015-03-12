package conf

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/config"
)

var (
	ReposRootPath string
	BooksRootPath string
)

var (
	Listen string
)

var (
	DriverName     string
	DataSourceName string
)

func Init() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	ReposRootPath = "./repos"
	BooksRootPath = "./books"
	Listen = cfg.Get("listen")
	DriverName = "mysql"
	DataSourceName = cfg.Get("datasource_name")
	return nil
}

func LoadConfig() (*config.Config, error) {
	return config.Load("./cfg.ini")
}
