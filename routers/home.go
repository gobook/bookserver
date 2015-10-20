package routers

import (
	"github.com/gobook/bookserver/models"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
)

type Home struct {
	NoAuthBase
	tango.Compress
	xsrf.Checker
}

func (h *Home) Get() error {
	books, err := models.RecentBooks()
	if err != nil {
		return err
	}

	var loginName string
	if h.IsLogin() {
		loginName = h.LoginUserName()
	}

	return h.Render("home.html", renders.T{
		"books":        books,
		"XsrfFormHtml": h.XsrfFormHtml(),
		"user":         loginName,
	})
}
