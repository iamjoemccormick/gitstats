package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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

func writeInfluxPoint(measurement string, tags map[string]string, fields map[string]int, time time.Time) error {

	client := influxdb2.NewClient(*dbUrl, *dbToken)
	writeAPI := client.WriteAPIBlocking(*dbOrg, *dbBucket)

	p := influxdb2.NewPointWithMeasurement(measurement)

	for k, v := range tags {
		p.AddTag(k, v)
	}

	for k, v := range fields {
		p.AddField(k, v)
	}

	p.SetTime(time)

	return writeAPI.WritePoint(context.Background(), p)

}
