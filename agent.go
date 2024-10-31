package space_traders_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

type Agent struct {
	AccountID       string //`json:"accountID"`
	Credits         int    //`json:"credits"`
	Headquarters    string //`json:"headquarters"`
	ShipCount       int    //`json:"shipCount"`
	StartingFaction string //`json:"startingFaction"`
	Symbol          string //`json:"symbol"`
}

func (self Agent) String() string {
	return fmt.Sprintf(
		"Agent %s\n"+
			"\tAccount ID:\t%s\n"+
			"\tCredits:\t%d\n"+
			"\tHeadquarters\t%s\n"+
			"\tShip Count:\t%d\n"+
			"\t(Starting) Faction:\t%s\n",
		self.Symbol,
		self.AccountID,
		self.Credits,
		self.Headquarters,
		self.ShipCount,
		self.StartingFaction,
	)
}


// Creates a spacetraders.io agent.
// Writes the (JSON) response body to a file: savefiles/[agent].json
// Returns the entire response body (including the token)
func CreateAgent(agent string, faction string) (SaveData, error) {
	const BUFFER_SIZE = 20000
	errPrefix := "Trying to create new agent."
	responseBody := make([]byte, BUFFER_SIZE)
	agentMap := make(map[string]string)
	agentMap["symbol"] = agent
	agentMap["faction"] = faction

	jsonObject := new(struct {
		Data *SaveData
	})

	agent_json_bytes, err := json.Marshal(agentMap)
	if err != nil {
		return SaveData{}, fmt.Errorf("%s Marshalling JSON. %w",
			errPrefix,
			err,
		)
	}

	resp, err := http.Post(
		URL_base.String()+"/v2/register",
		"application/json",
		strings.NewReader(string(agent_json_bytes)),
	)
	if err != nil {
		return SaveData{}, fmt.Errorf(
			"%s Trying to POST new agent request.\n%s\n%w",
			errPrefix,
			string(agent_json_bytes),
			err,
		)
	}
	defer resp.Body.Close()

	bodySize, err := resp.Body.Read(responseBody)
	if err != nil {
		return SaveData{}, fmt.Errorf(
			"%s Trying to read body of new agent request. %w",
			errPrefix,
			err,
		)
	}
	responseBody = bytes.TrimRight(responseBody, "\x00")
	responseBody = append(responseBody, byte('\n'))

	err = os.WriteFile("savefiles/"+agent+".json", responseBody, fs.ModePerm)
	if err != nil {
		return SaveData{}, fmt.Errorf(
			"%s Trying to write new agent file. %w",
			errPrefix,
			err,
		)
	}

	if bodySize >= BUFFER_SIZE {
		return SaveData{}, fmt.Errorf(
			"%s Response body too big for buffer (%d bytes).",
			errPrefix,
			BUFFER_SIZE,
		)
	}

	bytes.TrimRight(responseBody, "\x00")

	err = json.Unmarshal(responseBody, jsonObject)
	if err != nil {
		return SaveData{}, fmt.Errorf(
			"%s Unmarshaling JSON.\n%w",
			errPrefix,
			err,
		)
	}

	if jsonObject.Data == nil || jsonObject.Data.Agent == nil {
		return SaveData{}, fmt.Errorf(
			"%s Problem with SaveData, something is nil.\n",
			errPrefix,
		)
	}

	if *jsonObject.Data == (SaveData{}) || *jsonObject.Data.Agent == (Agent{}) {
		return *jsonObject.Data, fmt.Errorf(
			"%s Problem with SaveData, something is empty.\n",
			errPrefix,
		)
	}

	return *jsonObject.Data, err
}

// TODO: refactor for new request pattern.
//
//	Give it an http.Request with the GET method and "Authorization: Bearer [token] header already added.
//	or nil to use token_GET
//
// Returns Agent and raw JSON from request.
func GetAgentDetails(requestTemplate *http.Request) (Agent, string, error) {
	const BUFFER_SIZE = 10000
	error_prefix := "STAPI: Trying to get agent details."
	responseBody := make([]byte, BUFFER_SIZE)
	var JSONobject map[string]*Agent

	if requestTemplate == nil {
		if token_GET == nil {
			return Agent{}, "", fmt.Errorf(
				"%s token_GET is nil",
				error_prefix,
			)
		}
		requestTemplate = token_GET
	}

	req := requestTemplate.Clone(requestTemplate.Context())
	req.Body = io.NopCloser(strings.NewReader(""))
	req.URL.Path = "/v2/my/agent"
	defer req.Body.Close()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Agent{}, resp.Status, fmt.Errorf(
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
		return Agent{}, string(responseBody), fmt.Errorf(
			"%s Trying to read response body. %w",
			error_prefix,
			err,
		)
	}
	if bodySize >= BUFFER_SIZE {
		return Agent{}, string(responseBody), fmt.Errorf(
			"%s Response body too big for buffer (%d bytes).",
			error_prefix,
			BUFFER_SIZE,
		)
	}

	err = json.Unmarshal(responseBody, &JSONobject)
	if err != nil {
		return Agent{}, string(responseBody), fmt.Errorf(
			"%s Unmarshalling JSON. %w",
			error_prefix,
			err,
		)
	}

	if JSONobject["data"] == nil {
		return Agent{}, string(responseBody), fmt.Errorf(
			"%s No data.\n%w",
			error_prefix,
			err,
		)
	}

	return *JSONobject["data"], string(responseBody), err
}
