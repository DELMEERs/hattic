from typing import Dict, List, Set

import pandas as pd

from internal.analyzer.base import Alert, BaseDetector, Level


class UnusualTTLDetector(BaseDetector):
    def __init__(self):
        super().__init__()
        # Standard OS values: 64 (Linux/Android), 128 (Windows), 255 (Network equipment)
        self.standard_ttls = {64, 128, 255}

    def analyze(self, df: pd.DataFrame) -> List[Alert]:
        alerts = []
        if df.empty or "src_ip" not in df.columns or "ttl" not in df.columns:
            return alerts

        clean_df = df.dropna(subset=["src_ip", "ttl"])
        clean_df = clean_df[
            (clean_df["src_ip"].str.strip() != "") & (clean_df["ttl"] != 0)
        ]

        if clean_df.empty:
            return alerts

        for _, row in clean_df.iterrows():
            ip = row["src_ip"]
            ttl = int(row["ttl"])

            # Ignore standard OS values
            if ttl in self.standard_ttls:
                continue

            is_local = ip.startswith(("192.168.", "10.", "172.16."))

            if is_local:
                # Ignore TTL=1 and TTL=2 for local traffic (standard for discovery)
                if ttl in {1, 2}:
                    continue
                # For local IPs, we follow the instruction to "Only alert on 'unusual' values from external IPs"
                continue
            else:
                # External IP: Only alert on "unusual" values (e.g., between 10 and 30)
                if 10 <= ttl <= 30:
                    if self._should_alert("UNUSUAL_TTL", ip):
                        alerts.append(
                            Alert(
                                timestamp=str(row["timestamp"]),
                                level=Level.INFO,
                                type="UNUSUAL_TTL",
                                message=f"External IP {ip} with unusual TTL: {ttl}",
                                src_ip=ip,
                            )
                        )

        return alerts
