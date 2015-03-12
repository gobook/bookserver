package main

import (
	"fmt"

	"github.com/gobook/bookserver/conf"
	"github.com/gobook/bookserver/models"
	"github.com/gobook/bookserver/routers"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

func main() {
	err := conf.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = models.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	t := tango.Classic()
	t.Use(
		tango.Static(tango.StaticOptions{
			RootPath: conf.BooksRootPath,
			Prefix:   "read",
		}),
		renders.New(renders.Options{
			Reload:    t.Mode == tango.Dev,
			Directory: "templates",
			//Extensions: []string{".html"},
			Vars: renders.T{
				"AppUrl":  "/",
				"AppLogo": "/",
			},
		}),
	)

	t.Get("/", new(routers.Home))
	t.Any("/github", new(routers.GithubPush))
	t.Get("/github/auth", new(routers.Github))
	t.Get("/github/callback", new(routers.GithubCallback))

	t.Run(conf.Listen)
}
