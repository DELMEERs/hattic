from typing import Dict, List, Set

import pandas as pd

from internal.analyzer.base import Alert, BaseDetector, Level


class ARPSpoofDetector(BaseDetector):
    def __init__(self):
        super().__init__()

    def analyze(self, df: pd.DataFrame) -> List[Alert]:
        alerts = []
        if df.empty or "src_ip" not in df.columns or "src_mac" not in df.columns:
            return alerts

        ignored_macs = ["ff:ff:ff:ff:ff:ff", "00:00:00:00:00:00", ""]
        clean_df = df.dropna(subset=["src_ip", "src_mac"])
        clean_df = clean_df[
            (clean_df["src_ip"].str.strip() != "")
            & (~clean_df["src_mac"].str.lower().isin(ignored_macs))
        ]

        if clean_df.empty:
            return alerts

        grouped = clean_df.groupby("src_ip")["src_mac"].nunique()
        spoofed_ips = grouped[grouped > 1].index.tolist()

        for ip in spoofed_ips:
            macs = sorted(
                clean_df[clean_df["src_ip"] == ip]["src_mac"].unique().tolist()
            )
            latest_timestamp = clean_df[clean_df["src_ip"] == ip]["timestamp"].max()

            if self._should_alert("ARP_SPOOF", ip):
                alerts.append(
                    Alert(
                        timestamp=str(latest_timestamp),
                        level=Level.CRITICAL,
                        type="ARP_SPOOF",
                        message=f"Address conflict: IP {ip} seen on different devices: {', '.join(macs)}",
                        src_ip=ip,
                    )
                )

        return alerts
