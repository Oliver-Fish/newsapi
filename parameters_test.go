package newsapi

import (
	"errors"
	"testing"
)

// func (p *parameters) buildURL(bURL string, ap *allowedParameters) (string, error) {
func TestBuildURL(t *testing.T) {
	tt := []struct {
		testName          string
		baseURL           string
		parameters        parameters
		allowedParameters allowedParameters
		expectedResult    string
		expectedError     string
	}{
		{
			"Country parameter - Valid",
			"https://newsapi.org/v2/",
			parameters{"country": "gb"},
			allowedParameters{"country": "string"},
			"https://newsapi.org/v2/country=gb",
			"",
		},
		{
			"Country parameter - Invalid",
			"https://newsapi.org/v2/",
			parameters{"country": "Invalid"},
			allowedParameters{"country": "string"},
			"",
			"Unsupported country code Invalid",
		},
		{
			"Category parameter - Valid",
			"https://newsapi.org/v2/",
			parameters{"category": "business"},
			allowedParameters{"category": "string"},
			"https://newsapi.org/v2/category=business",
			"",
		},
		{
			"Category parameter - Invalid",
			"https://newsapi.org/v2/",
			parameters{"category": "Invalid"},
			allowedParameters{"category": "string"},
			"",
			"Invalid category Invalid",
		},
		{
			"language parameter - valid",
			"https://newsapi.org/v2/",
			parameters{"language": "en"},
			allowedParameters{"language": "string"},
			"https://newsapi.org/v2/language=en",
			"",
		},
		{
			"language parameter - Invalid",
			"https://newsapi.org/v2/",
			parameters{"language": "Invalid"},
			allowedParameters{"language": "string"},
			"",
			"Unsupported language Invalid",
		},
		{
			"sortBy parameter - valid",
			"https://newsapi.org/v2/",
			parameters{"sortBy": "publishedAt"},
			allowedParameters{"sortBy": "string"},
			"https://newsapi.org/v2/sortBy=publishedAt",
			"",
		},
		{
			"sortBy parameter - Invalid",
			"https://newsapi.org/v2/",
			parameters{"sortBy": "Invalid"},
			allowedParameters{"sortBy": "string"},
			"",
			"Invalid sort by type Invalid",
		},
		{
			"sources parameter - valid",
			"https://newsapi.org/v2/",
			parameters{"sources": []string{"source1", "source2", "source3"}},
			allowedParameters{"sources": "[]string"},
			"https://newsapi.org/v2/sources=source1%2Csource2%2Csource3",
			"",
		},
		{
			"sources parameter - Invalid over Max 20",
			"https://newsapi.org/v2/",
			parameters{"sources": []string{"source1", "source2", "source3", "source4", "source5", "source6", "source7", "source8", "source9", "source10", "source11", "source12", "source13", "source14", "source15", "source16", "source17", "source18", "source19", "source20", "source21"}},
			allowedParameters{"sources": "[]string"},
			"",
			"Maximum of 20 sources got 21",
		},
		{
			"sources parameter - Invalid 0 sources",
			"https://newsapi.org/v2/",
			parameters{"sources": []string{}},
			allowedParameters{"sources": "[]string"},
			"",
			"Empty list of sources",
		},
		{
			"domains parameter - valid",
			"https://newsapi.org/v2/",
			parameters{"domains": []string{"domain1", "domain2", "domain3"}},
			allowedParameters{"domains": "[]string"},
			"https://newsapi.org/v2/domains=domain1%2Cdomain2%2Cdomain3",
			"",
		},
		{
			"sources parameter - Invalid 0 domains",
			"https://newsapi.org/v2/",
			parameters{"domains": []string{}},
			allowedParameters{"domains": "[]string"},
			"",
			"Empty list of domains",
		},
		{
			"query parameter - valid",
			"https://newsapi.org/v2/",
			parameters{"q": "RandomQuery"},
			allowedParameters{"q": "string"},
			"https://newsapi.org/v2/q=RandomQuery",
			"",
		},
		{
			"query parameter - Invalid",
			"https://newsapi.org/v2/",
			parameters{"q": ""},
			allowedParameters{"q": "string"},
			"",
			"Expected query got empty string",
		},
		{
			"pageSize parameter - valid",
			"https://newsapi.org/v2/",
			parameters{"pageSize": 5},
			allowedParameters{"pageSize": "int"},
			"https://newsapi.org/v2/pageSize=5",
			"",
		},
		{
			"Unhandled parameter",
			"https://newsapi.org/v2/",
			parameters{"Invalid": ""},
			allowedParameters{"Invalid": "string"},
			"",
			"Unhandled parameter Invalid",
		},
		{
			"Failed Verify",
			"https://newsapi.org/v2/",
			parameters{"Invalid": ""},
			allowedParameters{"Invalid": "int"},
			"",
			"Invalid type for parameter Invalid expected type int got type string",
		},
	}

	for _, v := range tt {
		t.Run(v.testName, func(t *testing.T) {
			result, err := v.parameters.buildURL(v.baseURL, &v.allowedParameters)
			if err != nil {
				if v.expectedError != err.Error() {
					t.Fatal(err)
				}
			}
			if v.expectedResult != result {
				t.Fatalf("Expected '%v' got '%v'", v.expectedResult, result)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	tt := []struct {
		testName      string
		allowedParams allowedParameters
		passedParmas  parameters
		returnError   bool
	}{
		{"Invalid Parameters", allowedParameters{"Test1": "string", "Test2": "bool", "Test3": "int"}, parameters{"Test1": "A String", "Test2": true, "Test3": 5, "Test4": "Another String"}, true},
		{"Incorrect Parameter Type", allowedParameters{"Test1": "string", "Test2": "bool", "Test3": "int"}, parameters{"Test1": "A String", "Test2": true, "Test3": "Not an int"}, true},
		{"All parameters valid", allowedParameters{"Test1": "string", "Test2": "bool", "Test3": "int"}, parameters{"Test1": "A String", "Test2": true, "Test3": 5}, false},
	}

	for _, v := range tt {
		t.Run(v.testName, func(t *testing.T) {
			err := v.allowedParams.verify(v.passedParmas)
			if err != nil && !v.returnError {
				t.Fatal(err)
			} else if err == nil && v.returnError {
				t.Fatal(errors.New("Expected verify to return an error but instead it returned nil"))
			}
		})
	}
}

func TestInterfaceToStringList(t *testing.T) {
	tt := []struct {
		testName       string
		data           interface{}
		expectedResult string
		expectedLength int
	}{
		{"Valid string slice", []string{"Test1", "Test2", "Test3", "Test4"}, "Test1,Test2,Test3,Test4", 4},
		{"Empty string slice", []string{}, "", 0},
	}

	for _, v := range tt {
		t.Run(v.testName, func(t *testing.T) {
			r, l := interfaceToStringList(v.data)
			if r != v.expectedResult {
				t.Fatalf("Expected '%v' got '%v'", v.expectedResult, r)
			}

			if l != v.expectedLength {
				t.Fatalf("Expected '%d' got '%d'", v.expectedLength, l)
			}
		})
	}
}
