package models

import "time"

type Theme struct {
	Id      int64
	Name    string    `xorm:"unique"`
	Created time.Time `xorm:"created"`
}

func initThemes() error {
	cnt, err := orm.Count(new(Theme))
	if err != nil {
		return err
	}

	if cnt > 0 {
		return nil
	}

	_, err = orm.Insert(&Theme{
		Name: "gitbook",
	})
	return err
}
