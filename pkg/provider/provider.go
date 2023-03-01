package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
	"strings"
)

type Provider interface {
	Fetch(tag string, artifact manifestv1alpha1.Artifact) ([]byte, error)
}

func New(repo manifestv1alpha1.Repo) Provider {
	switch repo.Provider {
	case manifestv1alpha1.Artifactory:
		return Artifactory{}
	case manifestv1alpha1.Docker:
		return Docker{}
	case manifestv1alpha1.Generic:
		return Generic{}
	case manifestv1alpha1.GitHub:
		s := strings.Split(repo.Name, "/")
		return GitHub{
			Group: s[0],
			Repo:  s[1],
		}
	case manifestv1alpha1.GitLab:
		return GitLab{}
	case manifestv1alpha1.Helm:
		return Helm{}
	default:
		return Generic{}
	}
}
