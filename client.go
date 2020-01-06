package kandy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client is the Kandy client struct, holds information regarding with the authenticated party
type Client struct {
	// Kandy CPaaS userId to be used in API calls
	preferredUsername string

	// OAuth and OpenID Connect tokens
	accessToken  string
	idToken      string
	refreshToken string

	// List of the telephone numbers that this project can use for SMS
	smsDidList []string

	// login status: loggedIn, notLoggedIn
	loginStatus string

	// Project credentials
	projectCredentials

	// identifier for this instance of the client, likely be set as UUID
	clientCorrelator string

	// Kandy CPaaS URL
	cpaasURL string
}

// Struct to hold project credentials
type projectCredentials struct {
	privateProjectKey    string
	privateProjectSecret string
}

// Initialize is the function to load key config parameters for the Kandy CPaaS client
func Initialize(url string, key string, secret string, cc string) Client {
	kandyClient := Client{
		cpaasURL: url,
		projectCredentials: projectCredentials{
			privateProjectKey:    key,
			privateProjectSecret: secret,
		},
		clientCorrelator: cc,
	}

	return kandyClient
}

func (c *Client) sendPOSTFormRequest(url string, form url.Values) ([]byte, int, string) {
	resp, err := http.PostForm(url, form)
	if err != nil {
		// Failed, set the status
		return nil, 0, "External error - no response for POST"
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// return error with some explanation
		return nil, resp.StatusCode, "Internal error - POST response body cannot be parsed"
	}

	// We have an HTTP response
	return body, resp.StatusCode, ""
}

func (c *Client) sendPOSTRequest(url string, reqBody []byte) ([]byte, int, string) {
	client := &http.Client{}

	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if reqErr != nil {
		return nil, 0, "Internal error - POST cannot be built"
	}

	// Format Authorization header
	authorizationHeader := "Bearer " + c.accessToken

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorizationHeader)

	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, 0, "External error - no response for POST"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// return error with some explanation
		return nil, resp.StatusCode, "Internal error - POST response body cannot be parsed"
	}

	// We have an HTTP response
	return body, resp.StatusCode, ""
}
