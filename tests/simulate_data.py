import sqlite3
import random
import time
from datetime import datetime, timedelta

DB_PATH = "data/traffic.db"

def simulate():
    print(f"Connecting to {DB_PATH}...")
    conn = sqlite3.connect(DB_PATH, isolation_level=None)
    cursor = conn.cursor()

    # Ensure tables exist (normally Go engine does this, but we want standalone test capability)
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS `traffic_logs` (
        `id` integer PRIMARY KEY AUTOINCREMENT,
        `timestamp` datetime,
        `src_ip` text,
        `dst_ip` text,
        `src_mac` text,
        `dst_mac` text,
        `src_port` integer,
        `dst_port` integer,
        `protocol` text,
        `ttl` integer,
        `payload_size` integer,
        `packet_count` integer
    )""")
    
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS `alerts` (
        `id` integer PRIMARY KEY AUTOINCREMENT,
        `timestamp` datetime,
        `level` text,
        `type` text,
        `message` text,
        `src_ip` text
    )""")

    print("Inserting 1000 traffic logs...")
    protocols = ["TCP", "UDP", "mDNS", "ARP", "HTTP", "HTTPS", "DNS"]
    ips = [f"192.168.0.{i}" for i in range(2, 255)] + ["45.33.22.11", "99.88.77.66"]
    
    now = datetime.now()
    
    traffic_data = []
    for i in range(1000):
        ts = (now - timedelta(seconds=random.randint(0, 3600))).strftime("%Y-%m-%d %H:%M:%S")
        src = random.choice(ips)
        dst = "192.168.0.1"
        proto = random.choice(protocols)
        traffic_data.append((
            ts, src, dst, "AA:BB:CC:DD:EE:FF", "00:11:22:33:44:55",
            random.randint(1024, 65535), 80, proto, 64, random.randint(64, 1500), random.randint(1, 100)
        ))
    
    cursor.executemany("""
        INSERT INTO traffic_logs (timestamp, src_ip, dst_ip, src_mac, dst_mac, src_port, dst_port, protocol, ttl, payload_size, packet_count)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    """, traffic_data)

    print("Inserting 50 alerts (clogging test)...")
    alert_types = [
        ("CRITICAL", "ARP_SPOOF", "Address conflict: IP seen on different devices"),
        ("CRITICAL", "TRAFFIC_FLOOD", "Potential DoS attack detected!"),
        ("WARNING", "PORT_SCAN", "Host is scanning ports"),
        ("WARNING", "SUSPICIOUS_PROTOCOL", "Unusual protocol on port"),
        ("INFO", "NEW_DEVICE_MDNS", "New device discovered via mDNS"),
        ("INFO", "UNUSUAL_TTL", "External IP with unusual TTL")
    ]

    alert_data = []
    for i in range(50):
        ts = (now - timedelta(seconds=random.randint(0, 3600))).strftime("%Y-%m-%d %H:%M:%S")
        level, atype, msg = random.choice(alert_types)
        src = random.choice(ips)
        alert_data.append((ts, level.capitalize(), atype, f"{msg} ({src})", src))

    cursor.executemany("""
        INSERT INTO alerts (timestamp, level, type, message, src_ip)
        VALUES (?, ?, ?, ?, ?)
    """, alert_data)

    conn.close()
    print("Simulation complete. Database 'clogged' successfully.")

if __name__ == "__main__":
    simulate()
