package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/irussellps/alarmserver/config"
)

type Bus struct {
	Debug    bool
	urls     []string
	telegram bool
	client   *http.Client
}

type WebhookPayload struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

func (webhooks *Bus) Initialize(config config.WebhooksConfig) {
	fmt.Println("Initializing Webhook bus...")
	webhooks.client = &http.Client{}
	webhooks.urls = config.Urls
	webhooks.telegram = config.Telegram
}

func (webhooks *Bus) SendMessage(topic string, data string) {
	for _, url := range webhooks.urls {
		if webhooks.telegram {
			//combinevalues := string(topic) + string(data)
			//print("Topic: " + string(topic) + "\n")
			//print("Data: " + string(data) + "\n")
			topic = strings.ToUpper(topic)
			data = strings.ToUpper(data)
			tstr := strings.ReplaceAll(topic, "/", " - ")
			dstr := strings.ReplaceAll(data, "/", " - ")
			//replacing spaces with url encoding
			//combinedurloutput := strings.ReplaceAll(combinevalues, " ", "%20") // replacing space with %20
			//printing entire
			//print("url&combinedoutput: " + url + tstr + dstr + "\n\n\n")
			url = url + tstr + " - " + dstr
			//print(url + "\n\n\n")
			payload := WebhookPayload{Topic: topic, Data: data}
			go webhooks.send(url, payload)
		} else {
			//print("Whooopsy, no telegram set")
			payload := WebhookPayload{Topic: topic, Data: data}
			go webhooks.send(url, payload)
		}

	}
}

func (webhooks *Bus) send(url string, payload WebhookPayload) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("WEBHOOKS: Error marshaling payload to JSON", err)
		return
	}
	response, err := webhooks.client.Post(url, "application/json", bytes.NewBuffer(payloadJson))
	if err != nil {
		fmt.Printf("WEBHOOKS: Error delivering payload %s to %s\n", payloadJson, url)
	}
	if response == nil {
		fmt.Printf("WEBHOOKS: Got no response from %s", url)
	}
	if response != nil && response.StatusCode != 200 {
		fmt.Printf("WEBHOOKS: Got bad status code delivering payload to %s: %v\n", url, response.StatusCode)
	}
}
