package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type requestHeader struct {
	header string
	value  string
}

var (
	githubUsername  = flag.String("username", "", "Username for the GitHub REST API.")
	githubToken     = flag.String("token", "", "Token for the GitHub REST API.")
	gitHubReposFile = flag.String("repos", "", "Path to a file containing a list of GitHub repositories to monitor in the format <user>/<repo> (separated by newlines).")
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
)

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)

	githubRepos = loadFileToSlice(*gitHubReposFile)

	//TODO:
	// Need to determine how we want to store data in InfluxDB.
	// For example do we want different buckets for each GitHub repo?

	for _, repo := range githubRepos {
		fmt.Println(repo)
		githubApiGetRequest(githubApis["Views"], repo)
	}

}
