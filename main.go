package main

import "github.com/lunny/tango"

var (
	ReposRootPath string
	BooksRootPath string
)

func main() {
	ReposRootPath = "./repos"
	BooksRootPath = "./books"

	t := tango.Classic()
	t.Use(tango.Static(tango.StaticOptions{
		RootPath: BooksRootPath,
		Prefix:   "read",
	}))
	t.Any("/github", new(GithubPush))
	t.Run()
}
