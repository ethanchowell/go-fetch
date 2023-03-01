package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
)

type Artifactory struct {
}

func (p Artifactory) Fetch(tag string, artifact manifestv1alpha1.Artifact) ([]byte, error) {
	auth.NewArtifactoryDetails()
	return nil, nil
}
