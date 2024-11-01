package space_traders_api

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

var new_save SaveData
var test_user_first_ship_symbol string

const TEST_USER_TOKEN = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGlmaWVyIjoiVEVTVF9VU0VSIiwidmVyc2lvbiI6InYyLjIuMCIsInJlc2V0X2RhdGUiOiIyMDI0LTEwLTI3IiwiaWF0IjoxNzMwMjk1Mzc2LCJzdWIiOiJhZ2VudC10b2tlbiJ9.G_gv8evLdVdAbcNpryY5lyc8TjdMhMuHm14LA7_7nHsiVJzxO9lEfIXcemAy0fHv2Yqw-kF5zNevMBqw0aCUXADhdBwq2mhPPVhGoZ_TGTy5AycIvjnhXH5kVSyhgYB1HlDYV5REn3IkNCTSbE1SEjiVLKRmH5GAvh0L4kv4cu41re015fiwD_zt89_K85FQ_2Q0WlkWUgmA6js9wEp4D8M9omYY7TsS0r5pi5_pyORVca5hDt9lKtj5LqtbRDJ-oToyDoKT5Almu_ilEzc0Ajz1j3Th9_EvguH6R8GAb53qK6Tx7a4Q6qMd0LjbrB1O0ZswyXd51kWplR3sr1R4DQ"

const TEST_USER_ACCOUNT_ID = "cm2vx6x981gqys60c8deszydq"

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// timer returns a function that prints the name argument and
// the elapsed time between the call to timer and the call to
// the returned function. The returned function is intended to
// be used in a defer statement:
//
//	defer timer("sum")()
func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

// Helpful if you're not sure what's in an interface{} or something.
func describe(i interface{}) string {
	return fmt.Sprintf("(%v, %T)", i, i)
}

func TestVector2(t *testing.T) {
	errPrefix := "TEST_Vector2:"

	a := Vector2{3, 0}
	b := Vector2{0, 4}
	want := 5.0

	timerDistance := timer("Vector2.Distance()")
	res := a.Distance(b)
	timerDistance()

	if res != want {
		t.Fatalf(
			"%s Bad distance\ta:%v b:%v d:%f want:%f",
			errPrefix,
			a,
			b,
			res,
			want,
		)
	}
}

func TestCreateAgent(t *testing.T) {
	var err error
	errPrefix := "TEST_CreateAgent:"
	agentName := []byte("dodd_")
	for i := 0; i < 8; i++ {
		randIdx := rand.Intn(len(CHARSET))
		agentName = append(agentName, CHARSET[randIdx])
	}

	timerCreateAgent := timer("CreateAgent(string(agentName), \"COSMIC\")")

	new_save, err = CreateAgent(string(agentName), "COSMIC")
	if err != nil {
		t.Fatalf(
			"%s Creating Agent.\n%s\n%s",
			errPrefix,
			err.Error(),
			describe(new_save),
		)
	}

	timerCreateAgent()

	fmt.Printf("Created agent:\n%s\n---BEGIN TOKEN---\n%s\n---END TOKEN---\n",
		new_save.Agent,
		new_save.Token,
	)
}

// method is a valid method for http.NewRequest()
// Caller should close the .Body of the generate request.
func generateTestUserTemplate(method string) (*http.Request, error) {
	errPrefix := "generateTestUserTemplate():"
	//TEST_USER
	token := TEST_USER_TOKEN

	req, err := http.NewRequest(
		method,
		"https://api.spacetraders.io",
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"%s Error creating template request. %w",
			errPrefix,
			err,
		)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}

func TestGetAgentDetails(t *testing.T) {
	errPrefix := "TEST_GetAgentDetails:"

	req, err := generateTestUserTemplate("GET")
	if err != nil {
		t.Fatalf(
			"%s Error creating template request. %s",
			errPrefix,
			err.Error(),
		)
	}
	defer req.Body.Close()

	timerGetAgentDetails := timer("First call to GetAgentDetails()")

	result, resultJSON, err := GetAgentDetails(req)
	if err != nil {
		t.Fatalf(
			"%s First call returned error. %s\nObject: %v\nStatus: %s",
			errPrefix,
			err.Error(),
			result,
			resultJSON,
		)
	}

	timerGetAgentDetails()

	if result.Symbol != "TEST_USER" ||
		result.AccountID != TEST_USER_ACCOUNT_ID {
		t.Fatalf(
			"%s Result of first call is not TEST_USER. User is %s\naccountID is %s",
			errPrefix,
			result.Symbol,
			result.AccountID,
		)
	}

	result = Agent{}

	timerGetAgentDetails = timer("Second call to GetAgentDetails()")

	//Do a second call using the same template to make sure we don't mutate the template or something.
	result, resultJSON, err = GetAgentDetails(req)
	if err != nil {
		t.Fatalf(
			"%s Second call returned error. %s\nObject: %v\nStatus: %s",
			errPrefix,
			err.Error(),
			result,
			resultJSON,
		)
	}

	timerGetAgentDetails()

	if result.Symbol != "TEST_USER" ||
		result.AccountID != TEST_USER_ACCOUNT_ID {
		t.Fatalf(
			"%s Result of second call is not TEST_USER. User is %s\naccountID is %s",
			errPrefix,
			result.Symbol,
			result.AccountID,
		)
	}
}

func TestGetShipsByAgent(t *testing.T) {
	errPrefix := "TEST_GetShipsByAgent:"

	timerGetShipsByAgent := timer("GetShipsByAgent(TEST_USER_TOKEN)")
	ships, err := GetShipsByAgent(TEST_USER_TOKEN)
	timerGetShipsByAgent()
	if err != nil {
		t.Fatalf(
			"%s Error from GetShipsByAgent(),\n%s",
			errPrefix,
			err.Error(),
		)
	}

	if len(ships) == 0 || ships[0].Symbol == "" {
		t.Fatalf(
			"%s Got no ships.",
			errPrefix,
		)
	}

	for _, ship := range ships {
		fmt.Printf("%v\n", &ship)
	}

	test_user_first_ship_symbol = ships[0].Symbol
}

func TestGetShip(t *testing.T) {
	errPrefix := "TEST_GetShip():"

	if test_user_first_ship_symbol == "" {
		t.Skipf(
			"%s There is no ship symbol. See TestGetShipsByAgent(). Skipping.",
			errPrefix,
		)
	}

	timerGetShip := timer("GetShip(test_user_first_ship_symbol, TEST_USER_TOKEN")
	ship, err := GetShip(test_user_first_ship_symbol, TEST_USER_TOKEN)
	timerGetShip()
	if err != nil {
		t.Fatalf(
			"%s Getting ship.\n%s",
			errPrefix,
			err.Error(),
		)
	}

	fmt.Printf("%v\n", ship)
}

func TestGetShipNav(t *testing.T) {
	errPrefix := "TEST_GetShipNav():"

	if test_user_first_ship_symbol == "" {
		t.Skipf(
			"%s There is no ship symbol. See TestGetShipsByAgent(). Skipping.",
			errPrefix,
		)
	}

	timerGetShipNav := timer("GetShipNav()")
	shipNav, err := GetShipNav(test_user_first_ship_symbol, TEST_USER_TOKEN)
	timerGetShipNav()
	if err != nil {
		t.Fatalf(
			"%s Getting ship nav.\n%s",
			errPrefix,
			err.Error(),
		)
	}

	if shipNav == nil || shipNav.WaypointSymbol == "" {
		t.Fatalf(
			"%s Bad data.\n%s",
			errPrefix,
			err.Error(),
		)
	}

	timerGetShipLocation := timer("GetShipLocation()")
	shipLocation, err := GetShipLocation(test_user_first_ship_symbol, TEST_USER_TOKEN)
	timerGetShipLocation()
	if err != nil {
		t.Fatalf(
			"%s Getting ship location.\n%s",
			errPrefix,
			err.Error(),
		)
	}

	if shipLocation == (Vector2{}) {
		t.Fatalf(
			"%s Bad data.\n%s",
			errPrefix,
			err.Error(),
		)
	}
}

func TestGetSystemWaypoints(t *testing.T) {
	errPrefix := "TEST_GetSystemWaypoints():"
	var waypoints []Waypoint
	var err error

	timerGetAllWaypointsInSystem := timer("GetAllWaypointsInSystem(\"X1-UQ22\")")
	waypoints, err = GetAllWaypointsInSystem("X1-UQ22")
	timerGetAllWaypointsInSystem()
	if err != nil {
		t.Fatalf(
			"%s Getting all waypoints in system. %s",
			errPrefix,
			err.Error(),
		)
	}

	for idx, waypoint := range waypoints {
		if waypoint.Symbol == (Waypoint{}.Symbol) {
			t.Fatalf(
				"%s Waypoint with index %d has empty symbol.",
				errPrefix,
				idx,
			)
		}
	}

	timerGetSystemWaypoints := timer("GetSystemWaypoints(\"X1-UQ22\", []string{\"STRIPPED\"}, \"ASTEROID\")")
	waypoints, err = GetSystemWaypoints("X1-UQ22", []string{"STRIPPED"}, "ASTEROID")
	timerGetSystemWaypoints()
	if err != nil {
		t.Fatalf(
			"%s Getting all STRIPPED ASTEROIDs in system. %s",
			errPrefix,
			err.Error(),
		)
	}

	for idx, waypoint := range waypoints {
		if waypoint.Symbol == (Waypoint{}.Symbol) {
			t.Fatalf(
				"%s Waypoint with index %d has empty symbol.",
				errPrefix,
				idx,
			)
		}

		if waypoint.Type != "ASTEROID" {
			t.Fatalf(
				"%s Waypoint with index %d is not an ASTEROID.",
				errPrefix,
				idx,
			)
		}

		traitsGood := false
		for _, trait := range waypoint.Traits {
			if trait.Symbol == "STRIPPED" {
				traitsGood = true
			}
		}
		if !traitsGood {
			t.Fatalf(
				"%s Waypoint with index %d has no trait STRIPPED.",
				errPrefix,
				idx,
			)
		}
	}
}

func TestGetMyContracts(t *testing.T) {
	errPrefix := "TEST_GetMyContracts():"
	var contracts []Contract
	var err error

	timerGetContracts := timer("GetMyContracts()")
	contracts, err = GetMyContracts(TEST_USER_TOKEN)
	timerGetContracts()
	if err != nil {
		if errors.Is(err, NoContentError) {
			err = nil
		} else {
			t.Fatalf(
				"%s Error from GetContracts()\n%s",
				errPrefix,
				err.Error(),
			)
		}
	}

	if contracts == nil {
		t.Fatalf(
			"%s GetContracts() returned nil but no error.",
			errPrefix,
		)
	}

	for i, contract := range contracts {
		if contract.ID == "" {
			t.Fatalf(
				"%s Contract with index %d has empty ID.",
				errPrefix,
				i,
			)
		}

		fmt.Printf("Idx %d: ", i)

		if contract.ID == "0" {
			fmt.Println("No contract")
		} else {
			fmt.Printf("%v", contract)
		}
	}
}
