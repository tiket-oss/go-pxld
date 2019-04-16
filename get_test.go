package pxld

import (
	"bytes"
	"testing"
	"time"

	"encoding/binary"
	"github.com/stretchr/testify/require"
)

func TestGetMessageLength(t *testing.T) {
	data := []byte{
		50, 0, 0, 0, 0, 0, 0, 0,
	}
	buf := bytes.NewReader(data)

	n, err := GetMessageLength(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(50), n)
}

func TestGetMessageLengthNegative(t *testing.T) {
	buf := bytes.NewReader([]byte{})

	_, err := GetMessageLength(buf)
	require.Error(t, err)
}

func TestGetEncodedLength(t *testing.T) {
	// <= fb
	data := make([]byte, 1)
	data[0] = 0x08
	buf := bytes.NewReader(data)

	n, err := GetEncodedLength(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(8), n)

	// fc
	data = make([]byte, 3)
	data[0] = 0xFC
	data[1] = 0xFF
	data[2] = 0xFF
	buf = bytes.NewReader(data)

	n, err = GetEncodedLength(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(65535), n)

	// fd
	data = make([]byte, 4)
	data[0] = 0xFD
	data[1] = 0xFF
	data[2] = 0xFF
	data[3] = 0xFF
	buf = bytes.NewReader(data)

	n, err = GetEncodedLength(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(16777215), n)

	// fe
	data = make([]byte, 9)
	data[0] = 0xFE
	data[1] = 0xFF
	data[2] = 0xFF
	data[3] = 0xFF
	data[4] = 0xFF
	data[5] = 0xFF
	data[6] = 0xFF
	data[7] = 0xFF
	data[8] = 0xFF
	buf = bytes.NewReader(data)

	n, err = GetEncodedLength(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(18446744073709551615), n)
}

func TestGetEncodedLengthNegative(t *testing.T) {
	// <= fb
	buf := bytes.NewReader([]byte{})

	_, err := GetEncodedLength(buf)
	require.Error(t, err)

	// fc
	buf = bytes.NewReader([]byte{})
	_, err = GetEncodedLength(buf)
	require.Error(t, err)

	buf = bytes.NewReader([]byte{0xFC})
	_, err = GetEncodedLength(buf)
	require.Error(t, err)

	// fd
	buf = bytes.NewReader([]byte{})
	_, err = GetEncodedLength(buf)
	require.Error(t, err)

	buf = bytes.NewReader([]byte{0xFD})
	_, err = GetEncodedLength(buf)
	require.Error(t, err)

	// fe
	buf = bytes.NewReader([]byte{})
	_, err = GetEncodedLength(buf)
	require.Error(t, err)

	buf = bytes.NewReader([]byte{0xFE})
	_, err = GetEncodedLength(buf)
	require.Error(t, err)
}

func TestGetString(t *testing.T) {
	data := []byte{
		0x02,
	}
	data = append(data, []byte("ok")...)
	buf := bytes.NewReader(data)

	s, err := GetString(buf)
	require.NoError(t, err)
	require.Equal(t, "ok", s)
}

func TestGetStringNegative(t *testing.T) {
	data := []byte{}
	buf := bytes.NewReader(data)

	_, err := GetString(buf)
	require.Error(t, err)

	data = []byte{0x02}
	buf = bytes.NewReader(data)

	_, err = GetString(buf)
	require.Error(t, err)
}

func TestGetTime(t *testing.T) {
	now := time.Unix(time.Now().Unix(), 0)

	data := []byte{0xFE}
	timeData := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeData, uint64(now.Unix()*1000*1000))
	data = append(data, timeData...)
	buf := bytes.NewReader(data)

	tm, err := GetTime(buf)
	require.NoError(t, err)
	require.Equal(t, now, tm)
}

func TestGetTimeNegative(t *testing.T) {
	buf := bytes.NewReader([]byte{})

	_, err := GetTime(buf)
	require.Error(t, err)
}

func TestGetThreadID(t *testing.T) {
	data := []byte{0x10}
	buf := bytes.NewReader(data)

	tid, err := GetThreadID(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(16), tid)
}

func TestGetUsername(t *testing.T) {
	data := []byte{0x02}
	data = append(data, []byte("ok")...)
	buf := bytes.NewReader(data)

	tid, err := GetUsername(buf)
	require.NoError(t, err)
	require.Equal(t, "ok", tid)
}

func TestGetSchema(t *testing.T) {
	data := []byte{0x02}
	data = append(data, []byte("ok")...)
	buf := bytes.NewReader(data)

	tid, err := GetSchema(buf)
	require.NoError(t, err)
	require.Equal(t, "ok", tid)
}

func TestGetClientAddr(t *testing.T) {
	data := []byte{0x02}
	data = append(data, []byte("ok")...)
	buf := bytes.NewReader(data)

	tid, err := GetClientAddr(buf)
	require.NoError(t, err)
	require.Equal(t, "ok", tid)
}

func TestGetHID(t *testing.T) {
	data := []byte{0x10}
	buf := bytes.NewReader(data)

	tid, err := GetHID(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(16), tid)
}

func TestGetServerAddr(t *testing.T) {
	data := []byte{0x02}
	data = append(data, []byte("ok")...)
	buf := bytes.NewReader(data)

	tid, err := GetServerAddr(buf)
	require.NoError(t, err)
	require.Equal(t, "ok", tid)
}

func TestGetStartAt(t *testing.T) {
	now := time.Unix(time.Now().Unix(), 0)

	data := []byte{0xFE}
	timeData := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeData, uint64(now.Unix()*1000*1000))
	data = append(data, timeData...)
	buf := bytes.NewReader(data)

	tm, err := GetStartAt(buf)
	require.NoError(t, err)
	require.Equal(t, now, tm)
}

func TestGetEndAt(t *testing.T) {
	now := time.Unix(time.Now().Unix(), 0)

	data := []byte{0xFE}
	timeData := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeData, uint64(now.Unix()*1000*1000))
	data = append(data, timeData...)
	buf := bytes.NewReader(data)

	tm, err := GetEndAt(buf)
	require.NoError(t, err)
	require.Equal(t, now, tm)
}

func TestGetQuery(t *testing.T) {
	data := []byte{0x02}
	data = append(data, []byte("ok")...)
	buf := bytes.NewReader(data)

	tid, err := GetQuery(buf)
	require.NoError(t, err)
	require.Equal(t, "ok", tid)
}

func TestGetQueryDigest(t *testing.T) {
	data := []byte{0xFE}
	data = append(data, []byte{0xD6, 0x1F, 0xBA, 0x14, 0x4D, 0x1F, 0x23, 0xAE}...)

	buf := bytes.NewReader(data)

	digestStr, err := GetQueryDigest(buf)
	require.NoError(t, err)
	require.Equal(t, "0xD61FBA144D1F23AE", digestStr)
}

func TestGetQueryDigestNegative(t *testing.T) {
	data := []byte{}
	buf := bytes.NewReader(data)

	_, err := GetQueryDigest(buf)
	require.Error(t, err)
}

func TestGetMessage(t *testing.T) {
	data := []byte{
		0x10,
	}
	buf := bytes.NewReader(data)

	raw, b, err := GetMessage(1, buf)
	require.NoError(t, err)
	require.Equal(t, data, raw)

	resultData := make([]byte, 1)
	n, err := b.Read(resultData)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, data, resultData)
}

func TestMessageNegative(t *testing.T) {
	data := []byte{
		0x10,
	}
	buf := bytes.NewReader(data)

	_, _, err := GetMessage(2, buf)
	require.Error(t, err)

	buf = bytes.NewReader([]byte{})

	_, _, err = GetMessage(2, buf)
	require.Error(t, err)
}
