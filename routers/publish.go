package routers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-xweb/uuid"
	"github.com/gobook/bookserver/conf"
	"github.com/gobook/bookserver/models"
	"github.com/google/go-github/github"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

// Publish book
type Publish struct {
	AuthBase
}

func (p *Publish) Get() error {
	path := p.Param("*path")
	if len(path) == 0 {
		client := p.GetGithubClient()
		var repos []github.Repository
		repos1, _, err := client.Repositories.List("", nil)
		if err != nil {
			return err
		}
		repos = append(repos, repos1...)
		orgs, _, err := client.Organizations.List("", nil)
		if err != nil {
			return err
		}
		for _, org := range orgs {
			fmt.Println(org)
			if org.Name != nil {
				repos2, _, err := client.Repositories.ListByOrg(*org.Name, nil)
				if err != nil {
					return err
				}
				repos = append(repos, repos2...)
			}
		}
		return p.Render("publish.html", renders.T{
			"repositories": repos,
		})
	}

	// TODO: check repos' format

	return p.Render("publish2.html", renders.T{
		"path": path,
	})
}

func (p *Publish) Post() error {
	path := p.Param("*path")
	if len(path) == 0 {
		return tango.NotFound()
	}

	name := p.Form("name")
	cover, _, err := p.Req().FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		return err
	}

	var saveName string
	if err == nil {
		saveName = uuid.New()
		savePath := filepath.Join(conf.ImagesRootPath, saveName)
		dest, err := os.Create(savePath)
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, cover)
		if err != nil {
			return err
		}
	}

	book := &models.Book{
		Name:     name,
		Cover:    saveName,
		AuthorId: p.LoginUserId(),
		RepoPath: path,
	}

	err = models.Insert(book)
	if err != nil {
		return err
	}

	// TODO: download the book and publish it

	return p.Render("publish_success.html", renders.T{
		"book": book,
	})
}
