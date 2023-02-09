package download

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"reflect"
	"strings"
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
	v.SetEnvPrefix("AM")
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
	return nil
}

func (o *Options) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *Options) Run(cmd *cobra.Command, args []string) error {
	if err := cmd.ParseFlags(args); err != nil {
		klog.Fatalln(err)
	}
	fmt.Println(o)
	return nil
}
