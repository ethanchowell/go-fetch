package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var VERSION = "undefined"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of the application",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("go-fetch %s (built with %s)\n", VERSION, runtime.Version())
		},
	}

	return cmd
}
