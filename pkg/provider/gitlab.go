package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
)

type GitLab struct {
}

func (p *GitLab) Fetch(release manifestv1alpha1.Release) ([]byte, error) {
	return nil, nil
}
