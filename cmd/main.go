package main

import (
	"fmt"

	logsnag "github.com/uchennaemeruche/go-logsnag"
)

func main() {
	apiToken := "YOUR_TOKEN"
	projectName := "ecommerce-monitor"

	l := logsnag.NewLogsnag(projectName, &logsnag.APIClient{Token: apiToken})

	result, err := l.Publish(
		"waitlist",
		"User Joined",
		logsnag.IPublishPayloadOptions{
			Icon: "🎉",
			Tags: map[string]interface{}{
				"name":  "john doe",
				"email": "johndoe@example.com",
			},
			Description: "A new user has joined the waitlist",
			Notify:      true,
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Sprintln("event published:  %/s", result)

	insights, err := l.Insight("User Count", "100", logsnag.InsightPayloadOptions{Icon: "👨"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Sprintln("Insights:  %/s", insights)
}
