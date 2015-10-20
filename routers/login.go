package routers

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
	"golang.org/x/oauth2"
	authGithub "golang.org/x/oauth2/github"

	"github.com/gobook/bookserver/conf"
	"github.com/gobook/bookserver/models"
)

type Register struct {
	NoAuthBase
	flash.Flash
	xsrf.Checker
}

func (r *Register) Get() error {
	gn := r.Session.Get("githubName")
	ge := r.Session.Get("githubEmail")
	var githubName, githubEmail string
	if gn != nil {
		githubName = gn.(string)
		githubEmail = ge.(string)
	}
	return r.Render("register.html", renders.T{
		"githubName":  githubName,
		"githubEmail": githubEmail,
		"xsrf_html":   r.XsrfFormHtml(),
		"from": func() string {
			if len(githubName) > 0 {
				return "/publish"
			}
			return "/"
		}(),
	})
}

func (r *Register) Post() error {
	name, email := r.Form("name"), r.Form("email")
	passwd, repasswd := r.Form("passwd"), r.Form("repasswd")
	githubName := r.Form("githubName")
	from := r.Form("from")
	if passwd != repasswd {
		r.Flash.Set("error", "两次输入密码不一致")
		r.Redirect("/register")
		return nil
	}

	err := models.Insert(&models.User{
		Name:       name,
		Email:      email,
		Passwd:     passwd,
		GithubName: githubName,
	})
	if err != nil {
		r.Flash.Set("error", "注册失败")
		r.Redirect("/register")
		return nil
	}

	r.Redirect(from)
	return nil
}

type Login struct {
	NoAuthBase
	xsrf.Checker
	flash.Flash
}

func (l *Login) Get() error {
	return l.Render("login.html", renders.T{
		"XsrfFormHtml": l.XsrfFormHtml(),
	})
}

func (l *Login) Post() error {
	name, passwd := l.Form("name"), l.Form("passwd")
	err := models.CheckPasswd(name, passwd)
	if err != nil {
		l.Set("error", "账号或密码错误")
	}
	l.Redirect("/")
	return nil
}

type LoginGithub struct {
	NoAuthBase
}

func (l *LoginGithub) Get() error {
	code := l.Form("code")
	conf := &oauth2.Config{
		ClientID:     conf.GithubClientId,
		ClientSecret: conf.GithubClientSecret,
		Endpoint:     authGithub.Endpoint,
	}
	if len(code) == 0 {
		url := conf.AuthCodeURL("")
		fmt.Println(url)
		l.Redirect(url)
		return nil
	}
	fmt.Println(code)
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}
	fmt.Println(tok)

	tc := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(tok))

	client := github.NewClient(tc)
	l.SetGithubClient(client)
	user, _, err := client.Users.Get("")
	if err != nil {
		return err
	}
	fmt.Println(user)

	u, err := models.GetUserByGithubName(*user.Login)
	if err == models.ErrNotExist {
		l.Session.Set("githubName", *user.Login)
		l.Session.Set("githubEmail", *user.Email)
		l.Redirect("/register")
		return nil
	}
	if err != nil {
		return err
	}

	l.SetLogin(u)
	l.Redirect("/")
	return nil
}

type Logout struct {
	AuthBase
}

func (l *Logout) Get() error {
	l.Logout()
	l.Redirect("/")
	return nil
}
