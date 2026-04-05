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
            conn = sqlite3.connect(self.db_path)

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

        return all_alerts
