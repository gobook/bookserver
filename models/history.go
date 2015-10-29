package models

import "time"

type OperationType int

const (
	ManualPublish OperationType = iota
	AutoPublish
)

type History struct {
	Id         int64
	BookId     int64 `xorm:"index"`
	OperatorId int64 `xorm:"index"`
	OperationType
	CommitId string
	Created  time.Time `xorm:"created"`
}

func PublishBook(book *Book, commitId string) error {
	sess := orm.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.Insert(book); err != nil {
		return err
	}
	if _, err := sess.Insert(&History{
		BookId:        book.Id,
		CommitId:      commitId,
		OperatorId:    book.AuthorId,
		OperationType: ManualPublish,
	}); err != nil {
		return err
	}
	return sess.Commit()
}
