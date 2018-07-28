package newsapi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New("TestAPIKey")
	ct := reflect.TypeOf(c).String()
	if ct != "*newsapi.Client" {
		t.Fatalf("Expected type of %v got %v", "*newsapi.Client", ct)
	}
}

func TestMakeRequest(t *testing.T) {
	sData, err := ioutil.ReadFile("testdata/everything_sucess.json")
	if err != nil {
		t.Fatal(err)
	}
	fData, err := ioutil.ReadFile("testdata/everything_failure.json")
	if err != nil {
		t.Fatal(err)
	}

	testServer := fakeServer(sData, fData)
	defer testServer.Close()

	tt := []struct {
		testName      string
		apiKey        string
		endPoint      string
		expectedData  string
		expectedError bool
	}{
		{
			"Error from API",
			"NotValid",
			testServer.URL + "/failure/",
			string(fData),
			true},
		{
			"Valid Response",
			"NotValid",
			testServer.URL,
			string(sData),
			false},
		{
			"Invalid Endpoint URL",
			"NotValid",
			"https://NotReal",
			"",
			true},
		{
			"Missing API Key",
			"",
			testServer.URL,
			"",
			true},
	}

	for _, v := range tt {
		t.Run(v.testName, func(t *testing.T) {
			c := New(v.apiKey)
			c.HTTPClient.Timeout = time.Second * 2
			d, err := c.makeRequest(v.endPoint)
			if err != nil {
				if !v.expectedError {
					t.Fatalf("Unexpected error - %v", err.Error())
				}
				return
			}
			if v.expectedData != string(d) {
				t.Fatal("Passed data doesn't match returned data")
			}
		})
	}
}

func fakeServer(sData, fData []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "/failure/") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(fData)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(sData)
		}
	}))
}

func TestGetTopHeadlines(t *testing.T) {
	sData, err := ioutil.ReadFile("testdata/topheadlines_sucess.json")
	if err != nil {
		t.Fatal(err)
	}
	fData, err := ioutil.ReadFile("testdata/topheadlines_failure.json")
	if err != nil {
		t.Fatal(err)
	}
	testServer := fakeServer(sData, fData)
	defer testServer.Close()

	c := New("invalid")
	c.APIUrl = testServer.URL
	tt := []struct {
		testName      string
		parameters    parameters
		expectedData  []byte
		expectedError bool
	}{
		{"Valid", parameters{}, sData, false},
		{"Error Returned from API", parameters{}, fData, true},
		{"Invalid BuildURL Options", parameters{"invalid": "string"}, fData, true},
	}
	for _, v := range tt {
		t.Run(v.testName, func(t *testing.T) {
			if v.expectedError {
				c.APIUrl = testServer.URL + "/failure/"
			}
			a, err := c.GetTopHeadlines(v.parameters)
			if err != nil {
				if !(v.expectedError) {
					t.Fatalf("Unexpected Error '%v'", err.Error())
				}
				return
			}
			if v.expectedError {
				t.Fatal("Expected error got nil")
			}

			if reflect.TypeOf(a).String() != "newsapi.ArticleResults" {
				t.Fatalf("Expected type of newsapi.ArticleResults got %v", reflect.TypeOf(a).String())
			}

		})
	}
}

func TestGetEverything(t *testing.T) {
	sData, err := ioutil.ReadFile("testdata/everything_sucess.json")
	if err != nil {
		t.Fatal(err)
	}
	fData, err := ioutil.ReadFile("testdata/everything_failure.json")
	if err != nil {
		t.Fatal(err)
	}
	testServer := fakeServer(sData, fData)
	defer testServer.Close()

	c := New("invalid")
	c.APIUrl = testServer.URL
	tt := []struct {
		testName      string
		parameters    parameters
		expectedData  []byte
		expectedError bool
	}{
		{"Valid", parameters{}, sData, false},
		{"Error Returned from API", parameters{}, fData, true},
		{"Invalid BuildURL Options", parameters{"invalid": "string"}, fData, true},
	}
	for _, v := range tt {
		t.Run(v.testName, func(t *testing.T) {
			if v.expectedError {
				c.APIUrl = testServer.URL + "/failure/"
			}
			a, err := c.GetEverything(v.parameters)
			if err != nil {
				if !(v.expectedError) {
					t.Fatalf("Unexpected Error '%v'", err.Error())
				}
				return
			}
			if v.expectedError {
				t.Fatal("Expected error got nil")
			}

			if reflect.TypeOf(a).String() != "newsapi.ArticleResults" {
				t.Fatalf("Expected type of newsapi.ArticleResults got %v", reflect.TypeOf(a).String())
			}

		})
	}
}

func TestGetSources(t *testing.T) {
	sData, err := ioutil.ReadFile("testdata/sources_sucess.json")
	if err != nil {
		t.Fatal(err)
	}
	fData, err := ioutil.ReadFile("testdata/sources_failure.json")
	if err != nil {
		t.Fatal(err)
	}
	testServer := fakeServer(sData, fData)
	defer testServer.Close()

	c := New("invalid")
	c.APIUrl = testServer.URL
	tt := []struct {
		testName      string
		parameters    parameters
		expectedData  []byte
		expectedError bool
	}{
		{"Valid", parameters{}, sData, false},
		{"Error Returned from API", parameters{}, fData, true},
		{"Invalid BuildURL Options", parameters{"invalid": "string"}, fData, true},
	}
	for _, v := range tt {
		t.Run(v.testName, func(t *testing.T) {
			if v.expectedError {
				c.APIUrl = testServer.URL + "/failure/"
			}
			a, err := c.GetSources(v.parameters)
			if err != nil {
				if !(v.expectedError) {
					t.Fatalf("Unexpected Error '%v'", err.Error())
				}
				return
			}
			if v.expectedError {
				t.Fatal("Expected error got nil")
			}

			if reflect.TypeOf(a).String() != "newsapi.SourceResults" {
				t.Fatalf("Expected type of newsapi.SourceResults got %v", reflect.TypeOf(a).String())
			}

		})
	}
}
