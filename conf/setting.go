package conf

import (
	_ "github.com/go-sql-driver/mysql"
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
	ReposRootPath = "./repos"
	BooksRootPath = "./books"
	Listen = ":8000"
	DriverName = "mysql"
	DataSourceName = "root:872a8956@/gobook?charset=utf8"
	return nil
}

func LoadConfig() error {
	return nil
}
