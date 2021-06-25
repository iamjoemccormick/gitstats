package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Clones struct {
	Count   int
	Uniques int
	Clones  []DailyClones
}

type DailyClones struct {
	Timestamp string
	Count     int
	Uniques   int
}

func githubApiGetRequest(url, repo string) string {

	req, _ := http.NewRequest("GET", githubApiBaseUrl+repo+url, nil)

	for _, header := range githubHeaders {
		req.Header.Add(header.header, header.value)
	}

	req.SetBasicAuth(*githubUsername, *githubToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(res)
		return ""
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
