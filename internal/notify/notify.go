package notify

import (
    "bytes"
    "encoding/json"
    "net/http"

    "github.com/r-dvl/wind-watcher/internal/config"
)

type WebhookMessage struct {
	Content string `json:"content"`
}

func SendDiscordNotification(msg string) error {
    payload := WebhookMessage{Content: msg}
    body, _ := json.Marshal(payload)
    _, err := http.Post(config.GetDiscordWebhookURL(), "application/json", bytes.NewBuffer(body))
    return err
}
