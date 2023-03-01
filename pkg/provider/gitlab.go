package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
	"github.com/xanzy/go-gitlab"
)

type GitLab struct {
}

func (p GitLab) Fetch(tag string, artifact manifestv1alpha1.Artifact) ([]byte, error) {
	gitlab.NewOAuthClient("", gitlab.WithBaseURL(""))
	return nil, nil
}
