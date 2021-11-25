package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
		"PopularPaths":     "/traffic/popular/paths",
		"PopularReferrers": "/traffic/popular/referrers",
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

		for api, url := range githubApis {
			log.Printf("Requesting %s (%s) for repo %s", api, url, repo)
			body := githubApiGetRequest(url, repo) // TODO: Add error handling if we receive an unexpected response. 

			if body[0] == 123 { // The json object contains a map.  

				var data map[string]interface{}
				if err := json.Unmarshal(body, &data); err != nil {	
					panic(err)
				}

				measurement := string(api)
				timestamp := time.Now()
				tags := map[string]string{"repo":repo}
				fields := map[string]float64{}
				
				for k, v := range data {

					switch t := v.(type) { 
					case string:
						tags[k] = v.(string)					
					case float64:
						fields[k] = v.(float64)
					case []interface{}:

						daily_measurement := "daily_" + string(k)
						
						for _, v2 := range v.([]interface{}) {

							dailyTags := map[string]string{"repo":repo}
							dailyFields := map[string]float64{}
							daily_timestamp := time.Now()

							for k3, v3 := range v2.(map[string]interface{}) {
								if k3 == "timestamp" {
									var err error
									daily_timestamp, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", v3.(string))
	
									if err != nil { 
										panic(err) // TODO: Handle this better.
									}									
								} else {
									switch t2 := v3.(type) {

									case float64:
										dailyFields[k3] = v3.(float64)
									default:
										fmt.Printf("Type of '%s' is %T\n", v3, t2)	//TODO: Add logging for currently unhandled data types.							

									}									
								}
							}
							// TODO: Write the following to InfluxDB. 							
							fmt.Println(daily_measurement)
							fmt.Println(daily_timestamp)
							fmt.Println(dailyTags)
							fmt.Println(dailyFields)
							//fmt.Println(v)
						}
												
					default:
						fmt.Printf("Type of '%s' is %T\n", k, t) //TODO: Add logging for currently unhandled data types.						
					}

				}

				//TODO: Write the following to InfluxDB. 
				fmt.Println(measurement)
				fmt.Println(timestamp) // TODO: Handle if the response from the API includes a timestamp. 
				fmt.Println(tags)
				fmt.Println(fields)
						

			} else if  body[0] == 91 { // The json object contains a slice of maps. 

				var data []map[string]interface{}
				if err := json.Unmarshal(body, &data); err != nil {	
					panic(err)
				}
				fmt.Println(data)
				
				//TODO: Handle parsing and writing to InfluxDB. 

			} else {
				log.Printf("Encountered an unsupported json object while processing endpoint '%s' (first byte of the response didn't match '[' or '{').", string(body[0]))
			}
			
		}

	}

}
