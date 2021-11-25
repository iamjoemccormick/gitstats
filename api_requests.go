package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func githubApiGetRequest(url, repo string) []byte {

	req, _ := http.NewRequest("GET", githubApiBaseUrl+repo+url, nil)

	for _, header := range githubHeaders {
		req.Header.Add(header.header, header.value)
	}

	req.SetBasicAuth(*githubUsername, *githubToken)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(res)
		return nil
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}
