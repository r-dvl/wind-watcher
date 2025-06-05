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
            weatherData, windSpeed, err := wind.GetWeatherDataAndWindSpeed()
            if err != nil {
                log.Println("❌ Error getting wind data:", err)
            } else {
                msg := fmt.Sprintf("🌬️ Wind in Tarifa: %.1f km/h.", windSpeed)
                err = notify.SendDiscordWeatherNotification(msg, weatherData)
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