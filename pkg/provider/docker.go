package provider

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
	"io"
	"os"
)

type Docker struct {
	repo   string
	client *client.Client
}

func (p Docker) Fetch(tag string, artifact manifestv1alpha1.Artifact) ([]byte, error) {
	imageTag := fmt.Sprintf("%s/%s:%s", p.repo, artifact.Name, tag)
	fmt.Println(imageTag)
	responsePull, err := p.client.ImagePull(context.Background(), imageTag, types.ImagePullOptions{})
	if err != nil {
		fmt.Println("Error pulling image")
		fmt.Println(err)
		return nil, err
	}

	defer responsePull.Close()
	io.Copy(os.Stdout, responsePull)

	r, err := p.client.ImageSave(context.Background(), []string{imageTag})
	if err != nil {
		fmt.Println("Error saving image")
		fmt.Println(err)
		return nil, err
	}

	defer r.Close()

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Println("Error compressing image")
		fmt.Println(err)
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		fmt.Println("Error writing image")
		fmt.Println(err)
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
