from typing import List, Set

import pandas as pd

from internal.analyzer.base import Alert, BaseDetector, Level


class MDNSScanner(BaseDetector):
    def __init__(self):
        super().__init__()
        self.known_ips: Set[str] = set()

    def analyze(self, df: pd.DataFrame) -> List[Alert]:
        alerts = []
        if df.empty or "protocol" not in df.columns or "src_ip" not in df.columns:
            return alerts

        clean_df = df.dropna(subset=["src_ip", "protocol"])
        clean_df = clean_df[clean_df["src_ip"].str.strip() != ""]

        if clean_df.empty:
            return alerts

        mdns_traffic = clean_df[clean_df["protocol"].str.upper() == "MDNS"]

        for _, row in mdns_traffic.iterrows():
            src_ip = row["src_ip"]
            ip_str = str(src_ip)

            if ip_str not in self.known_ips:
                self.known_ips.add(ip_str)
                hostname = row.get("hostname", "Unknown Hostname")
                if pd.isna(hostname) or not str(hostname).strip():
                    hostname = "Unknown Hostname"

                alerts.append(
                    Alert(
                        timestamp=str(row["timestamp"]),
                        level=Level.INFO,
                        type="NEW_DEVICE_MDNS",
                        message=f"New device discovered via mDNS: {ip_str} ({hostname})",
                        src_ip=ip_str,
                    )
                )

        return alerts
