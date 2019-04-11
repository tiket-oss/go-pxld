package pxld

import (
	"bytes"
	"fmt"
	"io"
)

const (
	ProxySQLQuery = 0
)

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
