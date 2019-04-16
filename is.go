package pxld

import (
	"bytes"
	"fmt"
	"io"
)

const (
	// ProxySQLQuery byte should be 0
	ProxySQLQuery = 0
)

// IsProxySQLQuery check if the data is a valid proxysql event
func IsProxySQLQuery(dataStream io.Reader) (err error) {
	data := make([]byte, 1)

	_, err = dataStream.Read(data)
	if err != nil {
		return
	}

	if !bytes.Equal([]byte{ProxySQLQuery}, data) {
		err = fmt.Errorf("not a valid proxy sql query log line")
		return
	}

	return
}
