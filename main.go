package main

import (
	"log"
	"net/http"
	"os"
  "encoding/json"
  "io/ioutil"
  "fmt"
)

type Event struct {
  ChildName string `json:"childName"`
  RelativePoints int `json:"relativePoints"`
}

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

func main() {
	http.HandleFunc("/events", handleEvents)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
