package provider

import (
	"context"
	"errors"
	"fmt"
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
	"github.com/ethanchowell/go-fetch/pkg/util"
	"github.com/google/go-github/v50/github"
	"io"
	"net/http"
)

type GitHub struct {
	Group string
	Repo  string
}

func (p GitHub) Fetch(tag string, artifact manifestv1alpha1.Artifact) ([]byte, error) {
	c := github.NewClient(nil)
	fmt.Printf("fetching artifact: %s\n", artifact.Name)

	releaseService, _, err := c.Repositories.GetReleaseByTag(context.Background(), p.Group, p.Repo, tag)
	if err != nil {
		return nil, err
	}

	for _, asset := range releaseService.Assets {
		if *asset.Name == artifact.Name {
			return downloadAsset(c, p.Group, p.Repo, *asset.ID, artifact.Checksum)
		}
	}
	return nil, nil
}

func downloadAsset(client *github.Client, group, repo string, id int64, checksum string) ([]byte, error) {
	rc, redirect, err := client.Repositories.DownloadReleaseAsset(context.Background(), group, repo, id, http.DefaultClient)
	defer rc.Close()
	if err != nil {
		return nil, err
	}

	if redirect != "" {
		fmt.Printf("Github wants to redirect: %s", redirect)
		return nil, nil
	}

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	ok, err := util.ValidateChecksum(checksum, data)
	if err != nil {
		return nil, fmt.Errorf("failed validating checksum: %w", err)
	}

	if !ok {
		return nil, errors.New("unable to validate checksum")
	}

	return data, err
}
