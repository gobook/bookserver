package models

type Class struct {
	Id   int64
	Name string
}

type BookClass struct {
	BookId  int64 `xorm:"index"`
	ClassId int64 `xorm:"index"`
}
