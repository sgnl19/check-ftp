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
	Host              string
	Port              int
	User              string
	Password          string
	Verbose           int
}

// CheckAvailableAddresses checks if the deployment has a minimum of available replicas
func (c *userRestrictionImpl) UserRestriction(options UserRestrictionOptions) icinga.Result {
	name := "User.Restriction"

	err := c.ServerConn.Login(options.User, options.Password)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusCritical, fmt.Sprintf(
			"can't login to host [%s] with user [%s]: %s", options.Host, options.User, err))
	}

	home_dir := "/home/" + options.User
	dir, err := c.ServerConn.CurrentDir()
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf(
			"cannot read current dir: %s", err))
	}

	if 0 != strings.Compare(dir, home_dir) {
		return icinga.NewResult(name, icinga.ServiceStatusCritical, fmt.Sprintf(
			"user [%v] is expected to be in [%s] but is in [%s]", options.User, home_dir, dir))
	}

	err = c.ServerConn.ChangeDirToParent()
	if err == nil {
		return icinga.NewResult(name, icinga.ServiceStatusCritical, fmt.Sprintf(
			"user [%s] is able to navigate out of [%s]", options.User, home_dir))
	}

	if options.Verbose > 0 {
		log.Printf("user [%s] is not able to navogate out of [%s]: %s", options.User, home_dir, err)
	}

	err = c.ServerConn.Logout()
	if err == nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf(
			"failed to log out of [%s]", options.Host))
	}

	return icinga.NewResult(name, icinga.ServiceStatusOk, fmt.Sprintf("user [%s] has expected restrictions", options.User))
}
