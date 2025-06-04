import requests
from .config import DISCORD_WEBHOOK_URL

def send_discord_notification(message: str):
    payload = {"content": message}
    response = requests.post(DISCORD_WEBHOOK_URL, json=payload)
    response.raise_for_status()
