package frappe

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const requestTimeout time.Duration = 7000 * time.Millisecond

// Client refers to frappe client
type Client struct {
	baseURI    string
	auth       Auth
	authHeader *string
	debug      bool
	httpClient HTTPClient
}

// Auth generic auth interface
type Auth interface{}

// LoginAuth performs normal login flow
type LoginAuth struct {
	UserName string
	Password string
}

// BasicAuth sends base64 encoded auth header
type BasicAuth struct {
	APIKey    string
	APISecret string
}

// TokenAuth sends token formed from apiKey and apiSecret
type TokenAuth struct {
	APIKey    string
	APISecret string
}

// New creates a new frappe client.
func New(baseURI string, auth Auth, debug bool) (*Client, error) {
	client := &Client{
		baseURI: baseURI,
		auth:    auth,
		debug:   debug,
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

	switch a := auth.(type) {
	case *LoginAuth:
		// Do login auth which sets the cookies in the jar
		err = client.Login()
		if err != nil {
			return nil, err
		}
	case *BasicAuth:
		tk := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.APIKey, a.APISecret)))
		basicTk := fmt.Sprintf("Basic %s", tk)
		client.authHeader = &basicTk
	case *TokenAuth:
		tk := fmt.Sprintf("token %s:%s", a.APIKey, a.APISecret)
		client.authHeader = &tk
	}

	return client, nil
}

// Login performs a login request and sets the cookies
func (c *Client) Login() error {
	var (
		auth        = c.auth.(*LoginAuth)
		loginParams = url.Values{}
	)

	loginParams.Set("cmd", "login")
	loginParams.Set("usr", auth.UserName)
	loginParams.Set("pwd", auth.Password)
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
	// Set custom headers
	if c.authHeader != nil {
		if headers == nil {
			headers = make(http.Header)
		}
		headers.Set("Authorization", *c.authHeader)
	}

	return c.httpClient.Do(
		httpMethod,
		c.baseURI+"api/method/"+frappeMethod,
		params,
		headers,
	)
}

// DoJSON proxy underlying http client doJSON request
func (c *Client) DoJSON(httpMethod, frappeMethod string, params url.Values, headers http.Header, obj interface{}) (HTTPResponse, error) {
	// Set custom headers
	if c.authHeader != nil {
		if headers == nil {
			headers = make(http.Header)
		}
		headers.Set("Authorization", *c.authHeader)
	}

	return c.httpClient.DoJSON(
		httpMethod,
		c.baseURI+"api/method/"+frappeMethod,
		params,
		headers,
		obj,
	)
}
