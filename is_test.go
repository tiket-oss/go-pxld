package pxld

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsProxySQLQuery(t *testing.T) {
	data := []byte{0}
	buf := bytes.NewReader(data)

	err := IsProxySQLQuery(buf)
	require.NoError(t, err)
}

func TestIsProxySQLQueryNegative(t *testing.T) {
	data := []byte{1}
	buf := bytes.NewReader(data)

	err := IsProxySQLQuery(buf)
	require.Error(t, err)

	buf = bytes.NewReader([]byte{})
	err = IsProxySQLQuery(buf)
	require.Error(t, err)
}
