import datetime
from .config import STATE_FILE

def has_notified_today() -> bool:
    try:
        with open(STATE_FILE, "r") as f:
            last_date = f.read().strip()
            return last_date == datetime.date.today().isoformat()
    except FileNotFoundError:
        with open(STATE_FILE, "w") as f:
            f.write("")
        return False

def mark_notified_today():
    with open(STATE_FILE, "w") as f:
        f.write(datetime.date.today().isoformat())
