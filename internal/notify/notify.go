package notify

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/r-dvl/wind-watcher/internal/config"
    "github.com/r-dvl/wind-watcher/internal/wind"
)

func SendWeatherNotification(
    location string,
    threshold float64,
    days []wind.DayWindInfo,
    bestDay *wind.DayWindInfo,
    mapURL string,
) error {
    msg := formatWeatherMessage(location, threshold, days, bestDay)
    webhook := config.GetDiscordWebhookURL()
    if webhook != "" {
        return SendDiscordWeatherNotification(msg, bestDay, mapURL, days)
    }
    fmt.Fprintln(os.Stdout, msg)
    return nil
}

func formatWeatherMessage(location string, threshold float64, days []wind.DayWindInfo, bestDay *wind.DayWindInfo) string {
    msg := fmt.Sprintf("🌬️ Wind forecast for %s (≤ %.1f km/h):\n", location, threshold)
    for _, day := range days {
        t, err := time.Parse("2006-01-02", day.Date)
        weekday := day.Date
        if err == nil {
            weekday = t.Weekday().String()
        }
        mainData, _ := day.Data["main"].(map[string]interface{})
        windData, _ := day.Data["wind"].(map[string]interface{})
        temp := "?"
        humidity := "?"
        windDir := "?"
        if mainData != nil {
            if v, ok := mainData["temp"].(float64); ok {
                temp = fmt.Sprintf("%.1f°C", v)
            }
            if v, ok := mainData["humidity"].(float64); ok {
                humidity = fmt.Sprintf("%d%%", int(v))
            }
        }
        if windData != nil {
            if v, ok := windData["deg"].(float64); ok {
                windDir = wind.WindDirLabel(v)
            }
        }
        best := ""
        if bestDay != nil && day.Date == bestDay.Date {
            best = " ⭐ Best day!"
        }
        msg += fmt.Sprintf(
            "- %s (%s): wind %.1f km/h, %s, %s, %s%s\n",
            day.Date, weekday, day.Speed, temp, humidity, windDir, best,
        )
    }
    return msg
}

func SendDiscordWeatherNotification(msg string, bestDay *wind.DayWindInfo, mapURL string, days []wind.DayWindInfo) error {
    var fields []EmbedField

    for _, day := range days {
        t, err := time.Parse("2006-01-02", day.Date)
        weekday := day.Date
        if err == nil {
            weekday = t.Weekday().String()
        }
        mainData, _ := day.Data["main"].(map[string]interface{})
        windData, _ := day.Data["wind"].(map[string]interface{})
        temp := "?"
        humidity := "?"
        windDir := "?"
        if mainData != nil {
            if v, ok := mainData["temp"].(float64); ok {
                temp = fmt.Sprintf("%.1f°C", v)
            }
            if v, ok := mainData["humidity"].(float64); ok {
                humidity = fmt.Sprintf("%d%%", int(v))
            }
        }
        if windData != nil {
            if v, ok := windData["deg"].(float64); ok {
                windDir = wind.WindDirLabel(v)
            }
        }
        best := ""
        if bestDay != nil && day.Date == bestDay.Date {
            best = "⭐"
        }
        fields = append(fields, EmbedField{
            Name:   fmt.Sprintf("%s (%s) %s", weekday, day.Date, best),
            Value:  fmt.Sprintf("**Wind:** %.1f km/h\n**Temp:** %s\n**Humidity:** %s\n**Dir:** %s", day.Speed, temp, humidity, windDir),
            Inline: true,
        })
    }

    embed := Embed{
        Title:       "🌬️ Wind forecast",
        Description: fmt.Sprintf("Wind forecast for the next days:\nThreshold: ≤ %.1f km/h", bestDay.Speed),
        Color:       3447003,
        URL:         mapURL,
        Fields:      fields,
    }

    payload := WebhookMessage{
        Content:   "",
        Embeds:    []Embed{embed},
        Username:  "Wind Watcher",
        AvatarURL: "https://openweathermap.org/themes/openweathermap/assets/img/logo_white_cropped.png",
    }
    body, _ := json.Marshal(payload)
    _, err := http.Post(config.GetDiscordWebhookURL(), "application/json", bytes.NewBuffer(body))
    return err
}

// --- Discord embed types ---

type EmbedField struct {
    Name   string `json:"name"`
    Value  string `json:"value"`
    Inline bool   `json:"inline"`
}

type EmbedThumbnail struct {
    URL string `json:"url"`
}

type Embed struct {
    Title       string          `json:"title,omitempty"`
    Description string          `json:"description,omitempty"`
    Color       int             `json:"color,omitempty"`
    URL         string          `json:"url,omitempty"`
    Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
    Fields      []EmbedField    `json:"fields,omitempty"`
}

type WebhookMessage struct {
    Content   string  `json:"content,omitempty"`
    Embeds    []Embed `json:"embeds,omitempty"`
    Username  string  `json:"username,omitempty"`
    AvatarURL string  `json:"avatar_url,omitempty"`
}