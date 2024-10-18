package space_traders_api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

func describe(i interface{}) string {
	return fmt.Sprintf("(%v, %T)", i, i)
}

func TestGetAgentDetails(t *testing.T) {
	errPrefix := "TEST_GetAgentDetails:"

	//TEST_USER
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGlmaWVyIjoiVEVTVF9VU0VSIiwidmVyc2lvbiI6InYyLjIuMCIsInJlc2V0X2RhdGUiOiIyMDI0LTEwLTA2IiwiaWF0IjoxNzI5MTk3NzY1LCJzdWIiOiJhZ2VudC10b2tlbiJ9.ahkvlMr8P0PPsEFPmi8KURTwj1Cn3Xa3Mq5LENBGWTIqovhhqsYOBMN_NJH1LcTQWsjx5zTBQeUi0h0mIPlNpiPULj7OioVPc1HvJRk6x2U_P8UMpW7bI9CDOZ6E2nKyUCpqsZb7qzq6zClC5cZGcqJA56Y9dFI0zI6qyKqN1IxBzG-_4jKKpoXvGsHibbJVHexiImZpCeo-ORuF531luYIdXfGcNCccPYVL5Drqq8sAIEFieapGpYnWYMeS7Vk4dgLOXyybnRtEbKeMERXE19cwHmeF01UJY4jCf0Vd6obOs5OP0Qn3jwtii-GoldFHVL3Z6AROe2S7WEaCYKuG2w"

	req, err := http.NewRequest(
		"GET",
		"https://api.spacetraders.io",
		io.NopCloser(strings.NewReader("")),
	)
	if err != nil {
		t.Fatalf(
			"%s Error creating template request. %s",
			errPrefix,
			err.Error(),
		)
	}
	defer req.Body.Close()

	req.Header.Add("Authorization", "Bearer "+token)

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

	log.Printf(
		"Got Agent %v",
		result,
	)

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

	log.Printf(
		"Got Agent %v",
		result,
	)

	if result.Symbol != "TEST_USER" ||
		result.AccountID != "cm2drpb4ua6gts60c84n8hy0a" {
		t.Fatalf(
			"%s Result of first call is not TEST_USER. User is %s\naccountID is %s",
			errPrefix,
			result.Symbol,
			result.AccountID,
		)
	}
}
