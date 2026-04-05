from abc import ABC, abstractmethod
from dataclasses import dataclass
from datetime import datetime, timedelta
from enum import Enum
from typing import Dict, List, Tuple

import pandas as pd


class Level(Enum):
    INFO = "Info"
    WARNING = "Warning"
    CRITICAL = "Critical"


@dataclass
class Alert:
    timestamp: str
    level: Level
    type: str
    message: str
    src_ip: str


class BaseDetector(ABC):
    def __init__(self):
        self.last_alert_times: Dict[Tuple[str, str], datetime] = {}
        self.cooldown = timedelta(minutes=5)

    def _should_alert(self, alert_type: str, src_ip: str) -> bool:
        """
        Check if an alert of the given type and source IP should be triggered,
        respecting the cooldown period.
        """
        key = (alert_type, src_ip)
        now = datetime.now()

        if key in self.last_alert_times:
            if now - self.last_alert_times[key] < self.cooldown:
                return False

        self.last_alert_times[key] = now
        return True

    @abstractmethod
    def analyze(self, dataframe: pd.DataFrame) -> List[Alert]:
        """
        Analyzes the given DataFrame and returns a list of Alert objects.
        """
        pass
