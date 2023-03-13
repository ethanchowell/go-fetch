package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"io"
	"net/url"
	"path"
	"strings"

	"log"
	"os"
)

type Helm struct {
	repo string

	Store
}

func (p Helm) Fetch(tag string, artifact string) error {
	chartRef := artifact

	repoDir := p.repo
	if isURL(p.repo) {
		repoDir = "helm"
	}

	var saveFile string
	targetDir := path.Join(p.rootDir, repoDir, tag)
	if tag != "" {
		saveFile = fmt.Sprintf("%s-%s.tgz", artifact, tag)
	}
	if tag == "" && strings.Count(artifact, ":") == 1 {
		s := strings.Split(artifact, ":")
		artifact = s[0]
		saveFile = fmt.Sprintf("%s-%s.tgz", s[0], s[1])
	}

	_, err := os.Stat(path.Join(targetDir, saveFile))
	if !os.IsNotExist(err) && err == nil {
		fmt.Printf("skipping download for %s\n", path.Join(targetDir, artifact))
		return nil
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0700); err != nil {
			return err
		}
	}

	settings := helm.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return err
	}

	pull := action.NewPullWithOpts(action.WithConfig(actionConfig))

	pull.Version = tag
	pull.Settings = settings

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

	if isURL(p.repo) {
		chartUrl, err := repo.FindChartInAuthAndTLSAndPassRepoURL(p.repo, pull.Username, pull.Password, chartRef, pull.Version, pull.CertFile, pull.KeyFile, pull.CaFile, pull.InsecureSkipTLSverify, pull.PassCredentialsAll, getter.All(pull.Settings))
		if err != nil {
			return err
		}
		chartRef = chartUrl
	} else {
		chartRef = fmt.Sprintf("%s/%s", p.repo, artifact)
	}

	filePath, _, err := c.DownloadTo(chartRef, tag, targetDir)
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write(data)
	assetSum := hex.EncodeToString(hash.Sum(nil))
	_, err = p.checksumData.Write([]byte(fmt.Sprintf("%s %s\n", assetSum, filePath)))
	return err
}

func isURL(s string) bool {
	_, err := url.Parse(s)
	return err == nil && strings.Count(s, "/") != 0 && strings.Count(s, ":") != 0
}
