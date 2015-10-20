package models

import (
	"time"
)

type User struct {
	Id         int64
	Name       string `xorm:"unique"`
	Email      string `xorm:"unique"`
	Passwd     string
	GithubName string
	Created    time.Time `xorm:"created"`
	Updated    time.Time `xorm:"updated"`
}

func GetUserByGithubName(githubName string) (*User, error) {
	var user User
	has, err := orm.Where("github_name = ?", githubName).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}

	return &user, nil
}

func CheckPasswd(name, passwd string) error {
	return nil
}
