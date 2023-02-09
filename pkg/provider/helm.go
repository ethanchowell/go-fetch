package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/artifact-manager/pkg/apis/manifest/v1alpa1"
)

type Helm struct {
}

func (p *Helm) Fetch(release manifestv1alpha1.Release) ([]byte, error) {
	return nil, nil
}
