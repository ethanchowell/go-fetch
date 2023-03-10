package provider

import (
	"context"
	docker "github.com/docker/docker/client"
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
	"io"
	"strings"
)

type Provider interface {
	Fetch(tag string, artifact string) error
}

type Store struct {
	rootDir      string
	checksumData io.Writer
}

func NewStore(rootDir string, w io.Writer) Store {
	return Store{
		rootDir:      rootDir,
		checksumData: w,
	}
}

func New(repo manifestv1alpha1.Repo, store Store) Provider {
	switch repo.Provider {
	case manifestv1alpha1.Artifactory:
		return Artifactory{
			Store: store,
		}
	case manifestv1alpha1.Docker:
		c, _ := docker.NewClientWithOpts(docker.FromEnv)
		c.NegotiateAPIVersion(context.Background())
		return Docker{
			repo:   repo.Name,
			client: c,
			Store:  store,
		}
	case manifestv1alpha1.Generic:
		return Generic{
			Store: store,
		}
	case manifestv1alpha1.GitHub:
		s := strings.Split(repo.Name, "/")
		return GitHub{
			group: s[0],
			repo:  s[1],
			Store: store,
		}
	case manifestv1alpha1.GitLab:
		return GitLab{}
	case manifestv1alpha1.Helm:
		return Helm{
			repo:  repo.Name,
			Store: store,
		}
	default:
		return Generic{
			Store: store,
		}
	}
}
