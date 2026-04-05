from typing import Dict, List, Set, Tuple

import pandas as pd

from internal.analyzer.base import Alert, BaseDetector, Level


class SuspiciousProtocolDetector(BaseDetector):
    def __init__(self):
        super().__init__()
        self.common_ports = {
            80: ["TCP", "HTTP"],
            443: ["TCP", "HTTPS"],
            53: ["UDP", "TCP", "DNS"],
            22: ["TCP", "SSH"],
            25: ["TCP", "SMTP"],
            5353: ["UDP", "mDNS"],
        }
        self.suspicious_mappings = {
            80: ["SSH"],
            22: ["UDP"],
        }

    def analyze(self, df: pd.DataFrame) -> List[Alert]:
        alerts = []
        if (
            df.empty
            or "src_ip" not in df.columns
            or "dst_port" not in df.columns
            or "protocol" not in df.columns
        ):
            return alerts

        clean_df = df.dropna(subset=["src_ip", "dst_port", "protocol"])
        clean_df = clean_df[clean_df["src_ip"].str.strip() != ""]

        if clean_df.empty:
            return alerts

        for _, row in clean_df.iterrows():
            ip = row["src_ip"]
            port = int(row["dst_port"])
            proto = str(row["protocol"]).upper()

            is_suspicious = False
            reason = ""

            if (
                port in self.suspicious_mappings
                and proto in self.suspicious_mappings[port]
            ):
                is_suspicious = True
                reason = f"Unusual protocol {proto} on port {port}"

            elif port in self.common_ports and proto not in self.common_ports[port]:
                pass

            if is_suspicious:
                if self._should_alert("SUSPICIOUS_PROTOCOL", ip):
                    alerts.append(
                        Alert(
                            timestamp=str(row["timestamp"]),
                            level=Level.WARNING,
                            type="SUSPICIOUS_PROTOCOL",
                            message=f"[{ip}] {reason}",
                            src_ip=ip,
                        )
                    )

        return alerts
