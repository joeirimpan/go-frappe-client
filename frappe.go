package frappe

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const requestTimeout time.Duration = 7000 * time.Millisecond

// Client refers to frappe client
type Client struct {
	baseURI    string
	userName   string
	password   string
	debug      bool
	httpClient HTTPClient
}

// New creates a new frappe client.
func New(baseURI, userName, password string, debug bool) (*Client, error) {
	client := &Client{
		baseURI:  baseURI,
		userName: userName,
		password: password,
		debug:    debug,
	}
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// Create a default http handler with default timeout.
	client.SetHTTPClient(&http.Client{
		Timeout: requestTimeout,
		Jar:     cookieJar,
	})

	// Do login auth which sets the cookies in the jar
	err = client.Login()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Login performs a login request and sets the cookies
func (c *Client) Login() error {
	loginParams := url.Values{}
	loginParams.Set("cmd", "login")
	loginParams.Set("usr", c.userName)
	loginParams.Set("pwd", c.password)
	_, err := c.httpClient.Do(http.MethodPost, c.baseURI, loginParams, nil)
	if err != nil {
		return err
	}

	return nil
}

// SetHTTPClient sets http client for frappe client
func (c *Client) SetHTTPClient(h *http.Client) {
	c.httpClient = NewHTTPClient(h, nil, c.debug)
}

// Do proxy underlying http client do request
func (c *Client) Do(httpMethod, frappeMethod string, params url.Values, headers http.Header) (HTTPResponse, error) {
	return c.httpClient.Do(
		httpMethod,
		c.baseURI+"api/method/"+frappeMethod,
		params,
		headers,
	)
}

// DoJSON proxy underlying http client doJSON request
func (c *Client) DoJSON(httpMethod, frappeMethod string, params url.Values, headers http.Header, obj interface{}) (HTTPResponse, error) {
	return c.httpClient.DoJSON(
		httpMethod,
		c.baseURI+"api/method/"+frappeMethod,
		params,
		headers,
		obj,
	)
}
