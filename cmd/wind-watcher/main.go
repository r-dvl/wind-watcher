package main

import (
    "fmt"
    "log"
    "time"

    "github.com/r-dvl/wind-watcher/internal/config"
    "github.com/r-dvl/wind-watcher/internal/notify"
    "github.com/r-dvl/wind-watcher/internal/state"
    "github.com/r-dvl/wind-watcher/internal/wind"
)

func main() {
    fmt.Println("✅ Wind Notifier started")

    for {
        now := time.Now()

        if now.Hour() >= config.GetNotifyHour() && !state.HasNotifiedToday() {
            threshold := float64(config.GetWindThreshold())
            weekWind, err := wind.GetWeeklyWindPrediction(threshold)
            if err != nil {
                log.Println("❌ No notification sent:", err)
            } else {
                location := config.GetLocation()
                mapURL := fmt.Sprintf("https://openweathermap.org/find?q=%s", location)
                err = notify.SendWeatherNotification(
                    location,
                    threshold,
                    weekWind.Days,
                    weekWind.BestDay,
                    mapURL,
                )
                if err != nil {
                    log.Println("❌ Error sending notification:", err)
                } else {
                    log.Println("✅ Notification sent")
                    state.MarkNotifiedToday()
                }
            }
        } else {
            log.Println("🕒 Not time or already notified today")
        }

        time.Sleep(1 * time.Hour)
    }
}