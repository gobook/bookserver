package main

import (
	"fmt"
	"time"

	"github.com/gobook/bookserver/conf"
	"github.com/gobook/bookserver/middlewares"
	"github.com/gobook/bookserver/models"
	"github.com/gobook/bookserver/routers"
	"github.com/lunny/tango"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"
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
	sess := session.New()
	t.Use(
		tango.Static(tango.StaticOptions{
			RootPath: conf.BooksRootPath,
			Prefix:   "read",
		}),
		xsrf.New(20*time.Minute),
		middlewares.Auth(sess),
		renders.New(renders.Options{
			Reload:    true,
			Directory: "templates",
			//Extensions: []string{".html"},
			Vars: renders.T{
				"AppUrl":  "/",
				"AppLogo": "/public/images/logo.png",
			},
		}),
		sess,
		flash.Flashes(sess),
	)

	t.Get("/", new(routers.Home))
	t.Any("/login", new(routers.Login))
	t.Get("/logout", new(routers.Logout))
	t.Any("/register", new(routers.Register))

	t.Group("/user", func(g *tango.Group) {
		//g.Get("/*name", new(routers.User))
	})

	t.Group("/myBooks", func(g *tango.Group) {
		g.Get("/", new(routers.MyBooks))
		g.Get("/:id/down", new(routers.DownBook))
		g.Get("/:id/up", new(routers.UpBook))
		g.Get("/:id/update", new(routers.UpdateBook))
		g.Get("/:id/del", new(routers.DelBook))
	})

	t.Group("/book", func(g *tango.Group) {
		//g.Get("/:id", new(routers.Book))
	})

	t.Group("/publish", func(g *tango.Group) {
		g.Any("/", new(routers.Publish))
		g.Any("/*path", new(routers.Publish))
	})

	t.Group("/auth", func(g *tango.Group) {
		g.Get("/login/github", new(routers.LoginGithub))
	})

	t.Group("/github", func(g *tango.Group) {
		g.Any("/", new(routers.GithubPush))
		//g.Get("/auth", new(routers.Github))
		//g.Get("/callback", new(routers.GithubCallback))
	})

	t.Run(conf.Listen)
}
