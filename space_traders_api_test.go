package space_traders_api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// method is a valid method for http.NewRequest()
// Caller should close the .Body of the generate request.
func generateTestUserTemplate(method string) (*http.Request, error) {
	errPrefix := "generateTestUserTemplate():"
	//TEST_USER
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGlmaWVyIjoiVEVTVF9VU0VSIiwidmVyc2lvbiI6InYyLjIuMCIsInJlc2V0X2RhdGUiOiIyMDI0LTEwLTA2IiwiaWF0IjoxNzI5MTk3NzY1LCJzdWIiOiJhZ2VudC10b2tlbiJ9.ahkvlMr8P0PPsEFPmi8KURTwj1Cn3Xa3Mq5LENBGWTIqovhhqsYOBMN_NJH1LcTQWsjx5zTBQeUi0h0mIPlNpiPULj7OioVPc1HvJRk6x2U_P8UMpW7bI9CDOZ6E2nKyUCpqsZb7qzq6zClC5cZGcqJA56Y9dFI0zI6qyKqN1IxBzG-_4jKKpoXvGsHibbJVHexiImZpCeo-ORuF531luYIdXfGcNCccPYVL5Drqq8sAIEFieapGpYnWYMeS7Vk4dgLOXyybnRtEbKeMERXE19cwHmeF01UJY4jCf0Vd6obOs5OP0Qn3jwtii-GoldFHVL3Z6AROe2S7WEaCYKuG2w"

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

func describe(i interface{}) string {
	return fmt.Sprintf("(%v, %T)", i, i)
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

	result, resultStatus, err := GetAgentDetails(req)
	if err != nil {
		t.Fatalf(
			"%s First call returned error. %s\nObject: %v\nStatus: %s",
			errPrefix,
			err.Error(),
			result,
			resultStatus,
		)
	}

	if result.Symbol != "TEST_USER" ||
		result.AccountID != "cm2drpb4ua6gts60c84n8hy0a" {
		t.Fatalf(
			"%s Result of first call is not TEST_USER. User is %s\naccountID is %s",
			errPrefix,
			result.Symbol,
			result.AccountID,
		)
	}

	result = Agent{}

	//Do a second call using the same template to make sure we don't mutate the template or something.
	result, resultStatus, err = GetAgentDetails(req)
	if err != nil {
		t.Fatalf(
			"%s Second call returned error. %s\nObject: %v\nStatus: %s",
			errPrefix,
			err.Error(),
			result,
			resultStatus,
		)
	}

	if result.Symbol != "TEST_USER" ||
		result.AccountID != "cm2drpb4ua6gts60c84n8hy0a" {
		t.Fatalf(
			"%s Result of first call is not TEST_USER. User is %s\naccountID is %s",
			errPrefix,
			result.Symbol,
			result.AccountID,
		)
	}

	t.Logf(
		"%s\n%v",
		errPrefix,
		result,
	)
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

	if contracts == nil {
		t.Fatalf(
			"%s GetContracts() returned nil but no error.",
			errPrefix,
		)
	}
}
