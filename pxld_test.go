package pxld

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

var (
	testData = []byte{
		0x5C, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x15, 0x06, 0x64,
		0x69, 0x64, 0x61, 0x73,
		0x79, 0x04, 0x74, 0x65,
		0x73, 0x74, 0x0F, 0x31,
		0x32, 0x37, 0x2E, 0x30,
		0x2E, 0x30, 0x2E, 0x31,
		0x3A, 0x33, 0x33, 0x36,
		0x38, 0x30, 0x01, 0x0E,
		0x31, 0x32, 0x37, 0x2E,
		0x30, 0x2E, 0x30, 0x2E,
		0x31, 0x3A, 0x33, 0x33,
		0x30, 0x36, 0xFE, 0x3A,
		0xF1, 0x74, 0x91, 0x28,
		0x86, 0x05, 0x00, 0xFE,
		0x3A, 0xF1, 0x74, 0x91,
		0x28, 0x86, 0x05, 0x00,
		0xFE, 0x42, 0x6F, 0x13,
		0xB3, 0x37, 0x1D, 0xDF,
		0x38, 0x12, 0x73, 0x65,
		0x6C, 0x65, 0x63, 0x74,
		0x20, 0x2A, 0x20, 0x66,
		0x72, 0x6F, 0x6D, 0x20,
		0x74, 0x65, 0x73, 0x74,
	}
	line = &LogLine{
		MessageLength: 92,
		ThreadID:      21,
		RawMessage:    testData[8:],
		Username:      "didasy",
		Schema:        "test",
		QueryDigest:   "0xB3136F4238DF1D37",
		HID:           1,
		ClientAddr:    "127.0.0.1:33680",
		ServerAddr:    "127.0.0.1:3306",
		Query:         "select * from test",
		Duration:      0,
	}
	lineJSON = `{
  "message_length": 92,
  "raw_message": "ABUGZGlkYXN5BHRlc3QPMTI3LjAuMC4xOjMzNjgwAQ4xMjcuMC4wLjE6MzMwNv468XSRKIYFAP468XSRKIYFAP5CbxOzNx3fOBJzZWxlY3QgKiBmcm9tIHRlc3Q=",
  "thread_id": 21,
  "username": "didasy",
  "schema": "test",
  "start_at": "2019-04-10T15:08:00.727354+07:00",
  "end_at": "2019-04-10T15:08:00.727354+07:00",
  "query_digest": "0xB3136F4238DF1D37",
  "hid": 1,
  "client_addr": "127.0.0.1:33680",
  "server_addr": "127.0.0.1:3306",
  "query": "select * from test",
  "duration_ns": 0
}`
)

func TestDecodeLine(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2019-04-10T15:08:00.727354+07:00")
	line.StartAt = tm
	line.EndAt = tm

	buf := bytes.NewReader(testData)

	l, err := decodeLine(buf)
	require.NoError(t, err)
	require.Equal(t, line, l)
}

func TestDecode(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2019-04-10T15:08:00.727354+07:00")
	line.StartAt = tm
	line.EndAt = tm

	buf := bytes.NewReader(testData)
	ls, err := Decode(buf)
	require.NoError(t, err)
	require.NotEmpty(t, ls)
	require.Len(t, ls, 1)
	require.Equal(t, []*LogLine{line}, ls)
}

func TestDecodeFile(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2019-04-10T15:08:00.727354+07:00")
	line.StartAt = tm
	line.EndAt = tm

	f, err := ioutil.TempFile("", "test.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(testData)
	if err != nil {
		panic(err)
	}

	ls, err := DecodeFile(f.Name())
	require.NoError(t, err)
	require.NotEmpty(t, ls)
	require.Len(t, ls, 1)
	require.Equal(t, []*LogLine{line}, ls)
}

func TestDecodeFileNegative(t *testing.T) {
	_, err := DecodeFile("")
	require.Error(t, err)
}

func TestLogLineString(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2019-04-10T15:08:00.727354+07:00")
	line.StartAt = tm
	line.EndAt = tm

	require.Equal(t, lineJSON, line.String())
}
