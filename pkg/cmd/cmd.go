package cmd

import (
	"github.com/ethanchowell/artifact-manager/pkg/cmd/download"
	"github.com/ethanchowell/artifact-manager/pkg/cmd/version"
	"github.com/spf13/cobra"
	"log"
)

const (
	long = `
Artifact Manager is a CLI tool to describe and fetch external artifacts through a yaml descriptor
`
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact-manager",
		Short: "A tool for fetching external artifacts",
		Long:  long,

		Run: runHelp,
	}

	//helpCmd := help.NewCmd()
	downloadCmd := download.NewCmd()
	versionCmd := version.NewCmd()

	cmd.AddCommand(
		//helpCmd,
		downloadCmd,
		versionCmd,
	)

	return cmd
}

func runHelp(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); err != nil {
		log.Fatalln(err)
	}
}
