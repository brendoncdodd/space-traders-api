package space_traders_api

import (
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"bytes"
	"io"
	"encoding/json"
)

func GetContracts(requestTemplate *http.Request) ([]Contract, error) {
	const BUFFER_SIZE = 10000
	error_prefix := "STAPI: Trying to get contracts."
	responseBody := make([]byte, BUFFER_SIZE)
	var contracts []Contract

	if requestTemplate == nil {
		if token_GET == nil {
			return []Contract{}, fmt.Errorf(
				"%s token_GET is nil",
				error_prefix,
			)
		}
		requestTemplate = token_GET
	}

	req := requestTemplate.Clone(requestTemplate.Context())
	req.Body = io.NopCloser(strings.NewReader(""))
	req.URL.Path = "/v2/my/contracts"
	defer req.Body.Close()

	currentPage := 1

	q := req.URL.Query()
	q.Add("limit", "20")
	req.URL.RawQuery = q.Encode()

	for lastPage := false; !lastPage; {
		page := new(ContractPage)

		q = req.URL.Query()
		if q.Has("page") {
			q.Del("page")
		}
		q.Add("page", strconv.Itoa(currentPage))
		req.URL.RawQuery = q.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return contracts, fmt.Errorf(
				"%s Trying to send GET request: %s %w",
				req.URL.String(),
				error_prefix,
				err,
			)
		}
		defer resp.Body.Close()

		bodySize, err := resp.Body.Read(responseBody)
		responseBody = bytes.TrimRight(responseBody, "\x00")
		if err != nil {
			return contracts, fmt.Errorf(
				"%s Trying to read response body. %w",
				error_prefix,
				err,
			)
		}
		if bodySize >= BUFFER_SIZE {
			return contracts, fmt.Errorf(
				"%s Response body too big for buffer (%d bytes).",
				error_prefix,
				BUFFER_SIZE,
			)
		}

		err = json.Unmarshal(responseBody, page)
		if err != nil {
			return contracts, fmt.Errorf(
				"%s Unmarshalling JSON. %w",
				error_prefix,
				err,
			)
		}

		currentPage++
		if currentPage > page.Meta["total"] {
			lastPage = true
		}

		contracts = append(contracts, page.Data...)
	}

	if contracts == nil {
		contracts = []Contract{{ID: "0"}}
		return contracts, fmt.Errorf(
			"%s %w No contracts.",
			error_prefix,
			NoContentError,
		)
	}

	return contracts, nil
}
