package frappe

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// HTTPClient represents an HTTP client.
type HTTPClient interface {
	Do(method, rURL string, params url.Values, headers http.Header) (HTTPResponse, error)
	DoRaw(method, rURL string, reqBody []byte, headers http.Header) (HTTPResponse, error)
	DoJSON(method, rURL string, params url.Values, headers http.Header, obj interface{}) (HTTPResponse, error)
	GetClient() *httpClient
}

// httpClient is the default implementation of HTTPClient.
type httpClient struct {
	client *http.Client
	hLog   *log.Logger
	debug  bool
}

// HTTPResponse encompasses byte body  + the response of an HTTP request.
type HTTPResponse struct {
	Body     []byte
	Response *http.Response
}

// NewHTTPClient returns a self-contained HTTP request object
// with underlying keep-alive transport.
func NewHTTPClient(h *http.Client, hLog *log.Logger, debug bool) HTTPClient {
	if hLog == nil {
		hLog = log.New(os.Stdout, "base.HTTP: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	if h == nil {
		h = &http.Client{
			Timeout: time.Duration(5) * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   10,
				ResponseHeaderTimeout: time.Second * time.Duration(5),
			},
		}
	}

	return &httpClient{
		hLog:   hLog,
		client: h,
		debug:  debug,
	}
}

// Do executes an HTTP request and returns the response.
func (h *httpClient) DoRaw(method, rURL string, reqBody []byte, headers http.Header) (HTTPResponse, error) {
	var (
		resp     = HTTPResponse{}
		err      error
		postBody io.Reader
	)

	// Encode POST / PUT params.
	if method == http.MethodPost || method == http.MethodPut {
		postBody = bytes.NewReader(reqBody)
	}

	req, err := http.NewRequest(method, rURL, postBody)
	if err != nil {
		h.hLog.Printf("Request preparation failed: %v", err)
		return resp, errors.New("request preparation failed")
	}

	if headers != nil {
		req.Header = headers
	}

	req.Header.Add("Accept", "application/json")

	// If a content-type isn't set, set the default one.
	if req.Header.Get("Content-Type") == "" {
		if method == http.MethodPost || method == http.MethodPut {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	// If the request method is GET or DELETE, add the params as QueryString.
	if method == http.MethodGet || method == http.MethodDelete {
		req.URL.RawQuery = string(reqBody)
	}

	r, err := h.client.Do(req)
	if err != nil {
		h.hLog.Printf("Request failed: %v", err)
		return resp, errors.New("request failed")
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.hLog.Printf("Unable to read response: %v", err)
		return resp, errors.New("error reading response")
	}

	resp.Response = r
	resp.Body = body
	if h.debug {
		h.hLog.Printf("%s %s -- %d %v", method, req.URL.RequestURI(), resp.Response.StatusCode, req.Header)
	}

	return resp, nil
}

// DoRaw executes an HTTP request and returns the response.
func (h *httpClient) Do(method, rURL string, params url.Values, headers http.Header) (HTTPResponse, error) {
	if params == nil {
		params = url.Values{}
	}

	return h.DoRaw(method, rURL, []byte(params.Encode()), headers)
}

// DoJSON makes an HTTP request and parses the JSON response.
func (h *httpClient) DoJSON(method, url string, params url.Values, headers http.Header, obj interface{}) (HTTPResponse, error) {
	resp, err := h.Do(method, url, params, headers)
	if err != nil {
		return resp, err
	}

	// We now unmarshal the body.
	if err := json.Unmarshal(resp.Body, &obj); err != nil {
		h.hLog.Printf("Error parsing JSON response: %v | %s", err, resp.Body)
		return resp, errors.New("error parsing response")
	}

	return resp, nil
}

// GetClient return's the underlying net/http client.
func (h *httpClient) GetClient() *httpClient {
	return h
}
