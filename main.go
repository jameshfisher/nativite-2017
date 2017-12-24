package main

import (
	"log"
	"net/http"
	"os"
  "encoding/json"
)

type Event struct {
  ChildName string `json:"childName"`
  RelativePoints int `json:"relativePoints"`
}

func events(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  j, err := json.Marshal([]Event{
    Event{ChildName: "sophie", RelativePoints: 1},
    Event{ChildName: "constance", RelativePoints: 1},
  })
  if err != nil {
    http.Error(w, `JSON marshal failure`, 500)
    return
  }
	w.Write(j)
}

func main() {
	http.HandleFunc("/events", events)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
