package main

import (
	"errors"
	"os"
	"github.com/spf13/cobra"
)

var (
	globalUsage = `This program makes ftp checks`
	version     string
)

func main() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "check-ftp",
		Short:        "check-ftp checks ftp",
		Long:         globalUsage,
		Version:      version,
		SilenceUsage: true,
	}

	cmd.PersistentFlags().Parse(args)

	out := cmd.OutOrStdout()

	cmd.AddCommand(
		// check commands
		newCheckUserRestriction(out),
	)

	return cmd
}

// NameArgs returns an error if there are not exactly 1 arg containing the resource name.
func NameArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("only one check command is allowed")
		}
		return nil
	}
}
