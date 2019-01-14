package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getRequest() {
	response, err := http.Get("http://golang.org/")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}
}

func main() {
	// getRequest()
	parseJSON()
	// dbExample()
}
