from typing import Dict, List, Set

import pandas as pd

from internal.analyzer.base import Alert, BaseDetector, Level


class PortScannerDetector(BaseDetector):
    WHITELIST_IPS = {"192.168.0.1", "127.0.0.1"}
    IGNORE_PORTS = {1900, 5353, 3702}

    def __init__(self, port_threshold: int = 20):
        super().__init__()
        self.port_threshold = port_threshold

    def analyze(self, df: pd.DataFrame) -> List[Alert]:
        alerts = []
        if df.empty or "src_ip" not in df.columns or "dst_port" not in df.columns:
            return alerts

        clean_df = df.dropna(subset=["src_ip", "dst_port"])
        clean_df = clean_df[
            (clean_df["src_ip"].str.strip() != "") & (clean_df["dst_port"] != 0)
        ]

        if clean_df.empty:
            return alerts

        clean_df = clean_df[~clean_df["src_ip"].isin(self.WHITELIST_IPS)]
        clean_df = clean_df[~clean_df["dst_port"].isin(self.IGNORE_PORTS)]

        if clean_df.empty:
            return alerts

        stats = clean_df.groupby("src_ip")["dst_port"].nunique()
        scanners = stats[stats >= self.port_threshold].index.tolist()

        for ip in scanners:
            port_count = int(stats[ip])
            latest_timestamp = clean_df[clean_df["src_ip"] == ip]["timestamp"].max()

            if self._should_alert("PORT_SCAN", ip):
                alerts.append(
                    Alert(
                        timestamp=str(latest_timestamp),
                        level=Level.WARNING,
                        type="PORT_SCAN",
                        message=f"IP {ip} is scanning ports: {port_count} unique ports hit.",
                        src_ip=ip,
                    )
                )

        return alerts
