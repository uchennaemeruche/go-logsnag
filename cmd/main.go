package main

import (
	"fmt"

	logsnag "github.com/uchennaemeruche/go-logsnag"
)

func main() {
	apiKey := "YOUR_TOKEN"
	projectName := "bank-api"

	// l := logsnag.NewLogsnag(&logsnag.NewLogParams{Token: apiKey, Project: projectName})
	l := logsnag.NewLogsnag(projectName, &logsnag.APIClient{Token: apiKey})
	resp, err := l.Publish(
		"user-create",
		"Uchenna created a new account",
		logsnag.IPublishPayloadOptions{
			Description: "A new user has been created",
			Notify:      true,
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Sprintln("event published:  %/s", resp)

	insights, err := l.Insight("user-create", "New Account", logsnag.InsightPayloadOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Sprintln("Insights:  %/s", insights)
}
