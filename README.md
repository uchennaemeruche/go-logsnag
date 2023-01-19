# Logsnag's Golang implementation

Get notifications and track your project events. Please refer to the [docs](https://docs.logsnag.com/) to learn more.


## Installation

```golang
go get github.com/uchennaemeruche/go-logsnag
```

## Usage

### Import Package
```golang
    logsnag "github.com/uchennaemeruche/go-logsnag"
```

### Initialize Client
```golang
    l := logsnag.NewLogsnag(projectName, &logsnag.APIClient{Token: apiToken})
```

### Publish Event
```golang
   	result, err := l.Publish(
		"user-create",
		"Uchenna created a new account",
		logsnag.IPublishPayloadOptions{
			Icon: "🎉",
			Tags: map[string]interface{}{
				"name":  "john doe",
				"email": "johndoe@example.com",
			},
			Description: "A new user has been created",
			Notify:      true,
		},
	)
```

### Create Insight
```golang
   insights, err := l.Insight("user-create", "New Account", logsnag.InsightPayloadOptions{Icon: "👨"})
```


## Example Implementation

```golang
package main

import (
	"fmt"
	logsnag "github.com/uchennaemeruche/go-logsnag"
)

func main() {
	apiToken := "<YOUR_API_TOKEN>"
	projectName := "<your_project_name_>"

	l := logsnag.NewLogsnag(projectName, &logsnag.APIClient{Token: apiToken})

    // Publish an event
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

    // Create insight
	insights, err := l.Insight("User Count", "100", logsnag.InsightPayloadOptions{Icon: "👨"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Sprintln("Insights:  %/s", insights)

	fmt.Println(resp)
}
```
**NOTE**
Check the `cmd` directory to see the sample implementation code.
