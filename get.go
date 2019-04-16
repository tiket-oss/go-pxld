package pxld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// GetMessageLength is used to get message data length
func GetMessageLength(dataStream io.Reader) (dataLength uint64, err error) {
	data := make([]byte, 8)

	_, err = dataStream.Read(data)
	if err != nil {
		return
	}

	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &dataLength)

	return
}

// GetThreadID is used to obtain thread id
func GetThreadID(dataStream io.Reader) (threadID uint64, err error) {
	return GetEncodedLength(dataStream)
}

// GetUsername is used to obtain querier username
func GetUsername(dataStream io.Reader) (username string, err error) {
	return GetString(dataStream)
}

// GetSchema is used to get schema of a query
func GetSchema(dataStream io.Reader) (schema string, err error) {
	return GetString(dataStream)
}

// GetClientAddr is used to get client address of querier
func GetClientAddr(dataStream io.Reader) (clientAddr string, err error) {
	return GetString(dataStream)
}

// GetHID is used to get HID of a query, if not 0xFFFFFFFF we can get server address
func GetHID(dataStream io.Reader) (hid uint64, err error) {
	return GetEncodedLength(dataStream)
}

// GetServerAddr is used to get target server address for the query
func GetServerAddr(dataStream io.Reader) (serverAddr string, err error) {
	return GetString(dataStream)
}

// GetStartAt is used to get query start time in UNIX microseconds
func GetStartAt(dataStream io.Reader) (startAt time.Time, err error) {
	return GetTime(dataStream)
}

// GetEndAt is used to get query end time in UNIX microseconds
func GetEndAt(dataStream io.Reader) (endAt time.Time, err error) {
	return GetTime(dataStream)
}

// GetQueryDigest is used to get query's digest
func GetQueryDigest(dataStream io.Reader) (digest string, err error) {
	var digestRaw uint64
	digestRaw, err = GetEncodedLength(dataStream)
	if err != nil {
		return
	}

	// first we need to convert it back into []byte
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, digestRaw)

	// then turn into hex string
	digest = fmt.Sprintf("0x%X", buf)

	return
}

// GetQuery is used to get string representation of the query
func GetQuery(dataStream io.Reader) (q string, err error) {
	return GetString(dataStream)
}

// GetTime is used to get time from query log data
func GetTime(dataStream io.Reader) (t time.Time, err error) {
	var unixMicrosecond uint64
	unixMicrosecond, err = GetEncodedLength(dataStream)
	if err != nil {
		return
	}

	t = time.Unix(0, int64(unixMicrosecond*1000))

	return
}

// GetString is used to get string from query log data
func GetString(dataStream io.Reader) (s string, err error) {
	var n uint64
	n, err = GetEncodedLength(dataStream)
	if err != nil {
		return
	}

	raw := make([]byte, n)
	_, err = dataStream.Read(raw)
	if err != nil {
		return
	}

	s = string(raw)

	return
}

// GetEncodedLength is used to get message length from query log data
func GetEncodedLength(dataStream io.Reader) (ln uint64, err error) {
	// get first byte
	lenFlag := make([]byte, 1)

	_, err = dataStream.Read(lenFlag)
	if err != nil {
		return
	}

	if lenFlag[0] <= 0xFB {
		// just use this value as the actual length
		tmp := append(lenFlag, 0, 0, 0, 0, 0, 0, 0)
		err = binary.Read(bytes.NewReader(tmp), binary.LittleEndian, &ln)
	} else if lenFlag[0] == 0xFC {
		// get 2 bytes from data stream
		tmp := make([]byte, 2)
		_, err = dataStream.Read(tmp)
		if err != nil {
			return
		}

		tmp = append(tmp, 0, 0, 0, 0, 0, 0)
		err = binary.Read(bytes.NewReader(tmp), binary.LittleEndian, &ln)
	} else if lenFlag[0] == 0xFD {
		// get 3 bytes from data stream
		tmp := make([]byte, 3)
		_, err = dataStream.Read(tmp)
		if err != nil {
			return
		}

		tmp = append(tmp, 0, 0, 0, 0, 0)
		err = binary.Read(bytes.NewReader(tmp), binary.LittleEndian, &ln)
	} else if lenFlag[0] == 0xFE {
		// get 8 bytes from data stream
		tmp := make([]byte, 8)
		_, err = dataStream.Read(tmp)
		if err != nil {
			return
		}

		err = binary.Read(bytes.NewReader(tmp), binary.LittleEndian, &ln)
	}

	return
}

// GetMessage is used to get the message from query log data
func GetMessage(messageLength uint64, dataStream io.Reader) (raw []byte, buf io.Reader, err error) {
	raw = make([]byte, int(messageLength))

	var n int
	n, err = dataStream.Read(raw)
	if err != nil {
		return
	}
	if n != int(messageLength) {
		err = fmt.Errorf("failed to read %d bytes, read %d bytes instead", messageLength, n)
		return
	}

	buf = bytes.NewReader(raw)

	return
}
