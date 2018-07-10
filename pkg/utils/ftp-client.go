package utils

import (
	"fmt"
	"github.com/jlaffaye/ftp"
)

// NewServerConn creates a new instance of ServerConn
func NewServerConn(host string) (*ftp.ServerConn, error) {

	if len(host) == 0 {
		return nil, fmt.Errorf("ftp host missing: %s", host)
	}
	conn, err := ftp.Connect(host)

	if err != nil {
		return nil, fmt.Errorf("could not connect to host [%v]: %s", host, err)
	}

	return conn, err
}
