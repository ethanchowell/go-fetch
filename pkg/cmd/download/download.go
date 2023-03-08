package download

import (
	"fmt"
	manifestv1alpha1 "github.com/ethanchowell/go-fetch/pkg/apis/manifest/v1alpa1"
	"github.com/ethanchowell/go-fetch/pkg/provider"
	"github.com/ethanchowell/go-fetch/pkg/util"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io"
	"k8s.io/klog/v2"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
)

const (
	long = `
Parse a given YAML manifest for artifacts that should be downloaded.
`
)

type Options struct {
	File string `flag:"manifest"`

	GitLabToken string `flag:"gitlab-token" yaml:"gitlabToken"`
	GitHubToken string `flag:"github-token" yaml:"githubToken"`
	ArtToken    string `flag:"artifactory-token" yaml:"artToken"`
}

func NewCmd() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download artifacts defined in a manifest.",
		Long:  long,

		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Check(cmd, args); err != nil {
				klog.Errorln(err)
			}
			if err := o.Validate(cmd, args); err != nil {
				klog.Errorln(err)
			}
			if err := o.Run(cmd, args); err != nil {
				klog.Errorln(err)
			}
		},
	}

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("GO_FETCH")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	cmd.Flags().StringVar(&o.File, "manifest", "./artifacts.yaml", "Path to the manifest containing artifacts to download. Can be set from AM_MANIFEST.")
	cmd.Flags().StringVar(&o.GitLabToken, "gitlab-token", "", "The API token for authenticating with GitLab. Can be set from AM_GITLAB_TOKEN.")
	cmd.Flags().StringVar(&o.GitHubToken, "github-token", "", "The API token for authenticating with GitHub. Can be set from AM_GITHUB_TOKEN.")
	cmd.Flags().StringVar(&o.ArtToken, "artifactory-token", "", "The API token for authenticating with Artifactory. Can be set from AM_ARTIFACTORY_TOKEN.")

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		klog.Fatalln(err)
	}

	if err := registerFlags(v, "", cmd.Flags(), o); err != nil {
		klog.Fatalln(err)
	}

	if err := v.Unmarshal(o, decodeFromFlagTag); err != nil {
		klog.Fatalln(err)
	}

	return cmd
}

func registerFlags(v *viper.Viper, prefix string, flagSet *pflag.FlagSet, options interface{}) error {
	val := reflect.ValueOf(options)
	var typ reflect.Type
	if val.Kind() == reflect.Ptr {
		typ = val.Elem().Type()
	} else {
		typ = val.Type()
	}

	for i := 0; i < typ.NumField(); i++ {
		// pull out the struct tags:
		//    flag - the name of the command line flag
		//    cfg - the name of the config file option
		field := typ.Field(i)
		fieldV := reflect.Indirect(val).Field(i)
		fieldName := strings.Join([]string{prefix, field.Name}, ".")

		if isUnexported(field.Name) {
			// Unexported fields cannot be set by a user, so won't have tags or flags, skip them
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			err := registerFlags(v, fieldName, flagSet, fieldV.Interface())
			if err != nil {
				return err
			}
			continue
		}

		flagName := field.Tag.Get("flag")
		if flagName == "" {
			return fmt.Errorf("field %q does not have required tags (flag)", fieldName)
		}

		if flagSet == nil {
			return fmt.Errorf("flagset cannot be nil")
		}

		f := flagSet.Lookup(flagName)
		if f == nil {
			return fmt.Errorf("field %q does not have a registered flag", flagName)
		}
		if err := v.BindEnv(flagName); err != nil {
			return fmt.Errorf("error binding flag for field %q: %w", fieldName, err)
		}
	}

	return nil
}

// decodeFromCfgTag sets the Viper decoder to read the names from the `cfg` tag
// on each struct entry.
func decodeFromFlagTag(c *mapstructure.DecoderConfig) {
	c.TagName = "flag"
}

// if it is unexported.
func isUnexported(name string) bool {
	if len(name) == 0 {
		// This should never happen
		panic("field name has len 0")
	}

	first := string(name[0])
	return first == strings.ToLower(first)
}

func (o *Options) Check(cmd *cobra.Command, args []string) error {
	_, err := os.Stat(o.File)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", o.File)
	}

	if err != nil {
		return fmt.Errorf("could not check file: %w", err)
	}
	return nil
}

func (o *Options) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *Options) Run(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(o.File)
	if err != nil {
		return err
	}

	m := manifestv1alpha1.Manifest{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, release := range m.Releases {
		wg.Add(1)
		go func(release manifestv1alpha1.Release) {
			defer wg.Done()

			p := provider.New(release.Repo)
			downloadArtifacts(m.Target, p, release)
		}(release)
	}

	wg.Wait()
	return nil
}

func downloadArtifacts(rootDir string, p provider.Provider, release manifestv1alpha1.Release) {
	var wg sync.WaitGroup
	for _, artifact := range release.Artifacts {
		wg.Add(1)
		go func(artifact manifestv1alpha1.Artifact) {
			defer wg.Done()
			targetDir := path.Join(rootDir, release.Repo.Name, release.Tag)
			err := checkFile(path.Join(targetDir, artifact.Name), artifact.Checksum)
			if err != nil {
				fmt.Printf("skipping download for %s: %s\n", path.Join(targetDir, artifact.Name), err)
				return
			}

			fmt.Printf("downloading artifact: %s\n", artifact.Name)
			data, err := p.Fetch(release.Tag, artifact)
			if err != nil {
				fmt.Printf("could not download artifact: %s\n", artifact.Name)
				return
			}

			err = saveArtifact(targetDir, artifact.Name, artifact.Checksum, data)
			if err != nil {
				fmt.Printf("could not save artifact to %s: %s\n", path.Join(targetDir, artifact.Name), err)
			}
		}(artifact)
	}
	wg.Wait()
}

func saveArtifact(targetDir, name, checksum string, data []byte) error {
	filePath := path.Join(targetDir, name)
	_, err := os.Stat(filePath)
	if err := checkFile(filePath, checksum); err != nil {
		return err
	} else {
		if err := os.MkdirAll(targetDir, 0700); err != nil {
			return err
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func checkFile(filePath, checksum string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	ok, err := util.ValidateChecksum(checksum, data)
	if err != nil {
		return err
	}

	if ok {
		return fmt.Errorf("file already exists with valid checksum: %s", filePath)
	}

	return nil
}
