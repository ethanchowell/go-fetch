package v1alpa1

type Provider string

var (
	// GitLab uses the GitLab API to fetch artifacts.
	GitLab Provider = "gitlab"

	// GitHub uses the GitHub API to fetch artifacts.
	GitHub Provider = "github"

	// Artifactory uses the Artifactory API to fetch artifacts.
	Artifactory Provider = "artifactory"

	// Helm uses a helm client for fetching charts.
	Helm Provider = "helm"

	// Docker uses the docker client to pull and save images.
	Docker Provider = "docker"

	// Generic makes a simple HTTP GET request to fetch an artifact.
	Generic Provider = "generic"
)

// Repo configures the provider to use for downloading and an auth token.
type Repo struct {
	// Name is the name of the repository.
	Name string `yaml:"name" json:"name"`

	// AuthToken is the access token needed by the provider. If left empty,
	// the flag --<provider>-token and environment variable GO_FETCH_<PROVIDER>_TOKEN
	// is checked, and used if not empty.
	AuthToken string `yaml:"token,omitempty" json:"token,omitempty"`

	// Provider is the provider to use for fetching artifacts.
	Provider Provider `yaml:"provider" json:"provider"`
}

// Release is a tagged release in some remote repo that you want to fetch
type Release struct {
	// Tag is the release tag to download from a Repo.
	Tag string `yaml:"tag,omitempty" json:"tag,omitempty"`

	// Artifacts outline the specific items to fetch from a release
	// if left empty, the entire release will be downloaded.
	Artifacts []string `yaml:"artifacts,omitempty" json:"artifacts,omitempty"`

	// Repo configures the provider to use for downloading and an auth token.
	Repo Repo `yaml:"repo" json:"repo"`
}

// Manifest is a descriptor of artifacts and where to download them
type Manifest struct {
	// Target defines the directory that artifacts will be downlaoded
	Target string `yaml:"target" json:"target"`

	// Releases define the releases to be downloaded by this manifest
	Releases []Release `yaml:"releases" json:"releases"`

	// Package toggles exporting the downloaded content to a tar.gz
	Package bool `yaml:"package" json:"package"`
}
