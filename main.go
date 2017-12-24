package main

import (
	"log"
	"net/http"
	"os"
  "encoding/json"
  "io/ioutil"
  "fmt"
  "github.com/pusher/pusher-http-go"
)

type Event struct {
  ChildName string `json:"childName"`
  RelativePoints int `json:"relativePoints"`
}

var pusherClient pusher.Client

var events = []Event{
  Event{ChildName: "sophie", RelativePoints: 1},
  Event{ChildName: "constance", RelativePoints: 1},
}

func getEvents(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Serving events")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  j, err := json.Marshal(events)
  if err != nil {
    http.Error(w, `JSON marshal failure`, 500)
    return
  }
	w.Write(j)
}

func postEvent(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Received new event")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, `Could not read HTTP body`, 500)
    return
  }
  var newEvent Event
  err = json.Unmarshal(body, &newEvent)
  if err != nil {
    http.Error(w, `Could not unmarshal JSON from body`, 400)
    return
  }
  events = append(events, newEvent)

  data := map[string]string{}
  _, err = pusherClient.Trigger("events", "new-event", data)
  if err != nil {
    fmt.Print("Error triggering event:", err.Error())
  } else {
    fmt.Println("Triggered Pusher event")
  }
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    getEvents(w, r)
  } else if r.Method == "POST" {
    postEvent(w, r)
  } else {
    http.Error(w, `Method not allowed`, 405)
  }
}

func handleMessengerWebhook(w http.ResponseWriter, r *http.Request) {
  modes := r.URL.Query()["hub.mode"]
  tokens := r.URL.Query()["hub.verify_token"]
  challenges := r.URL.Query()["hub.challenge"]
  if 0 < len(modes) && string(modes[0]) == "subscribe" && 0 < len(tokens) && string(tokens[0]) == "y64wu657e" {
    fmt.Println("WEBHOOK_VERIFIED")
    w.Write([]byte(challenges[0]))
    return
  }

  fmt.Println("Received messenger webhook")
  bytes, _ := ioutil.ReadAll(r.Body)
  fmt.Println("Body", string(bytes))
}

func main() {
  pusherClient = pusher.Client{
    AppId: os.Getenv("PUSHER_APP_ID"),
    Cluster: os.Getenv("PUSHER_CLUSTER"),
    Key: os.Getenv("PUSHER_KEY"),
    Secret: os.Getenv("PUSHER_SECRET"),
    Secure: true,
  }

	http.HandleFunc("/events", handleEvents)
  http.HandleFunc("/messenger-webhook", handleMessengerWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
