package routers

import (
	"github.com/lunny/tango"
	"github.com/google/go-github"
)

type Github struct {
	tango.Compress
}

func (Github) Get() {
	
	"https://github.com/login/oauth/authorize"

	t := &oauth.Transport{
	  Token: &oauth.Token{AccessToken: ""},
	}

	client := github.NewClient(t.Client())

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List("", nil)
}

type GithubCallback struct {
	tango.Compress
}

func (GithubCallback) Get() {

}