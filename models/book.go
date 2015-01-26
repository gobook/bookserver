package models

import "time"

type Book struct {
	Id       int64
	Name     string
	AuthorId int64 `xorm:"index"`
	RepoPath string
	Cover    string
	Theme    string
	Token    string
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

func AddBook(book *Book) error {
	_, err := orm.Insert(book)
	return err
}

func RecentBooks() ([]*Book, error) {
	var books = make([]*Book, 0)
	err := orm.Desc("id").Find(&books)
	return books, err
}

func LastUpdatedBooks() ([]*Book, error) {
	var books = make([]*Book, 0)
	err := orm.Desc("updated").Find(&books)
	return books, err
}
