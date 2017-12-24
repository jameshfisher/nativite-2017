package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func scores(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, `{"sophie": 12, "constance": 12, "victoire": 12, "felicite": 12}`)
}

func main() {
	http.HandleFunc("/scores", scores)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
