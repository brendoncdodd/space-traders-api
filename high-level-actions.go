package space_traders_api

import (
	"fmt"
	"math"
)

// Not implemented yet
// TODO: Get symbols and locations of waypoint with traits.
// TODO: Remove lines that assign test values.
func FindNearestWaypointWithTraits(
	shipSymbol string,
	traits []string,
	token string,
) (
	waypointSymbol string,
	err error,
) {
	errPrefix := "Trying to find nearest waypoint with traits."
	minDistance := math.Inf(1)
	waypointLocations := make(map[string]Vector2)

	shipLocation, err := GetShipLocation(shipSymbol, token)
	if err != nil {
		fmt.Errorf(
			"%s Getting ship location.%w",
			errPrefix,
			err,
		)
	}

	waypoints, err := GetSystemWaypoints("", traits, "")
	if err != nil {
		return waypointSymbol, fmt.Errorf(
			"%s Getting waypoints.%w",
			errPrefix,
			err,
		)
	}

	for _, waypoint := range waypoints {
		waypointLocations[waypoint.Symbol] = Vector2{
			waypoint.X,
			waypoint.Y,
		}
	}

	for waypoint, location := range waypointLocations {
		d := shipLocation.Distance(location)
		if d < minDistance {
			waypointSymbol = waypoint
			minDistance = d
		}
	}

	return
}
