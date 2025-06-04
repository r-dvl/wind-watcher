package config

import (
	"log"
	"os"
	"strconv"
)

func GetStateFile() string {
	return getEnv("STATE_FILE", "/tmp/last_notification_date")
}

func GetWindThreshold() int {
	return getEnvAsInt("WIND_THRESHOLD", 11)
}

func GetOpenWeatherAPIKey() string {
	return mustGetEnv("OPENWEATHER_API_KEY")
}

func GetDiscordWebhookURL() string {
	return mustGetEnv("DISCORD_WEBHOOK_URL")
}

func GetWindKMHFactor() float64 {
	return getEnvAsFloat("WIND_KMH_FACTOR", 3.6)
}

func GetNotifyHour() int {
	return getEnvAsInt("NOTIFY_HOUR", 9)
}

// Helpers

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return val
}

func getEnvAsInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Invalid integer value for %s: %s", key, valStr)
	}
	return val
}

func getEnvAsFloat(key string, fallback float64) float64 {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		log.Fatalf("Invalid float value for %s: %s", key, valStr)
	}
	return val
}
