import sqlite3
from datetime import datetime, timedelta
from typing import List

import pandas as pd

from internal.analyzer.base import Alert, BaseDetector


class DetectionManager:
    def __init__(self, db_path: str = "data/traffic.db"):
        self.db_path = db_path
        self.detectors: List[BaseDetector] = []

    def register_detector(self, detector: BaseDetector):
        """Adds a new detector to the registry."""
        self.detectors.append(detector)

    def _load_data(self) -> pd.DataFrame:
        """Loads latest traffic logs from the SQLite database."""
        try:
            conn = sqlite3.connect(self.db_path, isolation_level=None)

            time_threshold = (datetime.now() - timedelta(minutes=1)).strftime(
                "%Y-%m-%d %H:%M:%S"
            )

            query = f"SELECT * FROM traffic_logs WHERE timestamp >= '{time_threshold}'"

            df = pd.read_sql_query(query, conn)
            conn.close()
            return df
        except Exception as e:
            print(f"Error loading data from {self.db_path}: {e}")
            return pd.DataFrame()

    def save_alerts(self, alerts: List[Alert]):
        """Persists alerts to the SQLite database."""
        if not alerts:
            return

        try:
            conn = sqlite3.connect(self.db_path, isolation_level=None)
            cursor = conn.cursor()

            for alert in alerts:
                cursor.execute(
                    "INSERT INTO alerts (timestamp, level, type, message, src_ip) VALUES (?, ?, ?, ?, ?)",
                    (
                        alert.timestamp,
                        alert.level.value,
                        alert.type,
                        alert.message,
                        alert.src_ip,
                    ),
                )

            conn.close()
        except Exception as e:
            print(f"Error saving alerts to {self.db_path}: {e}")

    def run_analysis(self) -> List[Alert]:
        """
        Loads data and runs all registered detectors.
        Returns a consolidated list of alerts.
        """
        df = self._load_data()
        all_alerts = []

        if df.empty:
            return all_alerts

        for detector in self.detectors:
            alerts = detector.analyze(df)
            all_alerts.extend(alerts)

        self.save_alerts(all_alerts)
        return all_alerts
