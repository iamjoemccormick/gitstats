package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type requestHeader struct {
	header string
	value  string
}

var (
	githubUsername  = flag.String("username", "", "Username for the GitHub REST API.")
	githubToken     = flag.String("token", "", "Token for the GitHub REST API.")
	gitHubReposFile = flag.String("config", "", "Path to a file containing a list of GitHub repositories to monitor in the format <user>/<repo> (separated by newlines).")
	githubHeaders   = []requestHeader{
		{header: "Accept", value: "application/vnd.github.v3+json"},
	}
	githubApiBaseUrl = "https://api.github.com/repos/"
	githubApis       = map[string]string{
		"Clones":           "/traffic/clones",
		"PopularPaths":     "/popular/paths",
		"PopularReferrers": "/popular/referrers",
		"Views":            "/traffic/views",
	}
	githubRepos = []string{}
	dbUrl       = flag.String("db_url", "", "InfluxDB Server URL (in the format http//url:port).")
	dbToken     = flag.String("db_token", "", "InfluxDB Token")
	dbOrg       = flag.String("db_org", "", "InfluxDB Organization")
	dbBucket    = flag.String("db_bucket", "", "InfluxDB Bucket")
)

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)

	githubRepos = loadFileToSlice(*gitHubReposFile)

	client := influxdb2.NewClient(*dbUrl, *dbToken)
	writeAPI := client.WriteAPIBlocking(*dbOrg, *dbBucket)

	for _, repo := range githubRepos {

		response := githubApiGetRequest(githubApis["Clones"], repo)
		var clones Clones
		json.Unmarshal([]byte(response), &clones)

		p := influxdb2.NewPointWithMeasurement("LifetimeClones").
			AddTag("Repo", repo).
			AddField("Count", clones.Count).
			AddField("Uniques", clones.Uniques).
			SetTime(time.Now())

		writeAPI.WritePoint(context.Background(), p)

		//TODO: Make this functionality reusable and expand it to the other endpoints.

		// for api, url := range githubApis {
		// 	fmt.Println(api, url)
		// 	body := githubApiGetRequest(url, repo)
		// 	fmt.Println(body)
		// }

	}

}
