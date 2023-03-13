package provider

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
	"path"
	"strings"
)

type Docker struct {
	repo   string
	client *client.Client

	Store
}

func (p Docker) Fetch(tag string, artifact string) error {
	var saveFile, imageTag string
	targetDir := path.Join(p.rootDir, p.repo, tag)
	if tag == "" && strings.Count(artifact, ":") == 1 {
		s := strings.Split(artifact, ":")
		saveFile = fmt.Sprintf("%s-%s.tar.gz", s[0], s[1])
		imageTag = path.Join(p.repo, artifact)
	} else {
		saveFile = fmt.Sprintf("%s.tar.gz", artifact)
		imageTag = path.Join(p.repo, strings.Join([]string{artifact, tag}, ":"))
	}

	if _, err := os.Stat(path.Join(targetDir, saveFile)); !os.IsNotExist(err) && err == nil {
		fmt.Printf("skipping download for %s\n", path.Join(targetDir, saveFile))
		return nil
	}

	responsePull, err := p.client.ImagePull(context.Background(), imageTag, types.ImagePullOptions{})
	if err != nil {
		fmt.Println("Error pulling image")
		fmt.Println(err)
		return err
	}

	defer responsePull.Close()
	_, err = io.ReadAll(responsePull)
	if err != nil {
		return err
	}

	r, err := p.client.ImageSave(context.Background(), []string{imageTag})
	if err != nil {
		fmt.Println("Error saving image")
		fmt.Println(err)
		return err
	}

	defer r.Close()

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Println("Error compressing image")
		fmt.Println(err)
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		fmt.Println("Error writing image")
		fmt.Println(err)
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write(buf.Bytes())
	assetSum := hex.EncodeToString(hash.Sum(nil))
	p.checksumData.Write([]byte(fmt.Sprintf("%s %s\n", assetSum, path.Join(targetDir, saveFile))))

	return saveAsset(targetDir, saveFile, buf.Bytes())
}
