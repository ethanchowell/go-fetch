package help

import (
	"github.com/spf13/cobra"
	"log"
	"strings"
)

const (
	long = `
Helm provides help for any supported sub command.
Just type artifact-manager help [path to command] for details.
`
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "help [command]",
		DisableFlagsInUseLine: true,
		Short:                 "Print the help for a command",
		Long:                  long,

		Run: RunHelp,
	}

	return cmd
}

// RunHelp checks given arguments and executes command
func RunHelp(cmd *cobra.Command, args []string) {
	foundCmd, _, err := cmd.Root().Find(args)

	// NOTE(andreykurilin): actually, I did not find any cases when foundCmd can be nil,
	//   but let's make this check since it is included in original code of initHelpCmd
	//   from github.com/spf13/cobra
	if foundCmd == nil {
		cmd.Printf("Unknown help topic %#q.\n", args)
		if err := cmd.Root().Usage(); err != nil {
			log.Fatalln(err)
		}
	} else if err != nil {
		// print error message at first, since it can contain suggestions
		cmd.Println(err)

		argsString := strings.Join(args, " ")
		var matchedMsgIsPrinted = false
		for _, foundCmd := range foundCmd.Commands() {
			if strings.Contains(foundCmd.Short, argsString) {
				if !matchedMsgIsPrinted {
					cmd.Printf("Matchers of string '%s' in short descriptions of commands: \n", argsString)
					matchedMsgIsPrinted = true
				}
				cmd.Printf("  %-14s %s\n", foundCmd.Name(), foundCmd.Short)
			}
		}

		if !matchedMsgIsPrinted {
			// if nothing is found, just print usage
			if err := cmd.Root().Usage(); err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		if len(args) == 0 {
			// help message for help command :)
			foundCmd = cmd
		}
		helpFunc := foundCmd.HelpFunc()
		helpFunc(foundCmd, args)
	}
}
