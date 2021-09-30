package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"testing"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

type testScenario struct {
	handler            handlerFunc
	endpoint           string
	expectedBody       string
	expectedStatusCode int
}

// The order is important
var expectedScenarios = []testScenario{
	// Zero value of int so counter should be equal to 0
	{handler: stateHandler, endpoint: "/", expectedBody: "0", expectedStatusCode: http.StatusOK},
	{handler: incrementHandler, endpoint: "/up", expectedBody: "", expectedStatusCode: http.StatusOK},
	{handler: incrementHandler, endpoint: "/up", expectedBody: "", expectedStatusCode: http.StatusOK},
	{handler: incrementHandler, endpoint: "/up", expectedBody: "", expectedStatusCode: http.StatusOK},
	// After three incrementations counter should be equal to 3
	{handler: stateHandler, endpoint: "/", expectedBody: "3", expectedStatusCode: http.StatusOK},
	{handler: decrementHandler, endpoint: "/down", expectedBody: "", expectedStatusCode: http.StatusOK},
	// Now after decrementation counter should be equal to 2
	{handler: stateHandler, endpoint: "/", expectedBody: "2", expectedStatusCode: http.StatusOK},
	// Three another decrementations to test whether counter will stop decrementing on 0
	{handler: decrementHandler, endpoint: "/down", expectedBody: "", expectedStatusCode: http.StatusOK},
	{handler: decrementHandler, endpoint: "/down", expectedBody: "", expectedStatusCode: http.StatusOK},
	{handler: decrementHandler, endpoint: "/down", expectedBody: "", expectedStatusCode: http.StatusOK},
	{handler: stateHandler, endpoint: "/", expectedBody: "0", expectedStatusCode: http.StatusOK},
}

func TestHandlers(t *testing.T) {
	for _, scenario := range expectedScenarios {
		statusCode, body := makeTestAPICall(scenario.handler)

		if statusCode != scenario.expectedStatusCode {
			t.Errorf("'%s' returned wrong status code: got %v want %v",
				getHandlerName(scenario.handler), statusCode, scenario.expectedStatusCode,
			)
		}

		if body != scenario.expectedBody {
			t.Errorf("'%s' returned unexpected body: got %v want %v",
				getHandlerName(scenario.handler), body, scenario.expectedBody,
			)
		}
	}
}

func TestRoutings(t *testing.T) {
	router := newRouter()
	mockedServer := httptest.NewServer(router)

	for _, scenario := range expectedScenarios {
		resp, err := http.Get(mockedServer.URL + scenario.endpoint)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != scenario.expectedStatusCode {
			t.Errorf("Endpoint '%s' returned wrong status code: got %v want %v",
				scenario.endpoint, resp.StatusCode, scenario.expectedStatusCode,
			)
		}

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		body := string(b)
		if body != scenario.expectedBody {
			t.Errorf("Endpoint '%s' returned unexpected body: got %v want %v",
				scenario.endpoint, body, scenario.expectedBody,
			)
		}
	}
}

func makeTestAPICall(handler handlerFunc) (int, string) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	h := http.HandlerFunc(handler)
	h.ServeHTTP(recorder, req)

	statusCode := recorder.Code
	body := recorder.Body.String()

	return statusCode, body
}

func getHandlerName(f handlerFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
