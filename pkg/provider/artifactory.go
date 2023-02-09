package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
)

type Artifactory struct {
}

func (p *Artifactory) Fetch(release manifestv1alpha1.Release) ([]byte, error) {
	return nil, nil
}
