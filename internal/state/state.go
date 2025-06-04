package state

import (
    "os"
    "strings"
    "time"

    "github.com/r-dvl/wind-watcher/internal/config"
)

func HasNotifiedToday() bool {
    data, err := os.ReadFile(config.GetStateFile())
    if err != nil {
        // create file if missing
        _ = os.WriteFile(config.GetStateFile(), []byte(""), 0644)
        return false
    }
    last := strings.TrimSpace(string(data))
    today := time.Now().Format("2006-01-02")
    return last == today
}

func MarkNotifiedToday() {
    today := time.Now().Format("2006-01-02")
    _ = os.WriteFile(config.GetStateFile(), []byte(today), 0644)
}
