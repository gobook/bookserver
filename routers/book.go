package routers

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gobook/bookserver/conf"
	"github.com/gobook/bookserver/models"
	"github.com/gobook/gobook"
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
)

type MyBooks struct {
	AuthBase
}

func (m *MyBooks) Get() error {
	books, err := models.FindBooksByUserId(m.LoginUserId())
	if err != nil {
		return err
	}
	return m.Render("mybooks.html", renders.T{
		"books": books,
	})
}

type DownBook struct {
	AuthBase
	flash.Flash
}

func (d *DownBook) Get() error {
	id := d.ParamInt64(":id")
	var book models.Book
	has, err := models.GetById(id, &book)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("书籍不存在")
	}

	if book.AuthorId != d.LoginUserId() {
		return errors.New("没有权限")
	}

	err = models.UpdateById(id, &models.Book{Status: models.BookOffline}, "status")
	if err != nil {
		return err
	}
	d.Flash.Set("info", "书籍 "+book.Name+" 下架成功")
	d.Redirect("/myBooks")
	return nil
}

type UpdateBook struct {
	AuthBase
	flash.Flash
}

type BookType int

const (
	UnknowBook BookType = iota
	GoBookType
)

func detectBookType(githubPath string) (BookType, error) {
	url := "https://" + githubPath + "/blob/master/SUMMARY.md"
	resp, err := http.Get(url)
	if err != nil {
		return UnknowBook, err
	}
	if resp.StatusCode == http.StatusOK {
		return GoBookType, nil
	}
	return UnknowBook, nil
}

func (u *UpdateBook) Get() error {
	id := u.ParamInt64(":id")
	var book models.Book
	has, err := models.GetById(id, &book)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("书籍不存在")
	}

	if book.AuthorId != u.LoginUserId() {
		return errors.New("没有权限")
	}

	dstDir := filepath.Join(conf.ReposRootPath, book.RepoPath)
	err = downGithubBook(dstDir, book.RepoPath)
	if err != nil {
		return err
	}

	_, err = gobook.MakeBook(filepath.Join(conf.BooksRootPath, book.RepoPath), dstDir)
	if err != nil {
		return err
	}

	u.Flash.Set("info", "书籍 "+book.Name+" 更新成功")
	u.Redirect("/myBooks")
	return nil
}

type DelBook struct {
	AuthBase
	flash.Flash
}

func (d *DelBook) Get() error {
	id := d.ParamInt64(":id")
	var book models.Book
	has, err := models.GetById(id, &book)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("书籍不存在")
	}
	os.RemoveAll(filepath.Join(conf.BooksRootPath, book.RepoPath))
	os.RemoveAll(filepath.Join(conf.ReposRootPath, book.RepoPath))
	err = models.DelById(id, new(models.Book))
	if err != nil {
		return err
	}
	d.Flash.Set("info", "书籍 "+book.RepoPath+" 删除成功")
	d.Redirect("/myBooks")
	return nil
}
