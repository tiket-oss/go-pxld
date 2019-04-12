# pxld
Decode your binary format ProxySQL query log.

[![Go Report Card](https://goreportcard.com/badge/github.com/tiket-oss/go-pxld)](https://goreportcard.com/report/github.com/tiket-oss/go-pxld)
[![Documentation](https://godoc.org/github.com/tiket-oss/go-pxld?status.svg)](http://godoc.org/github.com/tiket-oss/go-pxld)
[![license](https://img.shields.io/github/license/tiket-oss/go-pxld.svg)](https://github.com/tiket-oss/go-pxld/LICENSE)
[![Build Status](https://travis-ci.org/tiket-oss/go-pxld.svg?branch=master)](https://travis-ci.org/tiket-oss/go-pxld)

## Requirement

- Go v1.11 or later

## Build

- `go mod download`
- `cd cmd/decoder`
- `go build`

## Run

After you've build the decoder, just run the executable.

## The File Format

- `5D 00 00 00  00 00 00 00` first 8 bytes is the length of a message, this one.
- `00` next byte denounce that this is a ProxySQL log message or not, if not `00`, not valid.

The next after this line needs `read_encoded_length` function which itself needs `mysql_decode_length` function.

- `0A` this is the thread id because it is less than or equal `0xFB`, convert to `uint64`.
- `06` this is the length of username because it is less than or equal `0xFB`, convert to `uint64`.
- `64  69 64 61 73  79` is the username in ASCII.
- `12` this is the length of schema because it is less than or equal `0xFB`, convert to `uint64`.
- `69 6E  66 6F 72 6D  61 74 69 6F  6E 5F 73 63  68 65 6D 61` is the schema in ASCII.
- `0F` this is the length of client address because it is less than or equal `0xFB`, convert to `uint64`.
- `31 32 37  2E 30 2E 30  2E 31 3A 33  32 38 32 30` is the client address in ASCII.
- `FE` this tell us to read the next 8 bytes as `uint64`.
- `FF FF FF  FF FF FF FF  FF` this is the HID in `uint64`, if `HID == UINT64_MAX` then we don't have server address to read.
- `FE` this tell us to read the next 8 bytes as `uint64`.
- `91 00  B3 4A 28 86  05 00` this is query start time in UNIX microseconds in `uint64`.
- `FE` this tell us to read the next 8 bytes as `uint64`.
- `91 00  B3 4A 28 86  05 00` this is query end time in UNIX microseconds in `uint64`.
- `FE` this tell us to read the next 8 bytes as `uint64`.
- `D6 1F BA 14  4D 1F 23 AE` this is query digest in `uint64`, but it needs to be separated into two uint32 and then it can be printed into hex `sprintf("0x%X%X", n1, n2) == 0x14BA1FD6AE231F4D`.
- `0C` this is the length of the actual query because it is less than or equal `0xFB`, convert to `uint64`.
- `2E 30 2E  31 3A 33 32  38 32 30 00  A5` this the actual query in ASCII.

Then we can go to the next line and repeat.

### Decoding Message Parts Length

To decode part length, first we must take the first byte of the part. Then we check if:

- It is less or equal than `0xFB`. If so we take this byte as the message length.
- It is equal to `0xFC`. If so, we take 2 bytes as the message length.
- It is equal to `0xFD`. If so, we take 3 bytes as the message length.
- It is equal to `0xFE`. If so, we take 8 bytes as the message length.
