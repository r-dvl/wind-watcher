package wind

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/r-dvl/wind-watcher/internal/config"
)

// Returns both the full weather data (as map[string]interface{}) and the wind speed (float64)
func GetWeatherDataAndWindSpeed() (map[string]interface{}, float64, error) {
    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/weather?lat=36.0131&lon=-5.6078&appid=%s&units=metric",
        config.GetOpenWeatherAPIKey(),
    )
    resp, err := http.Get(url)
    if err != nil {
        return nil, 0, err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        return nil, 0, err
    }

    wind, ok := data["wind"].(map[string]interface{})
    if !ok {
        return data, 0, fmt.Errorf("wind data missing")
    }
    speed, ok := wind["speed"].(float64)
    if !ok {
        return data, 0, fmt.Errorf("wind speed missing")
    }

    // Convert to km/h if needed
    speed = speed * config.GetWindKMHFactor()

    return data, speed, nil
}