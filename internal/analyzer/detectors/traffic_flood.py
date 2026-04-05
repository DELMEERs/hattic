from typing import Dict, List, Optional

import pandas as pd

from internal.analyzer.base import Alert, BaseDetector, Level


class TrafficFloodDetector(BaseDetector):
    def __init__(
        self,
        info_threshold: int = 100000,
        warning_threshold: int = 500000,
        critical_threshold: int = 1000000,
    ):
        super().__init__()
        self.info_threshold = info_threshold
        self.warning_threshold = warning_threshold
        self.critical_threshold = critical_threshold

    def analyze(self, df: pd.DataFrame) -> List[Alert]:
        alerts = []
        if df.empty or "src_ip" not in df.columns:
            return alerts

        clean_df = df.dropna(subset=["src_ip"])
        clean_df = clean_df[clean_df["src_ip"].str.strip() != ""]

        if clean_df.empty:
            return alerts

        agg_dict = {"timestamp": "max"}
        has_packet_count = "packet_count" in clean_df.columns
        if has_packet_count:
            agg_dict["packet_count"] = "sum"

        stats = clean_df.groupby("src_ip").agg(agg_dict)

        for ip, row in stats.iterrows():
            ip_str = str(ip)
            p_count = int(row["packet_count"]) if has_packet_count else 0

            level: Optional[Level] = None
            desc = ""

            if p_count >= self.critical_threshold:
                level = Level.CRITICAL
                desc = "Potential DoS attack detected!"
            elif p_count >= self.warning_threshold:
                level = Level.WARNING
                desc = "Extreme network activity (possible flood)."
            elif p_count >= self.info_threshold:
                level = Level.INFO
                desc = "High network activity (streaming/downloading)."

            if level:
                if self._should_alert("TRAFFIC_FLOOD", ip_str):
                    alerts.append(
                        Alert(
                            timestamp=str(row["timestamp"]),
                            level=level,
                            type="TRAFFIC_FLOOD",
                            message=f"[{ip_str}] {desc} ({p_count} packets)",
                        )
                    )

        return alerts
