package routers

import (
	"errors"
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

	err = models.UpdateById(id, &models.Book{Status: BookOffline}, "status")
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

	err = gobook.MakeBook(filepath.Join(conf.BooksRootPath, book.RepoPath), dstDir)
	if err != nil {
		return err
	}

	d.Flash.Set("info", "书籍 "+book.Name+" 更新成功")
	u.Redirect("/myBooks")
	return nil
}
