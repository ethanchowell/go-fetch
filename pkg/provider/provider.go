package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/artifact-manager/pkg/apis/manifest/v1alpa1"
)

type Provider interface {
	Fetch(manifestv1alpha1.Release) ([]byte, error)
}
