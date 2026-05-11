<div align="center">
  <img src="assets/logo.png" width="200" alt="Hattic Logo">
  <h1>Hattic</h1>
  <p><strong>Defensive Intelligence for the Modern Network.</strong></p>

  <p>
    <a href="https://github.com/DELMEERs/hattic/actions/workflows/release.yml">
      <img src="https://img.shields.io/github/actions/workflow/status/DELMEERs/hattic/release.yml?style=for-the-badge&logo=github&label=Build" alt="GitHub Build Status">
    </a>
    <a href="https://github.com/DELMEERs/hattic/releases">
      <img src="https://img.shields.io/github/v/release/DELMEERs/hattic?style=for-the-badge&logo=github&color=blue" alt="Release Version">
    </a>
    <a href="LICENSE">
      <img src="https://img.shields.io/github/license/DELMEERs/hattic?style=for-the-badge&color=success" alt="MIT License">
    </a>
    <br>
    <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
    <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-lightgrey?style=for-the-badge" alt="Platforms">
  </p>
</div>

---

## 🛡️ The Essence of Hattic

**Hattic** is a high-performance, cross-platform Network Sniffer and Intrusion Detection System (IDS) engineered for clarity and speed. By marrying the low-level precision of **Go** with the fluid interactivity of **Vue.js**, Hattic provides a glass-pane view into your network's soul.

Whether you are a security researcher, a network administrator, or a developer debugging microservices, Hattic offers a sophisticated environment to capture, analyze, and defend your traffic without the overhead of traditional enterprise tools.

## ✨ Premium Features

-   🔍 **Real-Time Packet Analysis** – Deep inspection of network traffic with sub-millisecond latency.
-   🧠 **Smart Guard Dependency Check** – Intelligent auto-detection of `Npcap` (Windows) or `libpcap` (Linux), ensuring the engine starts only when the environment is ready.
-   📊 **Unified Dashboard** – A sleek, high-end interface built with Vite and Tailwind CSS for visualizing traffic spikes and protocol distribution.
-   🪶 **Lightweight Footprint** – Minimal CPU and memory usage, optimized for background monitoring.
-   🚨 **Intrusion Detection** – Built-in analyzers for detecting ARP spoofing, port scanning, and traffic floods.

---

## 🚀 Installation & Setup

Hattic is distributed as a portable binary. Select your platform below for a foolproof setup.

### 🪟 Windows
1.  **Download:** Grab the latest `.exe` from the [GitHub Releases](https://github.com/DELMEERs/hattic/releases) page.
2.  **Driver Requirement:** Hattic requires **Npcap**.
    -   [Download Npcap here](https://npcap.com/#download).
    -   **Important:** During installation, you **must** check the box for **"Install Npcap with WinPcap compatibility mode"** for the sniffer engine to function correctly.
3.  **Run:** Launch `hattic.exe` and grant Administrative privileges when prompted.

### 🐧 Linux
1.  **Download:** Grab the binary from the [GitHub Releases](https://github.com/DELMEERs/hattic/releases) page.
2.  **Permissions:** To capture raw packets without running the entire UI as root (highly recommended), execute the following:
    ```bash
    sudo setcap cap_net_raw,cap_net_admin=eip ./hattic
    ```
    *This grants the binary specific network capabilities while maintaining a secure, non-root execution environment.*
3.  **Run:** `./hattic`

---

## 🛠️ Tech Stack

Hattic is built using a modern, type-safe architecture.

| Layer | Technologies |
| :--- | :--- |
| **Frontend** | [Vue 3](https://vuejs.org/), [Vite](https://vitejs.dev/), [Tailwind CSS](https://tailwindcss.com/), [Lucide Icons](https://lucide.dev/) |
| **Backend** | [Go 1.22+](https://go.dev/), [gopacket](https://github.com/google/gopacket), [Wails v2](https://wails.io/) |
| **Infrastructure** | GitHub Actions (CI/CD), Go-Releaser |

---

## 🔮 Future Vision

We are constantly evolving. Here is our strategic roadmap:

-   [x] **Initial Core development in Go** – High-concurrency packet processing engine.
-   [x] **Cross-platform CI/CD pipeline** – Automated builds for Windows & Linux.
-   [x] **Intelligent dependency guard** – Robust handling of Npcap/libpcap environments.
-   [ ] **CLI/TUI Version** – A dedicated terminal interface for power-users and remote SSH sessions.
-   [ ] **REST/gRPC API** – Headless mode for integration with external monitoring stacks (ELK, Grafana).

---

## 🤝 Contributing & License

Hattic is an Open Source project. We welcome contributions that align with our vision of high-performance network security.

-   **License:** Distributed under the [MIT License](LICENSE).
-   **Security:** To report a vulnerability, please open a private GitHub Advisory.

<div align="center">
  <p>Built with ❤️ by dylan extreme</p>
  <a href="https://github.com/DELMEERs/hattic/stargazers">⭐ Star us on GitHub</a>
</div>

---

<div align="center">
  <p><b>Language / Язык</b></p>
  <a href="#hattic">English Version</a> •
  <a href="#hattic-1">Русская версия</a>
</div>