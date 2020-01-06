package kandy

import "fmt"

type requestErrorPayload struct {
	RequestError requestError `json:"requestError"`
}

type requestError struct {
	ServiceException serviceException `json:"serviceException"`
	PolicyException  policyException  `json:"policyException"`
}

type serviceException struct {
	MessageID string   `json:"messageId"`
	Text      string   `json:"text"`
	Variables []string `json:"variables"`
}

type policyException struct {
	MessageID string   `json:"messageId"`
	Text      string   `json:"text"`
	Variables []string `json:"variables"`
}

type kandyError struct {
	code int
	text string
}

func (k *kandyError) Error() string {
	return fmt.Sprintf("%d - %s", k.code, k.text)
}
