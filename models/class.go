package models

import "time"

type Class struct {
	Id       int64
	Name     string
	ParentId int64     `xorm:"index"`
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

func FindFirstClasses() (classes []Class, err error) {
	err = orm.Where("parent_id = ?", 0).Find(&classes)
	return
}

type BookClass struct {
	BookId  int64 `xorm:"index"`
	ClassId int64 `xorm:"index"`
}

var (
	defaultClasses = []Class{
		{Name:"计算机"},
		{Name:"工业技术"},
		{Name:"数理化"},
		{Name:"经济管理"},
		{Name:"人文社科"},
	}
)

func initClasses() error {
	total, err := orm.Count(new(Class))
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}
	_, err = orm.Insert(&defaultClasses)
	return err
}