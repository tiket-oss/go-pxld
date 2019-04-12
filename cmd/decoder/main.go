package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tiket-oss/go-pxld"
	"gopkg.in/alecthomas/kingpin.v2"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	targetFile  = kingpin.Flag("target", "Target file to decode").Required().ExistingFile()
	output      = kingpin.Flag("output", "Output of this, can be file path, http address (60s timeout), or omit to stdout").Default("").String()
	repeatEvery = kingpin.Flag("repeat", "Repeat reading from the target file every n seconds, useful for reading logrotated file").Duration()
)

func main() {
	kingpin.Parse()

	log.Infof("Starting ProxySQL query log decoder")

	if *repeatEvery > 0 {
		t := time.Tick(*repeatEvery)

		for {
			do()

			<-t
		}
	} else {
		do()
	}

	log.Infof("Finished ProxySQL query log decoder")
}

func do() {
	logs, err := pxld.DecodeFile(*targetFile)
	if err != nil {
		log.Fatalf("Unexpected error while decoding file %s: %v", *targetFile, err)
	}

	if *output == "" {
		for _, l := range logs {
			fmt.Println(l)
		}
	} else {
		raw, err := json.Marshal(logs)
		if err != nil {
			log.Fatalf("Unexpected error while marshaling  file %s to JSON: %v", *targetFile, err)
		}

		if isValidURL(*output) {
			cli := &http.Client{
				Timeout: time.Minute,
			}
			res, err := cli.Post(*output, "application/json", bytes.NewReader(raw))
			if err != nil {
				log.Fatalf("Unexpected error while sending file %s JSON to %s: %v", *targetFile, *output, err)
			}
			if res.StatusCode >= http.StatusBadRequest {
				err = fmt.Errorf("invalid response status code %d", res.StatusCode)
				log.Fatalf("Unexpected error while sending file %s JSON to %s: %v", *targetFile, *output, err)
			}
		} else {
			err = ioutil.WriteFile(*output, raw, 0644)
			if err != nil {
				log.Fatalf("Unexpected error while writing file %s JSON to %s: %v", *targetFile, *output, err)
			}
		}
	}
}

func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	return err == nil
}
