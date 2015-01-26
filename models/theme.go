package models

import "time"

type Theme struct {
	Id      int64
	Name    string
	Created time.Time `xorm:"created"`
}
