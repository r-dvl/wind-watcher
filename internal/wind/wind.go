package wind

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/r-dvl/wind-watcher/internal/config"
)

type DayWindInfo struct {
    Date  string
    Speed float64
    Data  map[string]interface{}
}

func GetForecastWindBelowThreshold(threshold float64) (*DayWindInfo, error) {
    url := fmt.Sprintf(
        "https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric",
        config.GetLocation(),
        config.GetOpenWeatherAPIKey(),
    )
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        return nil, err
    }

    list, ok := data["list"].([]interface{})
    if !ok {
        return nil, fmt.Errorf("forecast data missing")
    }

    // Group by day, find the minimum wind speed for each day
    windByDay := make(map[string]DayWindInfo)
    for _, entry := range list {
        item, ok := entry.(map[string]interface{})
        if !ok {
            continue
        }
        dtTxt, _ := item["dt_txt"].(string)
        day := ""
        if len(dtTxt) >= 10 {
            day = dtTxt[:10]
        }
        wind, ok := item["wind"].(map[string]interface{})
        if !ok {
            continue
        }
        speed, ok := wind["speed"].(float64)
        if !ok {
            continue
        }
        speed = speed * config.GetWindKMHFactor()
        if v, exists := windByDay[day]; !exists || speed < v.Speed {
            windByDay[day] = DayWindInfo{
                Date:  day,
                Speed: speed,
                Data:  item,
            }
        }
    }

    // Check today and the next two days
    for i := 0; i <= 2; i++ {
        d := time.Now().AddDate(0, 0, i).Format("2006-01-02")
        if info, ok := windByDay[d]; ok && info.Speed <= threshold {
            return &info, nil
        }
    }

    return nil, fmt.Errorf("no day with wind below or equal to threshold")
}