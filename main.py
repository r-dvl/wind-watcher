import time
import datetime
from src.config import NOTIFY_HOUR, WIND_THRESHOLD_KMH
from src.wind import get_wind_speed_kmh
from src.notify import send_discord_notification
from src.state import has_notified_today, mark_notified_today

def main():
    print("✅ Wind Notifier started")
    while True:
        now = datetime.datetime.now()
        print(f"[{now}] Checking wind...")

        if now.hour == NOTIFY_HOUR and not has_notified_today():
            try:
                wind_speed = get_wind_speed_kmh()
                print(f"[{now}] Wind speed: {wind_speed:.1f} km/h")

                if wind_speed < WIND_THRESHOLD_KMH:
                    message = f"🌬️ Low wind in Tarifa: {wind_speed:.1f} km/h. Perfect day!"
                    send_discord_notification(message)
                    print(f"[{now}] Notification sent.")
                    mark_notified_today()
                else:
                    print(f"[{now}] Wind too strong, no notification sent.")
            except Exception as e:
                print(f"[{now}] ❌ Error checking wind: {e}")
        else:
            print(f"[{now}] Not time to notify or already notified today.")

        time.sleep(3600)


if __name__ == "__main__":
    main()
