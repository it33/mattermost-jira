package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Data struct for JIRA JSON parsing.
type Data struct {
	WebhookEvent string
	User         struct {
		Name        string
		AvatarUrls  map[string]string
		DisplayName string
	}
	Issue struct {
		Self   string
		Key    string
		Fields struct {
			Issuetype struct {
				IconURL string
				Name    string
			}
			Summary string
		}
	}
	Comment struct {
		Body string
	}
	Changelog struct {
		Items []struct {
			Field      string
			FromString string
			ToString   string
		}
	}
}

// Message structure for Mattermost JSON creation.
type Message struct {
	Text     string `json:"text"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
}

// RoomMapping map holds ProjectKey:MattermostRoomName
type RoomMapping map[string]string

func getMessage(request *http.Request) []byte {
	// Parse JSON from JIRA
	decoder := json.NewDecoder(request.Body)
	var data Data
	decoder.Decode(&data)

	// Get JIRA URL from "issue" section in JSON
	u, _ := url.Parse(data.Issue.Self)

	// Select action
	var action, comment string
	switch data.WebhookEvent {
	case "jira:issue_created":
		action = "created"
	case "jira:issue_updated":
		action = "updated"
	case "jira:issue_deleted":
		action = "deleted"
	}

	//Process new comment
	if len(data.Comment.Body) > 0 {
		comment = fmt.Sprintf("\nComment:\n```\n%s\n```", data.Comment.Body)
	}

	// Process changelog
	var changelog string
	if len(data.Changelog.Items) > 0 {
		for _, item := range data.Changelog.Items {
			itemName := strings.ToUpper(string(item.Field[0])) + item.Field[1:]
			if item.FromString == "" {
				item.FromString = "None"
			}
			if itemName == "Description" {
				changelog += fmt.Sprintf(
					"\nNew Description:\n```\n%s\n```",
					item.ToString,
				)
			} else {
				changelog += fmt.Sprintf(
					"\n%s: ~~%s~~ %s",
					itemName,
					item.FromString,
					item.ToString,
				)
			}
		}
	}

	// Create message for Mattermost
	text := fmt.Sprintf(
		//Message format:
		//![user_icon](user_icon_link)[UserFirstName UserSecondName](user_link) commented task ![task_icon](task_icon link)[TSK-42](issue_link) "Test task"
		//Status: ~~Done~~ Finished
		//>Comment text
		"![user_icon](%s) [%s](%s://%s/secure/ViewProfile.jspa?name=%s) %s %s ![task_icon](%s) [%s](%s://%s/browse/%s) \"%s\"%s%s",
		data.User.AvatarUrls["16x16"],
		data.User.DisplayName,
		u.Scheme,
		u.Host,
		data.User.Name,
		action,
		strings.ToLower(data.Issue.Fields.Issuetype.Name),
		data.Issue.Fields.Issuetype.IconURL,
		data.Issue.Key,
		u.Scheme,
		u.Host,
		data.Issue.Key,
		data.Issue.Fields.Summary,
		changelog,
		comment,
	)

	channel := ""
	issueKeySplit := strings.Split(data.Issue.Key, "-")
	fmt.Println("Mappings", mappings)
	fmt.Println("Length", len(issueKeySplit))
	fmt.Println("First", issueKeySplit[0])
	if mappings != nil && len(issueKeySplit) > 1 && issueKeySplit[0] != "" {
		channel, _ = mappings[issueKeySplit[0]]
		fmt.Printf("Room Lookup: %q\n", channel)
	}

	message := Message{
		Text:     text,
		Channel:  channel,
		Username: "JIRA",
		IconURL:  "https://raw.githubusercontent.com/csduarte/mattermost-jira/master/logo-02.png",
	}
	fmt.Printf("Input Data: \n %v", data)
	fmt.Printf("Output Message: \n %v", message)

	JSONMessage, _ := json.Marshal(message)

	return JSONMessage
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request Incoming: %v%v\n", r.Host, r.RequestURI)
	// Get mattermost URL
	mattermostHookURL := r.URL.Query().Get("mattermost_hook_url")

	if len(mattermostHookURL) > 0 {
		// Get message from JIRA JSON request
		message := getMessage(r)

		// Create http-client
		req, _ := http.NewRequest("POST", mattermostHookURL, bytes.NewBuffer(message))
		req.Header.Set("Content-Type", "application/json")

		// Send data to Mattermost
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("WARN: Request missing query param 'mattermost_hook_url`")
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Write([]byte("OK"))
}

var mappings RoomMapping

func main() {
	var mapfile = flag.String("map", "", "JSON Mapping File")
	flag.Parse()
	fmt.Println(*mapfile)

	if *mapfile != "" {
		mapdata, err := ioutil.ReadFile(*mapfile)
		if err != nil {
			fmt.Printf("Failed to load map file from path: %q\n", *mapfile)
			fmt.Printf("Error: %v\n", err.Error())
			os.Exit(1)
		}
		if err := json.Unmarshal(mapdata, &mappings); err != nil {
			fmt.Printf("Failed to decode JSON from mapping: %v\n", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("WARN: No Channel Map file specified (-map=/path/to/file)")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Printf("Server starting on 0.0.0.0:%v\n", port)
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
