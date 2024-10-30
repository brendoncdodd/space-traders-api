package space_traders_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const MAX_PAGE_LIMIT = "20"

type WaypointTrait struct {
	Symbol      string
	Name        string
	Description string
}

type WaypointModifier struct {
	Symbol      string
	Name        string
	Description string
}

type WaypointChart struct {
	WaypointSymbol string
	Submittedby    string
	SubmittedOn    string
}

type Waypoint struct {
	Symbol              string
	Type                string
	SystemSymbol        string
	X                   int
	Y                   int
	Orbitals            *struct{ symbol string }
	Orbits              string
	Faction             *struct{ symbol string }
	Traits              []WaypointTrait
	Modifires           []WaypointModifier
	Chart               *WaypointChart
	IsUnderConstruction bool
}

func GetAllWaypointsInSystem(systemSymbol string) (ret []Waypoint, err error) {
	errPrefix := fmt.Sprintf(
		"Getting all waypoints in system %s.",
		systemSymbol,
	)

	ret, err = GetSystemWaypoints(systemSymbol, []string{}, "")
	if err != nil {
		return ret, fmt.Errorf(
			"%s Getting waypoints in system %s.%w",
			errPrefix,
			systemSymbol,
			err,
		)
	}

	return
}

func GetSystemWaypoints(
	systemSymbol string,
	traits []string,
	waypointType string,
) (ret []Waypoint, err error) {
	var pageWaypoints struct {
		data []Waypoint
		meta map[string]int
	}

	errPrefix := fmt.Sprintf("Getting system waypoints:\n\ttraits\t%v\n\ttype\t%s\n",
		traits,
		waypointType,
	)

	req, err := http.NewRequest(
		"GET",
		"https://api.spacetraders.io/v2/systems/"+systemSymbol+"/waypoints",
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		return ret, fmt.Errorf(
			"%s Creating waypoints request. %w",
			errPrefix,
			err,
		)
	}
	defer req.Body.Close()

	q := req.URL.Query()
	q.Add("limit", MAX_PAGE_LIMIT)

	if traits != nil && len(traits) > 0 {
		for _, trait := range traits {
			q.Add("traits", trait)
		}
	}

	if waypointType != "" {
		q.Add("type", waypointType)
	}

	req.URL.RawQuery = q.Encode()
	req.Close = true

	buf, err := SendRequest(req)
	if err != nil {
		return ret, fmt.Errorf(
			"%s Sending request for page 1.\n%w",
			errPrefix,
			err,
		)
	}
	req.Body.Close()
	req.Body = io.NopCloser(strings.NewReader(""))
	defer req.Body.Close()

	err = json.Unmarshal(buf, &pageWaypoints)
	if err != nil {
		return ret, fmt.Errorf(
			"%s Unmarshaling JSON for page 1.\n\tlength: %d\n\tJSON:%s\n%w",
			errPrefix,
			len(buf),
			string(buf),
			err,
		)
	}

	ret = append(ret, pageWaypoints.data...)

	for pageWaypoints.meta["page"]*pageWaypoints.meta["limit"] < pageWaypoints.meta["total"] {
		if q.Has("page") {
			q.Del("page")
		}
		q.Add("page", strconv.Itoa(pageWaypoints.meta["page"]+1))

		req.URL.RawQuery = q.Encode()

		buf, err = SendRequest(req)
		if err != nil {
			return ret, fmt.Errorf(
				"%s Sending request for page %d.%w",
				errPrefix,
				pageWaypoints.meta["page"]+1,
				err,
			)
		}

		err = json.Unmarshal(buf, pageWaypoints)
		if err != nil {
			return ret, fmt.Errorf(
				"%s Unmarshaling JSON.\n\t%v%w",
				errPrefix,
				pageWaypoints,
				err,
			)
		}

		ret = append(ret, pageWaypoints.data...)
	}

	return
}

// https://api.spacetraders.io/v2/systems/{systemSymbol}/waypoints/{waypointSymbol}
func GetWaypointLocation(waypointSymbol string) (Vector2, error) {
	var waypoint *Waypoint
	errPrefix := "Getting waypoint location."
	systemSymbol :=
		strings.Split(waypointSymbol, "-")[0] +
			strings.Split(waypointSymbol, "-")[1]

	req, err := http.NewRequest(
		"GET",
		"https://api.spacetraders.io/v2/systems/"+
			systemSymbol+
			"/waypoints/"+
			waypointSymbol,
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		return Vector2{}, fmt.Errorf(
			"%s Creating request. %w",
			errPrefix,
			err,
		)
	}

	buf, err := SendRequest(req)
	if err != nil {
		return Vector2{}, fmt.Errorf(
			"%s Sending request.%w",
			errPrefix,
			err,
		)
	}

	err = json.Unmarshal(buf, waypoint)
	if err != nil {
		return Vector2{}, fmt.Errorf(
			"%s Unmarshaling response.%w",
			errPrefix,
			err,
		)
	}

	return Vector2{waypoint.X, waypoint.Y}, nil
}
