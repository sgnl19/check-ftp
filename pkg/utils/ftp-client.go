package utils

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"strconv"
)

// NewServerConn creates a new instance of ServerConn
func NewServerConn(host string, port int) (*ftp.ServerConn, error) {

	if len(host) == 0 {
		return nil, fmt.Errorf("ftp host is missing: %s", host)
	}
	if port <= 0 {
		return nil, fmt.Errorf("ftp port is invalid: %v", port)
	}

	address := host + ":" + strconv.Itoa(port)
	conn, err := ftp.Connect(address)

	if err != nil {
		return nil, fmt.Errorf("could not connect to host [%v]: %s", address, err)
	}

	return conn, err
}
