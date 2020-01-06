package kandy

import (
	"encoding/json"
	"strconv"
	"strings"
)

type outboundSMSMessageRequestPayload struct {
	OutboundSMSMessageRequest outboundSMSMessageRequest `json:"outboundSMSMessageRequest"`
}

type outboundSMSMessageRequest struct {
	Address                []string               `json:"address"`
	ClientCorrelator       string                 `json:"clientCorrelator"`
	OutboundSMSTextMessage outboundSMSTextMessage `json:"outboundSMSTextMessage"`
}

type outboundSMSTextMessage struct {
	Message string `json:"message"`
}

type outboundSMSMessageResponseSuccesPayload struct {
	OutboundSMSMessageRequestSuccess outboundSMSMessageRequestSuccess `json:"outboundSMSMessageRequest"`
}

type outboundSMSMessageRequestSuccess struct {
	// Only interested in resourceURL
	ResourceURL string `json:"resourceURL"`
}

// SendSMS sends SMS via the provided sender and destination addresses, returns resourceURL to the SMS message in history when SMS sent successfully
func (c *Client) SendSMS(sender string, destination string, text string) (string, error) {
	a := []string{destination}
	payload := outboundSMSMessageRequestPayload{
		OutboundSMSMessageRequest: outboundSMSMessageRequest{
			Address:          a,
			ClientCorrelator: c.clientCorrelator,
			OutboundSMSTextMessage: outboundSMSTextMessage{
				Message: text,
			},
		},
	}

	var jsonData []byte
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", &kandyError{0, "Internal error - cannot parse JSON body"}
	}

	outboundSMSUrl := c.cpaasURL + "/cpaas/smsmessaging/v1/" + c.preferredUsername + "/outbound/" + sender + "/requests"
	respBody, statusCode, errorText := c.sendPOSTRequest(outboundSMSUrl, jsonData)

	if statusCode == 0 {
		// An internal error occured
		return "", &kandyError{0, errorText}
	}

	if statusCode == 201 {
		// Message sent successfully, extract messageId from the body
		var successResp outboundSMSMessageResponseSuccesPayload
		err := json.Unmarshal(respBody, &successResp)
		if err != nil {
			// This should not happen
			return "", &kandyError{201, "Response body cannot be parsed"}
		}

		resourceURL := successResp.OutboundSMSMessageRequestSuccess.ResourceURL
		s := strings.Split(resourceURL, "/")

		// ResourceURL format: /cpaas/smsmessaging/v1/{userId}/remoteAddresses/{remoteAddress}/localAddresses/{localAddress}/messages/{messageId}
		return s[10], nil
	}

	if statusCode == 401 {
		// Unauthorized, sth wrong about the access token
		return "", &kandyError{401, "Authentication failed, please relogin"}
	}

	if statusCode == 400 || statusCode > 401 {
		// Let's parse the error and return
		var failureResp requestErrorPayload
		err := json.Unmarshal(respBody, &failureResp)
		if err != nil {
			// This should not happen
			return "", &kandyError{statusCode, "Response body cannot be parsed"}
		}

		serviceExp := failureResp.RequestError.ServiceException
		policyExp := failureResp.RequestError.PolicyException
		var errText string
		var variables []string
		if serviceExp.MessageID != "" {
			errText = serviceExp.Text
			variables = serviceExp.Variables
		} else {
			errText = policyExp.Text
			variables = policyExp.Variables
		}

		if variables != nil {
			for i, item := range variables {
				errText = strings.Replace(errText, "%"+strconv.Itoa(i+1), item, -1)
			}
		}

		return "", &kandyError{statusCode, errText}
	}

	// We should not hit here so just in case..
	return "", &kandyError{statusCode, "An error occurred"}
}
