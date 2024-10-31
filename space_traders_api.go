package space_traders_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	URL_base       *url.URL
	token_GET      *http.Request
	token_POST     *http.Request
	NoContentError = fmt.Errorf("No content from server.")
)

type STJsonError struct {
	Code    int
	Message string
}

func (e *STJsonError) Error() string {
	return fmt.Sprintf(
		"Code: %d\tMessage: %s",
		e.Code,
		e.Message,
	)
}

type SaveData struct {
	Token    string
	Agent    *Agent
	Contract *Contract
	Faction  *any
	Ship     *any
}

func LoadSaveData(filename string) (data *SaveData, err error) {
	const errPrefix = "Loading save data."
	var saveFile *os.File

	saveFile, e := os.Open(filename)
	if e != nil {
		return nil, fmt.Errorf(
			"%s Opening save file %s %w",
			errPrefix,
			filename,
			e,
		)
	}

	fileData, e := io.ReadAll(saveFile)
	if e != nil {
		return nil, fmt.Errorf(
			"%s Reading data from save file %s %w",
			errPrefix,
			filename,
			e,
		)
	}

	e = json.Unmarshal(fileData, data)
	if e != nil {
		return data, fmt.Errorf(
			"%s Decoding JSON from %s %w",
			errPrefix,
			filename,
			e,
		)
	}

	return
}

func init() {
	var err error

	if URL_base == nil {
		URL_base, err = url.Parse("https://api.spacetraders.io")
	}

	if err != nil {
		log.Panicln("STAPI: Init. Failed to parse URL.")
	}
}

func SetBaseURL(rawURL string) (err error) {
	errPrefix := "STAPI: Setting base URL."

	URL_base, err = url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf(
			"%s Parsing URL %w",
			errPrefix,
			err,
		)
	}

	return
}

func LoadToken(token string) (err error) {
	errPrefix := "STAPI: While trying to load token."

	token_GET, err = http.NewRequest(
		"GET",
		URL_base.String(),
		strings.NewReader(""),
	)
	if err != nil || token_GET == nil {
		return fmt.Errorf(
			"%s Creating GET request template with agent token. %w",
			errPrefix,
			err,
		)
	}
	token_GET.Header.Add(
		"Authorization",
		"Bearer "+string(token),
	)

	token_POST, err = http.NewRequest(
		"POST",
		URL_base.String(),
		strings.NewReader(""),
	)
	if err != nil || token_POST == nil {
		return fmt.Errorf(
			"%s Creating POST request template with agent token. %w",
			errPrefix,
			err,
		)
	}
	token_POST.Header.Add(
		"Authorization",
		"Bearer "+string(token),
	)

	return nil
}

// Get a spacetraders.io agent token from some JSON.
// Give this some JSON that follows the pattern:
// { "data": { "token": [TOKEN] } }
// Trailing null bytes cause errors.
// TODO: Fix the trailing bytes thing.
func DecodeToken(JSONdata []byte) ([]byte, error) {
	var jsonMap map[string]any
	error_prefix := "STAPI: While trying to decode JSON and get token."

	err := json.Unmarshal(JSONdata, &jsonMap)
	if err != nil {
		return []byte{}, fmt.Errorf("%s %w", error_prefix, err)
	}

	if data, ok := jsonMap["data"].(map[string]any); !ok {
		return nil, fmt.Errorf(
			"%s data is not a map of strings.",
			error_prefix,
		)
	} else if token, ok := data["token"].(string); !ok {
		return nil, fmt.Errorf(
			"%s data.token is not string.",
			error_prefix,
		)
	} else {
		return []byte(token), err
	}
}

// Executes req with the default client.
// Stores response in a BUFFER_SIZE buffer, then trims null bytes.
// I've been getting some weird results with larger items, even when they don't come close to the buffer size.
//
//	I'm not sure if this is a spacetraders.io thing or a problem with my approach.
func SendRequest(req *http.Request) (response []byte, err error) {
	const BUFFER_SIZE = 1024

	errPrefix := "Sending GET"

	if req == nil {
		req, err = http.NewRequest(
			"GET",
			"https://api.spacetraders.io/v2",
			io.NopCloser(strings.NewReader("")),
		)
		if req != nil {
			defer req.Body.Close()
		}
		if err != nil {
			return nil, fmt.Errorf(
				"%s\n\tCreating request.%w",
				errPrefix,
				err,
			)
		}
	}

	errPrefix += fmt.Sprintf(
		"\tres %s",
		req.URL.Path,
	)

	if !req.Close {
		req.Close = true
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"%s\n\tExecuting request.%w",
			errPrefix,
			err,
		)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return nil, fmt.Errorf(
			"%s Non-200 status from request.\n\t%s",
			errPrefix,
			resp.Status,
		)
	}

	for err != io.EOF {
		buf := []byte(nil)
		buf = make([]byte, BUFFER_SIZE)

		_, err = resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return response, fmt.Errorf(
				"%s\nReading response.\n%w",
				errPrefix,
				err,
			)
		}
		buf = bytes.TrimRight(buf, "\x00")

		response = append(response, buf...)

		if len(response) > 1024*1024 {
			return response, fmt.Errorf(
				"%s Response buffer has become very large.",
				errPrefix,
			)
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}
