package routers

import (
	"github.com/google/go-github/github"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"

	"github.com/gobook/bookserver/middlewares"
)

type Base struct {
	session.Session
	tango.Ctx
	renders.Renderer
}

func (b *Base) SetGithubClient(client *github.Client) {
	b.Session.Set("github_client", client)
}

func (b *Base) GetGithubClient() *github.Client {
	return b.Session.Get("github_client").(*github.Client)
}

type NoAuthBase struct {
	Base
	middlewares.AuthUser
}

func (b *NoAuthBase) Render(tmpl string, t renders.T) error {
	return b.Renderer.Render(tmpl, t.Merge(renders.T{
		"islogin": b.IsLogin(),
	}))
}

type AuthBase struct {
	Base
	middlewares.Auther
}

func (b *AuthBase) Render(tmpl string, t renders.T) error {
	return b.Renderer.Render(tmpl, t.Merge(renders.T{
		"islogin": b.IsLogin(),
	}))
}
