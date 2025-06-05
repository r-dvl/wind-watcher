package notify

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/r-dvl/wind-watcher/internal/config"
)

type EmbedField struct {
    Name   string `json:"name"`
    Value  string `json:"value"`
    Inline bool   `json:"inline"`
}

type EmbedThumbnail struct {
    URL string `json:"url"`
}

type Embed struct {
    Title       string         `json:"title,omitempty"`
    Description string         `json:"description,omitempty"`
    Color       int            `json:"color,omitempty"`
    URL         string         `json:"url,omitempty"`
    Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
    Fields      []EmbedField   `json:"fields,omitempty"`
}

type WebhookMessage struct {
    Content   string  `json:"content,omitempty"`
    Embeds    []Embed `json:"embeds,omitempty"`
    Username  string  `json:"username,omitempty"`
    AvatarURL string  `json:"avatar_url,omitempty"`
}

func windType(deg float64) string {
    switch {
    case deg >= 45 && deg <= 135:
        return "Levante (East wind)"
    case deg >= 225 && deg <= 315:
        return "Poniente (West wind)"
    default:
        return "Other"
    }
}

func SendDiscordWeatherNotification(msg string, weatherData map[string]interface{}) error {
    main, _ := weatherData["main"].(map[string]interface{})
    wind, _ := weatherData["wind"].(map[string]interface{})
    weatherArr, _ := weatherData["weather"].([]interface{})
    var weather map[string]interface{}
    if len(weatherArr) > 0 {
        weather, _ = weatherArr[0].(map[string]interface{})
    }
    icon, _ := weather["icon"].(string)

    // Try to get coordinates from forecast "city" field if not present
    var lat, lon interface{}
    if coord, ok := weatherData["coord"].(map[string]interface{}); ok {
        lat = coord["lat"]
        lon = coord["lon"]
    } else if city, ok := weatherData["city"].(map[string]interface{}); ok {
        lat = city["coord"].(map[string]interface{})["lat"]
        lon = city["coord"].(map[string]interface{})["lon"]
    }

    deg, _ := wind["deg"].(float64)
    windKind := windType(deg)

    embed := Embed{
        Title:       "🌬️ Wind Alert!",
        Description: msg,
        Color:       3447003,
        URL:         fmt.Sprintf("https://openweathermap.org/weathermap?basemap=map&cities=true&layer=wind&lat=%v&lon=%v&zoom=10", lat, lon),
        Thumbnail:   &EmbedThumbnail{URL: fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", icon)},
        Fields: []EmbedField{
            {
                Name:   "Temperature",
                Value:  fmt.Sprintf("%.1f°C", main["temp"].(float64)),
                Inline: true,
            },
            {
                Name:   "Humidity",
                Value:  fmt.Sprintf("%d%%", int(main["humidity"].(float64))),
                Inline: true,
            },
            {
                Name:   "Wind",
                Value:  fmt.Sprintf("%.1f m/s, %.0f° (%s)", wind["speed"].(float64), deg, windKind),
                Inline: true,
            },
        },
    }

    payload := WebhookMessage{
        Content:   msg,
        Embeds:    []Embed{embed},
        Username:  "Wind Watcher",
        AvatarURL: fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", icon),
    }
    body, _ := json.Marshal(payload)
    _, err := http.Post(config.GetDiscordWebhookURL(), "application/json", bytes.NewBuffer(body))
    return err
}