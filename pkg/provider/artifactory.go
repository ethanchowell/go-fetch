package provider

import (
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
)

type Artifactory struct {
	Store
}

func (p Artifactory) Fetch(tag string, artifact string) error {
	auth.NewArtifactoryDetails()
	return nil
}
