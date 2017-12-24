package main

import (
	"log"
	"net/http"
	"os"
  "encoding/json"
  "io/ioutil"
  "fmt"
  "github.com/pusher/pusher-http-go"
  "strings"
)

type Event struct {
  ChildName string `json:"childName"`
  RelativePoints int `json:"relativePoints"`
}

type MessengerMessage struct {
 Text string `json:"text"`
}

type MessengerRecipient struct {
  Id string `json:"id"`
}

type MessengerRequestBody struct {
  MessagingType string `json:"messaging_type"`
  Recipient MessengerRecipient `json:"recipient"`
  Message MessengerMessage `json:"message"`
}

var pusherClient pusher.Client

var events = []Event{}

var realNames = map[string]string{
  "sophie": "Sophie",
  "constance": "Constance",
  "victoire": "Victoire",
  "felicite": "Félicité",
  "james": "James",
}

var messengerRecipients = map[string]struct{}{}

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

func messageText(doerNom string, bonne bool) string {
  if bonne {
    return realNames[doerNom] + " a fait une bonne chose! Dépêchez-vous, elle pourrait gagner le prix mystère!"
  } else {
    return realNames[doerNom] + " a fait une mauvaise chose! :O"
  }
}

func sendMessengerMsg(recipientId string, msgText string) error {
    messengerReqBodyBytes, err := json.Marshal(MessengerRequestBody{
      MessagingType: "UPDATE",
      Recipient: MessengerRecipient{
        Id: recipientId,
      },
      Message: MessengerMessage{
        Text: msgText,
      },
    })
    if err != nil {
      return err
    }
    resp, err := http.Post(
      "https://graph.facebook.com/v2.6/me/messages?access_token=" + os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN"),
      "application/json",
      strings.NewReader(string(messengerReqBodyBytes)),
    )
    if err != nil {
      return err
    }
    messengerRespBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return err
    }
    fmt.Println("Sent message; body: " + string(messengerRespBody))
    return nil
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
    http.Error(w, `Could not trigger Pusher event`, 500)
    return
  }

  msgText := messageText(newEvent.ChildName, newEvent.RelativePoints > 0)

  for recipientId, _ := range messengerRecipients {
    err := sendMessengerMsg(recipientId, msgText)
    if err != nil {
      fmt.Println("Could not send Messenger message: " + err.Error())
      http.Error(w, `Could not send Messenger message`, 500)
      return
    }
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

type MessengerWebhookBody struct {
  Entries []MessengerWebhookEntry `json:"entry"`
}
type MessengerWebhookEntry struct {
  Messagings []MessengerWebhookMessagings `json:"messaging"`
}
type MessengerWebhookMessagings struct {
  Sender MessengerWebhookSender `json:"sender"`
}
type MessengerWebhookSender struct {
  Id string `json:"string"`
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
  bytes, err := ioutil.ReadAll(r.Body)
  fmt.Println("Body", string(bytes))

  var messengerWebhookBody MessengerWebhookBody
  err = json.Unmarshal(bytes, &messengerWebhookBody)
  if err != nil {
    fmt.Println("Could not unmarshal body")
    return
  }

  senderId := messengerWebhookBody.Entries[0].Messagings[0].Sender.Id

  _, alreadyAdded := messengerRecipients[senderId]
  if !alreadyAdded {
    fmt.Println("TODO send a welcome message")
  }
  messengerRecipients[senderId] = struct{}{}
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
