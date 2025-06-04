package wind

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/r-dvl/wind-watcher/internal/config"
)

type WindData struct {
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

func GetWindSpeed() (float64, error) {
    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/weather?lat=36.0131&lon=-5.6078&appid=%s",
        config.GetOpenWeatherAPIKey(),
    )
    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var data WindData
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        return 0, err
    }

    return data.Wind.Speed * config.GetWindKMHFactor(), nil
}
