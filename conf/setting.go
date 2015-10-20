package conf

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/config"
)

var (
	ReposRootPath  string
	BooksRootPath  string
	ImagesRootPath string
)

var (
	Listen string
)

var (
	DriverName     string
	DataSourceName string
)

var (
	GithubClientId     string
	GithubClientSecret string
)

func Init() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	ReposRootPath = "./repos"
	os.MkdirAll(ReposRootPath, os.ModePerm)

	BooksRootPath = "./books"
	os.MkdirAll(BooksRootPath, os.ModePerm)

	ImagesRootPath = "./images"
	os.MkdirAll(ImagesRootPath, os.ModePerm)

	Listen = cfg.Get("listen")
	DriverName = "mysql"
	DataSourceName = cfg.Get("datasource_name")
	GithubClientId = cfg.Get("github_clientid")
	GithubClientSecret = cfg.Get("github_clientsecret")
	return nil
}

func LoadConfig() (*config.Config, error) {
	return config.Load("./cfg.ini")
}
