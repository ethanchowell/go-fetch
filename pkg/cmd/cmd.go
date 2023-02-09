package cmd

import (
	"github.com/ethanchowell/go-fetch/pkg/cmd/download"
	"github.com/ethanchowell/go-fetch/pkg/cmd/version"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

const (
	long = `
Go-Fetch is a CLI tool to describe and fetch external artifacts through a yaml descriptor
`
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-fetch",
		Short: "A tool for fetching external artifacts",
		Long:  long,

		Run: runHelp,
	}

	downloadCmd := download.NewCmd()
	versionCmd := version.NewCmd()

	cmd.AddCommand(
		downloadCmd,
		versionCmd,
	)

	return cmd
}

func runHelp(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		klog.Fatalln(err)
	}
}
