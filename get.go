package pxld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

func GetMessageLength(dataStream io.Reader) (dataLength uint64, err error) {
	data := make([]byte, 8)

	_, err = dataStream.Read(data)
	if err != nil {
		return
	}

	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &dataLength)

	return
}

func GetThreadID(dataStream io.Reader) (threadID uint64, err error) {
	return GetEncodedLength(dataStream)
}

func GetUsername(dataStream io.Reader) (username string, err error) {
	return GetString(dataStream)
}

func GetSchema(dataStream io.Reader) (schema string, err error) {
	return GetString(dataStream)
}

func GetClientAddr(dataStream io.Reader) (clientAddr string, err error) {
	return GetString(dataStream)
}

func GetHID(dataStream io.Reader) (hid uint64, err error) {
	return GetEncodedLength(dataStream)
}

func GetServerAddr(dataStream io.Reader) (serverAddr string, err error) {
	return GetString(dataStream)
}

func GetStartAt(dataStream io.Reader) (startAt time.Time, err error) {
	return GetTime(dataStream)
}

func GetEndAt(dataStream io.Reader) (endAt time.Time, err error) {
	return GetTime(dataStream)
}

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

func GetQuery(dataStream io.Reader) (q string, err error) {
	return GetString(dataStream)
}

func GetTime(dataStream io.Reader) (t time.Time, err error) {
	var unixMicrosecond uint64
	unixMicrosecond, err = GetEncodedLength(dataStream)
	if err != nil {
		return
	}

	t = time.Unix(0, int64(unixMicrosecond*1000))

	return
}

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
