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
		result.AccountID != TEST_USER_ACCOUNT_ID{
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
		result.AccountID != TEST_USER_ACCOUNT_ID{
		t.Fatalf(
			"%s Result of second call is not TEST_USER. User is %s\naccountID is %s",
			errPrefix,
			result.Symbol,
			result.AccountID,
		)
	}
}

func TestGetSystemWaypoints(t *testing.T) {
	errPrefix := "TEST_GetSystemWaypoints():"
	var waypoints []Waypoint
	var err error

	timerGetAllWaypointsInSystem := timer("GetAllWaypointsInSystem(\"X1-UQ22\")")

	waypoints, err = GetAllWaypointsInSystem("X1-UQ22")
	if err != nil {
		t.Fatalf(
			"%s Getting all waypoints in system. %s",
			errPrefix,
			err.Error(),
		)
	}

	timerGetAllWaypointsInSystem()

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
	if err != nil {
		t.Fatalf(
			"%s Getting all STRIPPED ASTEROIDs in system. %s",
			errPrefix,
			err.Error(),
		)
	}

	timerGetSystemWaypoints()

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

func TestGetContracts(t *testing.T) {
	errPrefix := "TEST_GetContracts():"
	var contracts []Contract
	var err error

	req, err := generateTestUserTemplate("GET")
	if err != nil {
		t.Fatalf(
			"%s Error creating template request. %s",
			errPrefix,
			err.Error(),
		)
	}
	defer req.Body.Close()

	timerGetContracts := timer("GetContracts()")

	if contracts, err = GetContracts(req); err != nil {
		if errors.Is(err, NoContentError) {
			err = nil
		} else {
			t.Fatalf(
				"%s Error from GetContracts() %s",
				errPrefix,
				err.Error(),
			)
		}
	}

	timerGetContracts()

	if contracts == nil {
		t.Fatalf(
			"%s GetContracts() returned nil but no error.",
			errPrefix,
		)
	}
}
