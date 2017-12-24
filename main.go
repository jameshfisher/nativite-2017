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
  fmt.Print("Received messenger webhook")
  bytes, _ := ioutil.ReadAll(r.Body)
  fmt.Print("Body", string(bytes))

  VERIFY_TOKEN := "y64wu657e"

  mode := r.URL.Query()["hub.mode"][0]
  token := r.URL.Query()["hub.verify_token"][0]
  challenge := r.URL.Query()["hub.challenge"][0]

  if mode != "" && token != "" {
    if mode == "subscribe" && token == VERIFY_TOKEN {
      fmt.Println("WEBHOOK_VERIFIED")
      w.Write([]byte(challenge))
    } else {
      http.Error(w, `Forbidden`, 403)
    }
  }
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
