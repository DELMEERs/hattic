import os
import sqlite3
import time
from datetime import datetime

import pandas as pd
import plotly.express as px
import plotly.graph_objects as go
import streamlit as st

st.set_page_config(
    page_title="hattic.",
    page_icon="assets/icons/icon.png",
    layout="wide",
    initial_sidebar_state="expanded",
)


DB_PATH = "data/traffic.db"


def get_connection():
    """Returns a sqlite3 connection compatible with WAL mode."""
    return sqlite3.connect(DB_PATH, isolation_level=None)


@st.cache_data(ttl=2)
def load_alerts():
    if not os.path.exists(DB_PATH):
        return pd.DataFrame()
    try:
        conn = get_connection()
        df = pd.read_sql_query("SELECT * FROM alerts ORDER BY timestamp DESC", conn)
        conn.close()
        return df
    except Exception as e:
        st.error(f"Error loading alerts: {e}")
        return pd.DataFrame()


@st.cache_data(ttl=2)
def load_traffic_stats():
    if not os.path.exists(DB_PATH):
        return pd.DataFrame()
    try:
        conn = get_connection()
        query = "SELECT protocol, SUM(packet_count) as total_packets FROM traffic_logs GROUP BY protocol"
        df = pd.read_sql_query(query, conn)
        conn.close()
        return df
    except Exception as e:
        st.error(f"Error loading traffic stats: {e}")
        return pd.DataFrame()


def clear_database():
    try:
        conn = get_connection()
        cursor = conn.cursor()
        cursor.execute("DELETE FROM alerts")
        cursor.execute("DELETE FROM traffic_logs")
        conn.close()
        st.success("Database cleared successfully!")
        time.sleep(1)
        st.rerun()
    except Exception as e:
        st.error(f"Error clearing database: {e}")


def sidebar():
    with st.sidebar:
        logo_path = "assets/icons/icon.png"
        if os.path.exists(logo_path):
            st.image(logo_path, width=80)

        st.title("hattic.")
        st.markdown("---")

        refresh_rate = st.slider("Refresh Rate (seconds)", 2, 30, 5)

        st.markdown("---")
        if st.button("🗑️ Clear Database", width="stretch"):
            clear_database()

        st.markdown("---")
        st.info(
            "Hattic Network IDS Dashboard. Monitoring active traffic and security threats."
        )

        return refresh_rate


def main_dashboard():
    col_title, col_status = st.columns([4, 1])
    with col_title:
        st.title("Network Security Operations Center")
    with col_status:
        st.success("System Live")

    df_alerts = load_alerts()
    df_traffic = load_traffic_stats()

    m1, m2, m3, m4 = st.columns(4)

    total_alerts = len(df_alerts)
    critical_alerts = (
        len(df_alerts[df_alerts["level"] == "Critical"]) if not df_alerts.empty else 0
    )
    warning_alerts = (
        len(df_alerts[df_alerts["level"] == "Warning"]) if not df_alerts.empty else 0
    )
    unique_ips = df_alerts["src_ip"].nunique() if not df_alerts.empty else 0

    m1.metric("Total Alerts", total_alerts)
    m2.metric("Critical Threats", critical_alerts, delta_color="inverse")
    m3.metric("Warnings", warning_alerts)
    m4.metric("Suspect IPs", unique_ips)

    st.markdown("---")

    c1, c2 = st.columns(2)

    with c1:
        st.subheader("Alert Distribution by Type")
        if not df_alerts.empty:
            type_counts = df_alerts["type"].value_counts().reset_index()
            type_counts.columns = ["Alert Type", "Count"]
            fig = px.pie(
                type_counts,
                values="Count",
                names="Alert Type",
                hole=0.4,
                color_discrete_sequence=px.colors.qualitative.Pastel,
            )
            fig.update_layout(margin=dict(t=0, b=0, l=0, r=0), height=300)
            st.plotly_chart(fig, width="stretch")
        else:
            st.info("No alert data available yet.")

    with c2:
        st.subheader("Traffic Volume by Protocol")
        if not df_traffic.empty:
            fig = px.bar(
                df_traffic,
                x="protocol",
                y="total_packets",
                labels={"protocol": "Protocol", "total_packets": "Total Packets"},
                color="protocol",
                color_discrete_sequence=px.colors.qualitative.Safe,
            )
            fig.update_layout(
                margin=dict(t=20, b=20, l=20, r=20), height=300, showlegend=False
            )
            st.plotly_chart(fig, width="stretch")
        else:
            st.info("No traffic data available yet.")

    st.markdown("---")

    c3, c4 = st.columns(2)

    with c3:
        st.subheader("Alert Activity over Time")
        if not df_alerts.empty:
            df_alerts["timestamp"] = pd.to_datetime(df_alerts["timestamp"])
            time_series = (
                df_alerts.set_index("timestamp").resample("1min").size().reset_index()
            )
            time_series.columns = ["Time", "Alert Count"]
            fig = px.line(time_series, x="Time", y="Alert Count", markers=True)
            fig.update_layout(height=300, margin=dict(t=20, b=20, l=20, r=20))
            st.plotly_chart(fig, width="stretch")
        else:
            st.info("Waiting for time-series data...")

    with c4:
        st.subheader("Top Suspect IPs")
        if not df_alerts.empty:
            top_ips = df_alerts["src_ip"].value_counts().head(5).reset_index()
            top_ips.columns = ["Source IP", "Alert Count"]
            fig = px.bar(
                top_ips,
                x="Alert Count",
                y="Source IP",
                orientation="h",
                color="Alert Count",
                color_continuous_scale="Reds",
            )
            fig.update_layout(height=300, margin=dict(t=20, b=20, l=20, r=20))
            st.plotly_chart(fig, width="stretch")
        else:
            st.info("No suspect IPs identified.")

    st.markdown("---")

    st.subheader("Real-time Alert Log")
    if not df_alerts.empty:

        def color_level(val):
            if val == "Critical":
                return "background-color: #ff4b4b; color: white"
            if val == "Warning":
                return "background-color: #ffa500; color: black"
            if val == "Info":
                return "background-color: #1c83e1; color: white"
            return ""

        display_df = df_alerts[["timestamp", "level", "type", "src_ip", "message"]]
        st.dataframe(
            display_df.style.map(color_level, subset=["level"]),
            width="stretch",
            height=400,
        )
    else:
        st.info("Listening for security events...")


if __name__ == "__main__":
    refresh = sidebar()
    main_dashboard()
    time.sleep(refresh)
    st.rerun()
