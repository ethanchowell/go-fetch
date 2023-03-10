package provider

import (
	"github.com/xanzy/go-gitlab"
)

type GitLab struct {
}

func (p GitLab) Fetch(tag string, artifact string) error {
	gitlab.NewOAuthClient("", gitlab.WithBaseURL(""))
	return nil
}
