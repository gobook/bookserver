package middlewares

import (
	"github.com/lunny/log"
	"github.com/lunny/tango"
	"github.com/tango-contrib/session"

	"github.com/gobook/bookserver/models"
)

var (
	LoginIdKey   = "auth_user_id"
	LoginNameKey = "auth_user_name"
)

type auther interface {
	AskAuth() bool
	IsLogin() bool
	SetUserSession(*session.Session)
}

var _ auther = &Auther{}

type Auther struct {
	s *session.Session
}

func (Auther) AskAuth() bool {
	return true
}

func (a *Auther) SetUserSession(s *session.Session) {
	a.s = s
}

func (a *Auther) SetLogin(u *models.User) {
	a.s.Set(LoginIdKey, u.Id)
	a.s.Set(LoginNameKey, u.Name)
}

func (a *Auther) Logout() {
	a.s.Release()
}

func (a *Auther) LoginUserId() int64 {
	return a.s.Get(LoginIdKey).(int64)
}

func (a *Auther) LoginUserName() string {
	return a.s.Get(LoginNameKey).(string)
}

func (a *Auther) LoginUser() *models.User {
	var user models.User
	has, err := models.GetById(a.LoginUserId(), &user)
	if err != nil {
		log.Error("LoginUser:", err)
		return nil
	}

	if !has {
		log.Error("LoginUser: no exist", a.LoginUserId())
		return nil
	}

	return &user
}

func (a *Auther) IsLogin() bool {
	v := a.s.Get(LoginIdKey)
	return v != nil && v.(int64) > 0
}

type AuthUser struct {
	Auther
}

func (AuthUser) AskAuth() bool {
	return false
}

func Auth(sessions *session.Sessions) tango.HandlerFunc {
	return func(ctx *tango.Context) {
		if auther, ok := ctx.Action().(auther); ok {
			auther.SetUserSession(sessions.Session(ctx.Req(), ctx.ResponseWriter))
			if auther.AskAuth() && !auther.IsLogin() {
				ctx.Redirect("/login")
				return
			}
		}

		ctx.Next()
	}
}
