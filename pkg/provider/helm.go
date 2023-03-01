package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
)

type Helm struct {
}

func (p Helm) Fetch(tag string, artifact manifestv1alpha1.Artifact) ([]byte, error) {
	return nil, nil
}
