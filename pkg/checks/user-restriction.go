package checks

import (
	"fmt"
	"github.com/benkeil/icinga-checks-library"
	"github.com/jlaffaye/ftp"
	"strings"
	"log"
)

type (
	// UserRestrictionCheck interface to check a user's restrictions
	UserRestrictionCheck interface {
		UserRestriction(UserRestrictionOptions) icinga.Result
	}

	userRestrictionImpl struct {
		ServerConn ftp.ServerConn
	}
)

// NewUserRestriction creates a new instance of UserRestrictionCheck
func NewUserRestriction(serverConn ftp.ServerConn) UserRestrictionCheck {
	return &userRestrictionImpl{ServerConn: serverConn}
}

// CheckAvailableAddressesOptions contains options needed to run CheckAvailableAddresses check
type UserRestrictionOptions struct {
	Host     string
	Port     int
	User     string
	Password string
	File     string
	Verbose  int
}

// CheckAvailableAddresses checks if the deployment has a minimum of available replicas
func (c *userRestrictionImpl) UserRestriction(options UserRestrictionOptions) icinga.Result {
	context := "User.Restriction"

	err := c.ServerConn.Login(options.User, options.Password)
	if err != nil {
		return icinga.NewResult(context, icinga.ServiceStatusCritical, fmt.Sprintf(
			"can't login to host [%s] with user [%s]: %s", options.Host, options.User, err))
	}

	home_dir := "/"
	result := staysInDir(c.ServerConn, home_dir, options.User, context)
	if result != nil {
		return result
	}

	err = c.ServerConn.ChangeDirToParent()
	if err != nil {
		return icinga.NewResult(context, icinga.ServiceStatusCritical, fmt.Sprintf(
			"change dir to parent failed: %s", err))
	}

	result = staysInDir(c.ServerConn, home_dir, options.User, context)
	if result != nil {
		return result
	}

	if len(options.File) > 0 {
		if options.Verbose > 0 {
			log.Printf("Searching for [%s] in [%s]", options.File, home_dir)
		}

		entries, err := c.ServerConn.List(home_dir)
		if err != nil {
			return icinga.NewResult(context, icinga.ServiceStatusUnknown, fmt.Sprintf(
				"failed to list [%v]: %s", options.File, err))
		}

		fileExists := false
		for _, entry := range entries {
			if options.Verbose > 0 {
				log.Printf("Searching [%s]", entry.Name)
			}
			if 0 == strings.Compare(entry.Name, options.File) && entry.Type == ftp.EntryTypeFile {
				fileExists = true
				break
			}
		}

		if !fileExists {
			return icinga.NewResult(context, icinga.ServiceStatusCritical, fmt.Sprintf(
				"check file [%s] does not exist", options.File))
		}
	}

	err = c.ServerConn.Logout()
	if err == nil {
		return icinga.NewResult(context, icinga.ServiceStatusUnknown, fmt.Sprintf(
			"failed to log out of [%s]", options.Host))
	}

	return icinga.NewResult(context, icinga.ServiceStatusOk, fmt.Sprintf("user [%s] has expected restrictions", options.User))
}

func staysInDir(conn ftp.ServerConn, dir string, user string, context string) icinga.Result {
	currentDir, err := conn.CurrentDir()
	if err != nil {
		return icinga.NewResult(context, icinga.ServiceStatusUnknown, fmt.Sprintf("cannot read current dir: %s", err))
	}

	if 0 != strings.Compare(dir, currentDir) {
		return icinga.NewResult(context, icinga.ServiceStatusCritical, fmt.Sprintf(
			"user [%s] is expected to be in [%s] but is in [%s]", user, dir, currentDir))
	}

	return nil
}
