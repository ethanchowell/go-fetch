package provider

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/go-github/v50/github"
	"github.com/minio/sha256-simd"
	"io"
	"net/http"
	"os"
	"path"
)

type GitHub struct {
	group string
	repo  string

	Store
}

func (p GitHub) Fetch(tag string, artifact string) error {
	c := github.NewClient(nil)
	fmt.Printf("fetching artifact: %s\n", artifact)

	targetDir := path.Join(p.rootDir, p.group, p.repo, tag)
	if _, err := os.Stat(path.Join(targetDir, artifact)); !os.IsNotExist(err) && err == nil {
		fmt.Printf("skipping download for %s\n", path.Join(targetDir, artifact))
		return nil
	}

	releaseService, _, err := c.Repositories.GetReleaseByTag(context.Background(), p.group, p.repo, tag)
	if err != nil {
		return err
	}

	for _, asset := range releaseService.Assets {
		if *asset.Name == artifact {
			data, err := downloadAsset(c, p.group, p.repo, *asset.ID)
			if err != nil {
				return err
			}

			hash := sha256.New()
			hash.Write(data)
			assetSum := hex.EncodeToString(hash.Sum(nil))
			p.checksumData.Write([]byte(fmt.Sprintf("%s %s\n", assetSum, path.Join(targetDir, artifact))))

			return saveAsset(targetDir, artifact, data)
		}
	}
	return nil
}

func downloadAsset(client *github.Client, group, repo string, id int64) ([]byte, error) {
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

	return data, err
}

func saveAsset(targetDir, name string, data []byte) error {
	filePath := path.Join(targetDir, name)
	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}
