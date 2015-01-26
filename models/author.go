package models

import "time"

type Author struct {
	Id      int64
	Name    string
	Email   string
	Created time.Time `xorm:"created"`
}
