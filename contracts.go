package space_traders_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Contract struct {
	ID            string
	FactionSymbol string
	Type          string
	Terms         struct {
		Deadline string
		Payment  struct {
			OnAccepted  int
			OnFulfilled int
		}
		Deliver []struct {
			TradeSymbol       string
			DestinationSymbol string
			UnitsRequired     int
			UnitsFulfilled    int
		}
	}
	Accepted         bool
	Fulfilled        bool
	Expiration       string
	DeadlineToAccept string
}

func (self Contract) String() string {
	ret := fmt.Sprintf(
		"Contract Expires %s\n"+
			"\t%s\n"+
			"\t%s\n"+
			"\tType\t%s\n"+
			"\tTerms\n"+
			"\t\tDeadline\t%s\n"+
			"\t\tUp Front\t%dc\n"+
			"\t\tFulfilled\t%dc\n"+
			"\t\tDeliver\t%v\n",
		self.Expiration,
		self.ID,
		self.FactionSymbol,
		self.Type,
		self.Terms.Deadline,
		self.Terms.Payment.OnAccepted,
		self.Terms.Payment.OnFulfilled,
		self.Terms.Deliver,
	)

	if !self.Accepted {
		ret += fmt.Sprintf(
			"\tDeadlineToAccept\t%s\n",
			self.DeadlineToAccept,
		)

	}

	if self.Fulfilled {
		ret += fmt.Sprintf("\tFULFILLED")
	}

	return ret
}

const BASE_URL = "https://api.spacetraders.io/v2"

func GetMyContracts(token string) ([]Contract, error) {
	errPrefix := "STAPI: Trying to get contracts."
	var contracts []Contract
	respObject := new(struct {
		Data  []Contract
		Meta  map[string]int
		Error *STJsonError
	})

	req, err := http.NewRequest(
		"GET",
		BASE_URL+"/my/contracts",
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		return []Contract{}, fmt.Errorf(
			"%s Creating request.\n%w",
			errPrefix,
			err,
		)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	for page := 1; respObject.Meta == nil || respObject.Meta["total"] > respObject.Meta["page"]*MAX_PAGE_LIMIT; page++ {
		q := req.URL.Query()
		q.Set("limit", strconv.Itoa(MAX_PAGE_LIMIT))
		q.Set("page", strconv.Itoa(page))
		req.URL.RawQuery = q.Encode()

		resp, err := SendRequest(req)
		if err != nil {
			return contracts, fmt.Errorf(
				"%s Trying to send GET request for page %d.\n%w",
				errPrefix,
				page,
				err,
			)
		}

		err = json.Unmarshal(resp, respObject)
		if err != nil {
			return contracts, fmt.Errorf(
				"%s Unmarshalling JSON.\n%w",
				errPrefix,
				err,
			)
		}
		if respObject.Error != nil {
			return contracts, fmt.Errorf(
				"%s spacetraders.io error on page %d.\n%w",
				errPrefix,
				page,
				respObject.Error,
			)
		}

		contracts = append(contracts, respObject.Data...)
		req.Body.Close()
		req.Body = io.NopCloser(strings.NewReader(""))
	}

	if contracts == nil {
		contracts = []Contract{{ID: "0"}}
		return contracts, fmt.Errorf(
			"%s %w No contracts.",
			errPrefix,
			NoContentError,
		)
	}

	return contracts, nil
}
