package routers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobook/bookserver/conf"
	"github.com/gobook/gobook"
	"github.com/lunny/tango"
)

// http://push.gobook.io/github?book=lunny%2Fxorm-manual-zh-cn&
// username=lunny&token=66940d9b-d51e-4028-99f3-a20250a024d3
type GithubPush struct {
	tango.Ctx
}

type PushInfo struct {
}

func downGithubBook(dstDir string, book string) error {
	if !strings.HasPrefix(book, "github.com/") {
		book = "github.com/" + book
	}
	url := fmt.Sprintf("https://codeload.%s/zip/master", book)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	os.RemoveAll(dstDir)
	os.MkdirAll(dstDir, os.ModePerm)
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bf := bytes.NewReader(bs)
	r, err := zip.NewReader(bf, int64(len(bs)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		p := filepath.Join(strings.Split(f.Name, "/")[1:]...)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filepath.Join(dstDir, p), os.ModePerm)
			continue
		}
		d, err := os.Create(filepath.Join(dstDir, p))
		if err != nil {
			return err
		}
		defer d.Close()

		rd, err := f.Open()
		if err != nil {
			return err
		}

		if _, err = io.Copy(d, rd); err != nil {
			return err
		}
	}

	return nil
}

func (g *GithubPush) Get() {
	g.Post()
}

func (g *GithubPush) makebook(book string) {
	dstDir := filepath.Join(conf.ReposRootPath, book)
	err := downGithubBook(dstDir, book)
	if err != nil {
		g.Error("GithubPush:", err)
		return
	}

	_, err = gobook.MakeBook(filepath.Join(conf.BooksRootPath, book), dstDir)
	if err != nil {
		g.Error("GithubPush:", err)
		return
	}
}

func (g *GithubPush) Post() {
	book := g.Req().FormValue("book")
	//userName := g.Req().FormValue("username")
	/*token := g.Req().FormValue("token")
	bs, err := ioutil.ReadAll(g.Req().Body)
	if err != nil {
		g.Error("GithubPush:", err)
		return
	}
	defer g.Req().Body.Close()*/

	/*var pushInfo GithubPush
	err = json.UnMarshal(bs, &pushInfo)
	if er != nil {
		g.Error("GithubPush:", err)
		return
	}*/

	go g.makebook(book)
}
