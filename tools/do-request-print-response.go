package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	st "github.com/brendoncdodd/space_traders_api"
)

var res string

func main() {
	errPrefix := "ERROR:"

	flag.StringVar(
		&res,
		"res",
		"/",
		"The path of a resource to append to \"https://spacetraders.io/v2\"",
	)

	flag.Parse()

	req, err := http.NewRequest(
		"GET",
		"https://api.spacetraders.io/v2"+res,
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		err = fmt.Errorf(
			"%s Creating GET request.\n%w",
			errPrefix,
			err,
		)
		log.Fatal(err.Error())
	}
	defer req.Body.Close()

	resp, err := st.SendRequest(req)
	if err != nil {
		err = fmt.Errorf(
			"%s Sending request.\n%w",
			errPrefix,
			err,
		)
		log.Fatal(err.Error())
	}

	fmt.Printf(
		"%s",
		string(resp),
	)
}
