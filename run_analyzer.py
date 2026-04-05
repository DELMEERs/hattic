import time

from internal.analyzer.detectors.arp_spoof import ARPSpoofDetector
from internal.analyzer.detectors.mdns_scanner import MDNSScanner
from internal.analyzer.detectors.port_scanner import PortScannerDetector
from internal.analyzer.detectors.suspicious_protocol import SuspiciousProtocolDetector
from internal.analyzer.detectors.traffic_flood import TrafficFloodDetector
from internal.analyzer.detectors.unusual_ttl import UnusualTTLDetector
from internal.analyzer.manager import DetectionManager


def main():
    manager = DetectionManager(db_path="data/traffic.db")


    manager.register_detector(ARPSpoofDetector())
    manager.register_detector(
        TrafficFloodDetector(
            info_threshold=100000, warning_threshold=500000, critical_threshold=1000000
        )
    )
    manager.register_detector(MDNSScanner())
    manager.register_detector(PortScannerDetector(port_threshold=25))
    manager.register_detector(SuspiciousProtocolDetector())
    manager.register_detector(UnusualTTLDetector())

    print("--- Traffic Analysis Engine Started ---")
    print("Monitoring database for security alerts...")

    try:
        while True:
            alerts = manager.run_analysis()

            if alerts:
                print(
                    f"\n[{time.strftime('%Y-%m-%d %H:%M:%S')}] Found {len(alerts)} alerts:"
                )
                for alert in alerts:
                    print(f"  [{alert.level.value}] {alert.type}: {alert.message}")
            else:
                # print(".", end="", flush=True)
                pass

            time.sleep(10)
    except KeyboardInterrupt:
        print("\nShutting down analysis engine...")


if __name__ == "__main__":
    main()
