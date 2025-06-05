# Wind Watcher

**Wind Watcher** is a simple Go application that notifies you via Discord when the wind forecast in your chosen location is favorable (i.e., below a configurable threshold) for today or the next two days. It is ideal for outdoor enthusiasts, kitesurfers, windsurfers, or anyone who needs to plan around the wind.

## Features

- Fetches wind forecasts from OpenWeatherMap for your configured location.
- Checks if the wind speed is below or equal to your specified threshold for today or the next two days.
- Sends a detailed notification to a Discord channel via webhook, including temperature, humidity, wind speed, and direction.
- Configurable via environment variables or Kubernetes ConfigMap.
- Prevents duplicate notifications within the same day.

## Configuration

Set the following environment variables (or use the provided Helm chart/config files):

| Variable                | Description                                 | Default           |
|-------------------------|---------------------------------------------|-------------------|
| `OPENWEATHER_API_KEY`   | Your OpenWeatherMap API key                 | **(required)**    |
| `DISCORD_WEBHOOK_URL`   | Discord webhook URL for notifications       | **(required)**    |
| `LOCATION`              | Location for weather forecast (e.g. Tarifa,ES) | Tarifa,ES      |
| `WIND_THRESHOLD`        | Wind speed threshold (km/h)                 | 11                |
| `NOTIFY_HOUR`           | Hour of the day to check and notify (0-23)  | 9                 |
| `STATE_FILE`            | Path to file for notification state         | /tmp/last_notification_date |
| `WIND_KMH_FACTOR`       | Conversion factor to km/h (default for m/s) | 3.6               |

## Usage

### Docker

```sh
docker build -t wind-watcher .
docker run -e OPENWEATHER_API_KEY=your_key -e DISCORD_WEBHOOK_URL=your_webhook_url wind-watcher
```

### Docker Compose

See [`docker-compose.yml`](docker-compose.yml) for an example.

### Kubernetes

A Helm chart is provided in the [`helm/wind-watcher`](helm/wind-watcher) directory.

```sh
cd helm/wind-watcher
helm install wind-watcher .
```

## How it works

1. At the configured hour, the app fetches the wind forecast for your location.
2. It checks today and the next two days for any period where the wind is below or equal to your threshold.
3. If a suitable day is found, it sends a notification to your Discord channel.
4. Only one notification is sent per day.

## License

This project is licensed under the [GNU GPL v3](LICENSE).

---

Made with ❤️ by r-dvl, life in Tarifa is hard.