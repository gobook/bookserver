package routers

import (
	"github.com/gobook/bookserver/models"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type Home struct {
	tango.Compress
	renders.Renderer
}

func (h *Home) Get() error {
	books, err := models.RecentBooks()
	if err != nil {
		return err
	}

	return h.Render("home.html", renders.T{
		"books": books,
	})
}
