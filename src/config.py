import os

OPENWEATHER_API_KEY = os.getenv("OPENWEATHER_API_KEY")
DISCORD_WEBHOOK_URL = os.getenv("DISCORD_WEBHOOK_URL")
LATITUDE = float(os.getenv("LATITUDE", "36.0139"))
LONGITUDE = float(os.getenv("LONGITUDE", "-5.606"))
WIND_THRESHOLD_KMH = float(os.getenv("WIND_THRESHOLD_KMH", "10"))
NOTIFY_HOUR = int(os.getenv("NOTIFY_HOUR", "9"))  # 24h format
STATE_FILE = os.getenv("STATE_FILE", "last_notification.txt")
