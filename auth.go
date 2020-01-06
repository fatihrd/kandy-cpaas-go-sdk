package kandy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type tokenResponsePayload struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

type idtokenClaims struct {
	PreferredUsername string `json:"preferred_username"`
	SmsDidList        string `json:"sms-did-list"`
}

type authErrorResponse struct {
	Message string `json:"message"`
}

// Login is the API function to authenticate and get Kandy CPaaS tokens
func (c *Client) Login() error {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", c.projectCredentials.privateProjectKey)
	form.Add("client_secret", c.projectCredentials.privateProjectSecret)
	form.Add("scope", "openid")

	tokenURL := c.cpaasURL + "/cpaas/auth/v1/token"

	body, status, err := c.sendPOSTFormRequest(tokenURL, form)

	if status == 200 {
		// Success, load the token values
		tokenResponse := fetchTokens(body)
		c.accessToken = tokenResponse.AccessToken
		c.refreshToken = tokenResponse.RefreshToken
		c.idToken = tokenResponse.IDToken

		result := setClaims(c)
		if result != "" {
			// preferred_user claim cannot be fetched, we cannot use the API
			c.loginStatus = "LoggedOut"
			return &kandyError{status, result}
		}

		c.loginStatus = "LoggedIn"
		return nil
	}

	// We should have a failure here
	var errorText string
	if err != "" {
		errorText = err
	} else {
		// Get the error text from the body
		var authError authErrorResponse
		err := json.Unmarshal(body, &authError)
		if err != nil {
			// Cannot parse response
			return &kandyError{status, "Response body cannot be parsed"}
		}
		errorText = authError.Message
	}

	return &kandyError{status, errorText}
}

func fetchTokens(b []byte) (t tokenResponsePayload) {
	var tokenResponse tokenResponsePayload
	err := json.Unmarshal(b, &tokenResponse)
	if err != nil {
		// too bad :(
		fmt.Println("Cannot decode tokens")
	}
	return tokenResponse
}

func setClaims(c *Client) string {

	payload := strings.Split(c.idToken, ".")

	// JWT contains 3 parts
	if len(payload) != 3 {
		return "id_token is not valid!"
	}

	var RawStdEncoding = base64.StdEncoding.WithPadding(base64.NoPadding)
	data, err := RawStdEncoding.DecodeString(payload[1])
	if err != nil {
		return "id_token cannot be decoded!"
	}

	var claims idtokenClaims
	err = json.Unmarshal(data, &claims)
	if err != nil {
		return "id_token cannot be parsed"
	}

	c.preferredUsername = claims.PreferredUsername
	c.smsDidList = strings.Split(claims.SmsDidList, ",")

	return ""
}
