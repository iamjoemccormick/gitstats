package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"
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

	for _, repo := range githubRepos {

		response := githubApiGetRequest(githubApis["Clones"], repo)
		var clones Clones
		json.Unmarshal([]byte(response), &clones)

		tags := map[string]string{"Repo": repo}
		fields := map[string]int{"Count": clones.Count, "Uniques": clones.Uniques}

		writeInfluxPoint("LifetimeClones", tags, fields, time.Now())

		for _, c := range clones.Clones {
			t, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", c.Timestamp)

			if err == nil {
				writeInfluxPoint("DailyClones", map[string]string{"Repo": repo}, map[string]int{"Count": c.Count, "Unique": c.Uniques}, t)
			} else {
				log.Printf("Error parsing timestamp '%s' for InfluxPoint 'DailyClones'.", err)
			}

		}

		//TODO: Make this functionality reusable and expand it to the other endpoints.

		// for api, url := range githubApis {
		// 	fmt.Println(api, url)
		// 	body := githubApiGetRequest(url, repo)
		// 	fmt.Println(body)
		// }

	}

}
