package main

import (
	"io"
	"github.com/spf13/cobra"
	"github.com/jlaffaye/ftp"
	"github.com/sgnl04/check-ftp/pkg/utils"
	"github.com/benkeil/icinga-checks-library"
	"github.com/sgnl04/check-ftp/pkg/checks"
)

type (
	checkUserRestriction struct {
		out           io.Writer
		FtpServerConn ftp.ServerConn
		Host          string
		User          string
		Password      string
		Verbose       int
	}
)

func newCheckUserRestriction(out io.Writer) *cobra.Command {
	c := &checkUserRestriction{out: out}

	cmd := &cobra.Command{
		Use:          "user-restriction",
		Short:        "check that a user can login, is initially located in his home dir and cannot leave it",
		SilenceUsage: false,
		Args:         NameArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			c.Host = args[0]
			client, err := utils.NewServerConn(c.Host)
			if err != nil {
				icinga.NewResult("NewServerConn", icinga.ServiceStatusUnknown, err.Error()).Exit()
			}
			c.FtpServerConn = *client
		},
		Run: func(cmd *cobra.Command, args []string) {
			c.run()
		},
	}

	cmd.Flags().StringVarP(&c.User, "user", "u", "*", "the ftp user")
	cmd.Flags().StringVarP(&c.Password, "password", "p", "*", "the ftp password")
	cmd.Flags().CountVarP(&c.Verbose, "verbose", "v", "enable verbose output")

	return cmd
}

func (c *checkUserRestriction) run() {
	userRestriction := checks.NewUserRestriction(c.FtpServerConn, c.Host)
	results := userRestriction.UserRestriction(checks.UserRestrictionOptions{
		Host:              c.Host,
		User:              c.User,
		Password:          c.Password,
		Verbose:           c.Verbose,
	})
	results.Exit()
}
