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
            windSpeed, err := wind.GetWindSpeed()
            if err != nil {
                log.Println("❌ Error getting wind speed:", err)
            } else {
                log.Printf("🌬️ Current wind speed: %.1f km/h\n", windSpeed)
                if windSpeed < float64(config.GetWindThreshold()) {
                    msg := fmt.Sprintf("🌬️ Low wind in Tarifa: %.1f km/h. Perfect day!", windSpeed)
                    err = notify.SendDiscordNotification(msg)
                    if err != nil {
                        log.Println("❌ Error sending notification:", err)
                    } else {
                        log.Println("✅ Notification sent")
                        state.MarkNotifiedToday()
                    }
                } else {
                    log.Println("💨 Wind too strong, skipping notification")
                }
            }
        } else {
            log.Println("🕒 Not time or already notified today")
        }

		time.Sleep(1 * time.Hour)
	}
}
