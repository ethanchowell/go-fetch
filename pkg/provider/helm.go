package provider

import (
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"io"
	"net/url"
	"strings"

	"log"
	"os"
)

type Helm struct {
	repo string
}

func (p Helm) Fetch(tag string, artifact manifestv1alpha1.Artifact) ([]byte, error) {
	chartRef := artifact.Name
	settings := helm.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	pull := action.NewPullWithOpts(action.WithConfig(actionConfig))

	pull.Version = tag
	pull.Settings = settings
	if isURL(p.repo) {
		chartUrl, err := repo.FindChartInAuthAndTLSAndPassRepoURL(p.repo, pull.Username, pull.Password, chartRef, pull.Version, pull.CertFile, pull.KeyFile, pull.CaFile, pull.InsecureSkipTLSverify, pull.PassCredentialsAll, getter.All(pull.Settings))
		if err != nil {
			return nil, err
		}
		chartRef = chartUrl
	}

	var out strings.Builder

	c := downloader.ChartDownloader{
		Out:     &out,
		Keyring: pull.Keyring,
		Verify:  downloader.VerifyNever,
		Getters: getter.All(pull.Settings),
		Options: []getter.Option{
			getter.WithBasicAuth(pull.Username, pull.Password),
			getter.WithPassCredentialsAll(pull.PassCredentialsAll),
			getter.WithTLSClientConfig(pull.CertFile, pull.KeyFile, pull.CaFile),
			getter.WithInsecureSkipVerifyTLS(pull.InsecureSkipTLSverify),
		},
		RegistryClient:   actionConfig.RegistryClient,
		RepositoryConfig: pull.Settings.RepositoryConfig,
		RepositoryCache:  pull.Settings.RepositoryCache,
	}

	u, err := c.ResolveChartVersion(chartRef, pull.Version)
	if err != nil {
		return nil, err
	}

	g, err := c.Getters.ByScheme(u.Scheme)
	if err != nil {
		return nil, err
	}

	data, err := g.Get(u.String(), c.Options...)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(data)
}

func isURL(s string) bool {
	_, err := url.Parse(s)
	return err == nil
}
