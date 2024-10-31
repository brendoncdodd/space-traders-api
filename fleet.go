package space_traders_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Ship struct {
	Symbol       string
	Registration *struct {
		Name          string
		FactionSymbol string
		Role          string
	}
	Nav  *ShipNav
	Crew *struct {
		Current  int
		Required int
		Capacity int
		Rotation string
		Morale   int
		Wages    int
	}
	Frame *struct {
		Symbol         string
		Name           string
		Description    string
		Condition      int
		Integrity      int
		ModuleSlots    int
		MountingPoints int
		FuelCapacity   int
		Requirements   *struct {
			Power int
			Crew  int
			Slots int
		}
	}
	Reactor *struct {
		Symbol       string
		Name         string
		Description  string
		Condition    int
		Integrity    int
		PowerOutput  int
		Requirements *struct {
			Power int
			Crew  int
			Slots int
		}
	}
	Engine *struct {
		Symbol       string
		Name         string
		Description  string
		Condition    int
		Integrity    int
		Speed        int
		Requirements *struct {
			Power int
			Crew  int
			Slots int
		}
	}
	Cooldown *struct {
		ShipSymbol       string
		TotalSeconds     int
		RemainingSeconds int
		Expiration       string
	}
	Modules []*struct {
		Symbol       string
		Capacity     int
		Range        int
		Name         string
		Description  string
		Requirements any //Schema says {}
	}
	Mounts []*struct {
		Symbol       string
		Name         string
		Description  string
		Strength     int
		Deposits     []any //Schema says [ null ]
		Requirements any   //Schema says {}
	}
	Cargo *struct {
		Capacity  int
		Units     int
		Inventory []any //Schema says [ {} ]
	}
	Fuel *struct {
		Current  int
		Capacity int
		Consumed *struct {
			Amount    int
			Timestamp string
		}
	}
}

func (self *Ship) String() string {
	return fmt.Sprintf(
		"Ship %s"+
			"\tRegistration [ %s, %s, %s ]"+
			"\tFuel %d/%d"+
			"\tCooldown %ds"+
			"\tCargo %d/%d"+
			"\t[%s (%d,%d) (%d,%d)]",
		self.Symbol,
		self.Registration.Name,
		self.Registration.FactionSymbol,
		self.Registration.Role,
		self.Fuel.Current,
		self.Fuel.Capacity,
		self.Cooldown.RemainingSeconds,
		self.Cargo.Units,
		self.Cargo.Capacity,
		self.Nav.WaypointSymbol,
		self.Nav.Route.Origin.X,
		self.Nav.Route.Origin.Y,
		self.Nav.Route.Destination.X,
		self.Nav.Route.Destination.Y,
	)
}

type ShipNav struct {
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

// Gets all of an agent's ships
func GetShipsByAgent(token string) (ships []Ship, err error) {
	errPrefix := "Getting agent's ships."
	respObject := new(struct {
		Data  []Ship
		Error *STJsonError
		Meta  map[string]int
	})

	req, err := http.NewRequest(
		"GET",
		"https://api.spacetraders.io/v2/my/ships",
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		return []Ship{}, fmt.Errorf(
			"%s Creating request.\n%w",
			errPrefix, err,
		)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	for page := 1; len(respObject.Meta) == 0 || respObject.Meta["total"] > respObject.Meta["page"]*respObject.Meta["limit"]; page++ {
		q := req.URL.Query()
		q.Set("limit", "20")
		q.Set("page", strconv.Itoa(page))
		req.URL.RawQuery = q.Encode()

		buf, err := SendRequest(req)
		if err != nil {
			return []Ship{}, fmt.Errorf(
				"%s Sending request for page %d.\n%w",
				errPrefix,
				page,
				err,
			)
		}

		err = json.Unmarshal(buf, respObject)
		if err != nil {
			return []Ship{}, fmt.Errorf(
				"%s Unmarshalling JSON page %d.\n%w",
				errPrefix,
				page,
				err,
			)
		}
		if respObject.Error != nil {
			return []Ship{}, fmt.Errorf(
				"%s spacetraders.io error for page %d.\n%w",
				errPrefix,
				page,
				respObject.Error,
			)
		}

		if len(respObject.Data) == 0 || respObject.Data[0].Symbol == "" {
			return ships, fmt.Errorf(
				"%s No ships on page %d",
				errPrefix,
				page,
			)
		}

		ships = append(ships, respObject.Data...)

		req.Body.Close()
		req.Body = io.NopCloser(strings.NewReader(""))
	}
	defer req.Body.Close()

	return ships, nil
}

// Gets a ship by symbol
func GetShip(shipSymbol string, token string) (ret *Ship, err error) {
	errPrefix := "Getting ship nav for " + shipSymbol + "."

	respObject := new(struct {
		Data  *Ship
		Error *STJsonError
	})

	req, err := http.NewRequest(
		"GET",
		"https://api.spacetraders.io/v2/my/ships/"+shipSymbol,
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"%s Creating ship nav request.\n%w",
			errPrefix,
			err,
		)
	}
	defer req.Body.Close()

	req.Header.Add("Authorization", "Bearer "+token)

	buf, err := SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf(
			"%s Sending request.\n%w",
			errPrefix,
			err,
		)
	}

	err = json.Unmarshal(buf, respObject)
	if err != nil {
		return nil, fmt.Errorf(
			"%s Decoding response body. %w",
			errPrefix,
			err,
		)
	}

	if respObject.Data == nil {
		err = fmt.Errorf("No data.")
	}

	if respObject.Error != nil {
		return respObject.Data, fmt.Errorf(
			"%s spacetraders.io error.\n%w",
			errPrefix,
			err,
		)
	}

	return respObject.Data, nil
}

// Gets a ship's nav by symbol
func GetShipNav(shipSymbol string, token string) (*ShipNav, error) {
	errPrefix := "Getting ship nav for " + shipSymbol + "."

	respObject := new(struct {
		Data  *ShipNav
		Error *STJsonError
	})

	req, err := http.NewRequest(
		"GET",
		"https://api.spacetraders.io/v2/my/ships/"+shipSymbol+"/nav",
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		return &ShipNav{}, fmt.Errorf(
			"%s Creating ship nav request.\n%w",
			errPrefix,
			err,
		)
	}
	defer req.Body.Close()

	req.Header.Add("Authorization", "Bearer "+token)

	buf, err := SendRequest(req)
	if err != nil {
		return &ShipNav{}, fmt.Errorf(
			"%s Sending request.\n%w",
			errPrefix,
			err,
		)
	}

	err = json.Unmarshal(buf, respObject)
	if err != nil {
		return &ShipNav{}, fmt.Errorf(
			"%s Decoding response body.\n%w",
			errPrefix,
			err,
		)
	}

	if respObject.Data == nil {
		err = fmt.Errorf("No data")
	}

	if respObject.Error != nil {
		return respObject.Data, fmt.Errorf(
			"%s spacetraders.io error: %w",
			errPrefix,
			err,
		)
	}

	return respObject.Data, err
}

// Identifies the waypoint where a ship is located,
// then returns the coords of the waypoint.
// I'm not sure how this behaves for ships in transit.
func GetShipLocation(shipSymbol string, token string) (Vector2, error) {
	errPrefix := "Getting ship location."

	nav, err := GetShipNav(shipSymbol, token)
	if err != nil {
		return Vector2{}, fmt.Errorf(
			"%s Getting ship nav.\n%w",
			errPrefix, err,
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
