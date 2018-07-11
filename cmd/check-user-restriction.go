package main

import (
	"io"
	"github.com/spf13/cobra"
	"github.com/jlaffaye/ftp"
	"github.com/sgnl04/check-ftp/pkg/checks"
	"github.com/sgnl04/check-ftp/pkg/utils"
	"github.com/benkeil/icinga-checks-library"
)

type (
	checkUserRestriction struct {
		out           io.Writer
		FtpServerConn ftp.ServerConn
		Host          string
		Port          int
		User          string
		Password      string
		File          string
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
		Run: func(cmd *cobra.Command, args []string) {
			c.run()
		},
	}

	cmd.Flags().StringVarP(&c.Host, "host", "H", "", "the ftp host")
	cmd.Flags().IntVarP(&c.Port, "port", "P", 21, "the ftp port")
	cmd.Flags().StringVarP(&c.User, "user", "u", "", "the ftp user")
	cmd.Flags().StringVarP(&c.Password, "password", "p", "", "the ftp password")
	cmd.Flags().StringVarP(&c.File, "file-exists", "f", "", "a file that will be checked for existence")
	cmd.Flags().CountVarP(&c.Verbose, "verbose", "v", "enable verbose output")

	return cmd
}

func (c *checkUserRestriction) run() {
	client, err := utils.NewServerConn(c.Host, c.Port)
	if err != nil {
		icinga.NewResult("checkUserRestriction.run", icinga.ServiceStatusUnknown, err.Error()).Exit()
	}
	c.FtpServerConn = *client

	userRestriction := checks.NewUserRestriction(c.FtpServerConn)
	results := userRestriction.UserRestriction(checks.UserRestrictionOptions{
		Host:     c.Host,
		Port:     c.Port,
		User:     c.User,
		Password: c.Password,
		File:     c.File,
		Verbose:  c.Verbose,
	})
	results.Exit()
}
