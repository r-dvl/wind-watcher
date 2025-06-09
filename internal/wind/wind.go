package wind

import (
    "encoding/json"
    "fmt"
    "math"
    "net/http"
    "strings"
    "time"

    "github.com/r-dvl/wind-watcher/internal/config"
)

type DayWindInfo struct {
    Date  string
    Speed float64
    Data  map[string]interface{}
}

type WeekWindInfo struct {
    Days    []DayWindInfo
    BestDay *DayWindInfo
}

func WindDirLabel(deg float64) string {
    toDeg := math.Mod(deg+180, 360)
    loc := strings.ToLower(config.GetLocation())

    if strings.Contains(loc, "tarifa") {
        switch {
        case toDeg >= 67.5 && toDeg < 112.5:
            return "Levante (→)"
        case toDeg >= 247.5 && toDeg < 292.5:
            return "Poniente (←)"
        default:
            return "Wind shift"
        }
    }

    switch {
    case toDeg >= 22.5 && toDeg < 67.5:
        return "↗"
    case toDeg >= 67.5 && toDeg < 112.5:
        return "→"
    case toDeg >= 112.5 && toDeg < 157.5:
        return "↘"
    case toDeg >= 157.5 && toDeg < 202.5:
        return "↓"
    case toDeg >= 202.5 && toDeg < 247.5:
        return "↙"
    case toDeg >= 247.5 && toDeg < 292.5:
        return "←"
    case toDeg >= 292.5 && toDeg < 337.5:
        return "↖"
    default:
        return "↑"
    }
}

func GetWeeklyWindPrediction(threshold float64) (*WeekWindInfo, error) {
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

    var qualifyingDays []DayWindInfo
    var bestDay *DayWindInfo
    for i := 0; i < 7; i++ {
        d := time.Now().AddDate(0, 0, i).Format("2006-01-02")
        if info, ok := windByDay[d]; ok && info.Speed <= threshold {
            qualifyingDays = append(qualifyingDays, info)
            if bestDay == nil || info.Speed < bestDay.Speed {
                tmp := info
                bestDay = &tmp
            }
        }
    }

    if len(qualifyingDays) == 0 {
        return nil, fmt.Errorf("no days with wind below or equal to threshold this week")
    }

    return &WeekWindInfo{
        Days:    qualifyingDays,
        BestDay: bestDay,
    }, nil
}