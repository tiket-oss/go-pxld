package pxld

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"time"
)

type LogLine struct {
	MessageLength uint64        `json:"message_length"`
	RawMessage    []byte        `json:"raw_message"` // this is without message length data prepended
	ThreadID      uint64        `json:"thread_id"`
	Username      string        `json:"username"`
	Schema        string        `json:"schema"`
	StartAt       time.Time     `json:"start_at"`
	EndAt         time.Time     `json:"end_at"`
	QueryDigest   string        `json:"query_digest"`
	HID           uint64        `json:"hid,omitempty"`
	ClientAddr    string        `json:"client_addr"`
	ServerAddr    string        `json:"server_addr,omitempty"` // this depends on HID value
	Query         string        `json:"query"`
	Duration      time.Duration `json:"duration_ns"`
}

func (l *LogLine) String() string {
	raw, _ := json.MarshalIndent(l, "", "  ")

	return string(raw)
}

func Decode(r io.Reader) (l []*LogLine, err error) {
	l = []*LogLine{}

	// read until encountering EOF or unexpected error
	for {
		var line *LogLine
		line, err = decodeLine(r)
		if err != nil {
			break
		}

		l = append(l, line)
	}
	if err == io.EOF {
		err = nil
	}

	return
}

func DecodeFile(fp string) (l []*LogLine, err error) {
	var f *os.File
	f, err = os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	defer f.Close()

	return Decode(f)
}

func decodeLine(dataStream io.Reader) (line *LogLine, err error) {
	line = &LogLine{}

	// first read message length, this is an uint64, so 8 bytes
	// somehow turn this 8 bytes into uint64
	// this consume first 8 bytes of the buffer
	line.MessageLength, err = GetMessageLength(dataStream)
	if err != nil {
		return
	}

	// read all the message and replace dataStream
	line.RawMessage, dataStream, err = GetMessage(line.MessageLength, dataStream)
	if err != nil {
		return
	}

	// then consume the next 1 byte, if 0 proceed, if not 0
	// then just return with error as this is not a valid ProxySQL Query Log
	err = IsProxySQLQuery(dataStream)
	if err != nil {
		return
	}

	// then read thread id
	line.ThreadID, err = GetThreadID(dataStream)
	if err != nil {
		return
	}

	// then username
	line.Username, err = GetUsername(dataStream)
	if err != nil {
		return
	}

	// then schema name
	line.Schema, err = GetSchema(dataStream)
	if err != nil {
		return
	}

	// then client addr
	line.ClientAddr, err = GetClientAddr(dataStream)
	if err != nil {
		return
	}

	// then HID
	line.HID, err = GetHID(dataStream)
	if err != nil {
		return
	}

	// if HID not null, read server addr
	// HID is null if the same as maximum of uint64
	if line.HID != math.MaxUint64 {
		line.ServerAddr, err = GetServerAddr(dataStream)
		if err != nil {
			return
		}
	}

	// then start time
	line.StartAt, err = GetStartAt(dataStream)
	if err != nil {
		return
	}

	// then end time
	line.EndAt, err = GetEndAt(dataStream)
	if err != nil {
		return
	}

	// then calculate duration
	line.Duration = line.EndAt.Sub(line.StartAt)

	// then query digest
	line.QueryDigest, err = GetQueryDigest(dataStream)
	if err != nil {
		return
	}

	// then get the actual query
	line.Query, err = GetQuery(dataStream)
	if err != nil {
		return
	}

	return
}
