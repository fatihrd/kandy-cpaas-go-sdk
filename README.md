# kandy-cpaas-go-sdk

## Install
`go get github.com/fatihrd/kandy-cpaas-go-sdk`

## Supported Features
* Login via project credentials
* SendSMS

## Example Usage
```go
package main

import (
    "fmt"
    "github.com/fatihrd/kandy-cpaas-go-sdk"
)

func main() {
	fmt.Println("Starting Kandy CPaaS")

	kandyClient := kandy.Initialize(
		"https://oauth-cpaas.att.com",
		"projectKey",
		"projectSecret",
		"exampleClientcorrelator",
	)

	err := kandyClient.Login()
	fmt.Println(err)

	messageID, err2 := kandyClient.SendSMS("+1xxxxxxxxx", "+1xxxxxxxxx", "Test")
	fmt.Println(messageID, err2)
}
```