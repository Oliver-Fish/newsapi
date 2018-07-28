package newsapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiPath           = "https://newsapi.org/v2" //Base path for API
	apiHeadlinePath   = "/top-headlines?"
	apiSourcePath     = "/sources?"
	apiEverythingPath = "/everything?"
)

type errorResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Client struct {
	APIUrl     string
	APIKey     string
	HTTPClient *http.Client
}

type Option func(*Client)

type Article struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
}
type ArticleResults struct {
	Status       string    `json:"status"`
	TotalResults int64     `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SourceResults struct {
	Status string   `json:"status"`
	Source []Source `json:"sources"`
}

//New creates our Client struct that we use to make requests
func New(apiKey string, options ...Option) *Client {
	c := Client{
		APIUrl: apiPath,
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
	for _, o := range options {
		o(&c) //Apply any caller options to the returned Client struct
	}
	return &c
}

//makeRequest handles all our requests to newsapi
func (c *Client) makeRequest(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	if c.APIKey == "" {
		return nil, errors.New("Expected API key got nothing")
	}
	req.Header.Add("X-Api-Key", c.APIKey)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		var errResp errorResponse
		err := json.Unmarshal(b, &errResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.Message)
	}
	return b, nil
}

func (c *Client) GetTopHeadlines(p parameters) (ArticleResults, error) {
	var o ArticleResults
	ap := allowedParameters{ //List of allowed parameters and their allowed types
		"country":  "string",
		"category": "string",
		"sources":  "[]string",
		"q":        "string",
		"pageSize": "int",
		"page":     "int"}
	u, err := p.buildURL(c.APIUrl+apiHeadlinePath, &ap)
	if err != nil {
		return o, err
	}

	d, err := c.makeRequest(u)
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(d, &o)
	if err != nil {
		return o, nil
	}
	return o, nil
}

func (c *Client) GetEverything(p parameters) (ArticleResults, error) {
	var o ArticleResults
	ap := allowedParameters{ //List of allowed parameters and their allowed types
		"q":        "string",
		"sources":  "[]string",
		"domains":  "[]string",
		"from":     "string",
		"to":       "string",
		"lanague":  "string",
		"sortBy":   "string",
		"pageSize": "int",
		"page":     "int"}
	u, err := p.buildURL(c.APIUrl+apiEverythingPath, &ap)
	if err != nil {
		return o, err
	}
	d, err := c.makeRequest(u)
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(d, &o)
	if err != nil {
		return o, nil
	}
	return o, nil
}

func (c *Client) GetSources(p parameters) (SourceResults, error) {
	var o SourceResults
	ap := allowedParameters{ //List of allowed parameters and their allowed types
		"country":  "string",
		"category": "string",
		"lanague":  "string"}
	u, err := p.buildURL(c.APIUrl+apiSourcePath, &ap)
	if err != nil {
		return o, err
	}
	d, err := c.makeRequest(u)
	if err != nil {
		return o, err
	}
	err = json.Unmarshal(d, &o)
	if err != nil {
		return o, nil
	}
	return o, nil
}
