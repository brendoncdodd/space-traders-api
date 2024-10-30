package space_traders_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"io"
)

type Nav struct {
	SystemSymbol   string
	WaypointSymbol string
	Route          struct {
		Destination struct {
			Symbol       string
			Type         string
			SystemSymbol string
			X            int
			Y            int
		}
		Origin struct {
			Symbol       string
			Type         string
			SystemSymbol string
			X            int
			Y            int
		}
		DepartureTime string
		Arrival       string
	}
	Status     string
	FlightMode string
}

// Identifies the waypoint where a ship is located,
// then returns the coords of the waypoint.
func GetShipLocation(shipSymbol string, token string) (Vector2, error) {
	const BUFFER_SIZE = 10000
	errPrefix := "Getting ship location."
	buf := make([]byte, BUFFER_SIZE)
	var nav *Nav

	req, err := http.NewRequest(
		"GET",
		"https://api.spacetraders.io/v2/my/ships/"+shipSymbol+"/nav",
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		return Vector2{0, 0}, fmt.Errorf(
			"%s Creating ship nave request. %w",
			errPrefix,
			err,
		)
	}
	defer req.Body.Close()

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Vector2{0, 0}, fmt.Errorf(
			"%s Sending ship nav request. %w",
			errPrefix,
			err,
		)
	}
	bodySize, err := resp.Body.Read(buf)
	if err != nil {
		return Vector2{0, 0}, fmt.Errorf(
			"%s Reading response body into buffer. %w",
			errPrefix,
			err,
		)
	}
	defer resp.Body.Close()
	if bodySize >= BUFFER_SIZE {
		return Vector2{0, 0}, fmt.Errorf(
			"%s Reading response body into buffer. %s%d/%d",
			errPrefix,
			"Response too big for buffer: ",
			bodySize,
			BUFFER_SIZE,
		)
	}
	buf = bytes.TrimRight(buf, "\x00")

	err = json.Unmarshal(buf, nav)
	if err != nil {
		return Vector2{0, 0}, fmt.Errorf(
			"%s Decoding response body. %w",
			errPrefix,
			err,
		)
	}

	shipLocation, err := GetWaypointLocation(nav.WaypointSymbol)
	if err != nil {
		return shipLocation, fmt.Errorf(
			"%s Getting ship's waypoint location. %w",
			errPrefix,
			err,
		)
	}
	return shipLocation, nil
}
