package routers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-xweb/uuid"
	"github.com/gobook/bookserver/conf"
	"github.com/gobook/bookserver/models"
	"github.com/gobook/gobook"
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

	bookType, err := detectBookType(path)
	if err != nil {
		return err
	}
	if bookType != GoBookType {
		return errors.New("该工程不符合书籍格式，未找到SUMMARY.md")
	}

	return p.Render("publish2.html", renders.T{
		"path": path,
	})
}

func (p *Publish) Post() error {
	path := p.Param("*path")
	if len(path) == 0 {
		return tango.NotFound()
	}

	cover, _, err := p.Req().FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		return err
	}

	hasCover := (err != http.ErrMissingFile)

	name := p.Form("name")
	if len(name) <= 0 {
		panic("name is empty")
	}

	dstDir := filepath.Join(conf.ReposRootPath, path)
	err = downGithubBook(dstDir, path)
	if err != nil {
		return err
	}

	bk, err := gobook.MakeBook(filepath.Join(conf.BooksRootPath, path), dstDir)
	if err != nil {
		return err
	}

	var saveName string
	if hasCover {
		saveName = uuid.New()
		savePath := filepath.Join(conf.BooksRootPath, path, saveName)
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

	if len(saveName) == 0 {
		saveName = bk.Cover
	}

	book := &models.Book{
		Name:     name,
		Cover:    saveName,
		AuthorId: p.LoginUserId(),
		RepoPath: path,
	}

	var commitId string
	err = models.PublishBook(book, commitId)
	if err != nil {
		return err
	}

	return p.Render("publish_success.html", renders.T{
		"book": book,
	})
}
