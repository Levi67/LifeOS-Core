# 🚀 CoreHub

A privacy-first, local-AI-powered personal dashboard and life operating system. Self-host your tasks, notes, habits, health metrics, and finances with zero cloud lock-in.

---

## ✨ Features

- **100% Offline AI:** Powered by local runtimes (like Ollama) for summaries, semantic search, and macro tracking without sending data to the cloud.
- **Modular Addon System:** Extend your dashboard with custom plugins for stocks, Google Fit, open banking, and more.
- **Custom Box Dashboard:** Drag-and-drop dynamic widgets to customize your layout.
- **Multi-User & Self-Hosted:** Host one server instance for isolated user accounts.
- **Cross-Platform:** Single binary deployment with a responsive web interface and mobile app access.

---

## 🏗️ Tech Stack

- **Backend:** Go / Rust (Fast, single-binary API server)
- **Database:** SQLite + Vector Extension
- **Frontend:** SvelteKit / React + Tailwind CSS
- **AI Engine:** Ollama / Local Embeddings API

---

## ⚡ Quick Start

```bash
git clone https://github.com/Levi67/LifeOS-Core/

cd LifeOS-Core

air

cd web

npm run dev
