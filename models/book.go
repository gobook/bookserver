package models

import "time"

type BookStatus int

const (
	BookOnline BookStatus = iota
	BookOffline
)

type Book struct {
	Id           int64
	Name         string
	AuthorId     int64  `xorm:"index"`
	RepoPath     string `xorm:"unique"`
	Cover        string
	Theme        string
	Token        string
	AutoUpdated  bool
	Lang         string // 书籍语言
	AuthorName   string // 显示的书籍作者名称
	Status       BookStatus
	LastUpdated  time.Time // 书籍最后更新时间
	LastCommitId string    // 最后一次更新的CommitId
	Created      time.Time `xorm:"created"`
	Updated      time.Time `xorm:"updated"`
}

func AddBook(book *Book) error {
	_, err := orm.Insert(book)
	return err
}

func RecentBooks() ([]*Book, error) {
	var books = make([]*Book, 0)
	err := orm.Where("status = ?", BookOnline).Desc("id").Find(&books)
	return books, err
}

func LastUpdatedBooks() ([]*Book, error) {
	var books = make([]*Book, 0)
	err := orm.Where("status = ?", BookOnline).Desc("updated").Find(&books)
	return books, err
}

func FindBooksByUserId(userId int64) ([]Book, error) {
	var books []Book
	err := orm.Where("author_id=?", userId).Desc("updated").Find(&books)
	return books, err
}
