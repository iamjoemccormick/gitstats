package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func loadFileToSlice(path string) (output []string) {

	contents, err := ioutil.ReadFile(path)

	if err != nil {
		log.Printf("The following error occured while attempting to read the file at %s: %s", path, err)
		os.Exit(1)
	}

	output = strings.Split(string(contents), "\n")

	return output
}
