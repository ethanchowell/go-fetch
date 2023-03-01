package provider

import (
	"github.com/docker/docker/client"
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
)

type Docker struct {
	client *client.Client
}

func (p Docker) Fetch(tag string, artifact manifestv1alpha1.Artifact) ([]byte, error) {
	client.NewClientWithOpts(client.FromEnv)
	return nil, nil
}
