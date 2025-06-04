import requests
from .config import OPENWEATHER_API_KEY, LATITUDE, LONGITUDE

def get_wind_speed_kmh():
    url = (
        f"https://api.openweathermap.org/data/2.5/weather?"
        f"lat={LATITUDE}&lon={LONGITUDE}&units=metric&appid={OPENWEATHER_API_KEY}"
    )
    response = requests.get(url)
    data = response.json()
    wind_speed_m_s = data["wind"]["speed"]
    return wind_speed_m_s * 3.6  # Convert m/s to km/h
