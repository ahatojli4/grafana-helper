package cmd

import (
	"fmt"
	"net/http"
	"net/url"
)

func checkCredentials(host, authString string) bool {
	c := http.Client{}
	response, err := c.Do(&http.Request{
		Method: "GET",
		URL:    &url.URL{Scheme: "https", Host: host, Path: "/api/org/"},
		Header: map[string][]string{
			"Authorization": {authString},
		},
	})
	if err != nil {
		panic(fmt.Errorf("error while checking credentials: %s", err))
	}
	success := response.StatusCode == 200
	if !success {
		fmt.Println("probably wrong credentials")
		fmt.Println(fmt.Sprintf("grafana server answer status code is %d", response.StatusCode))
	}

	return success
}
