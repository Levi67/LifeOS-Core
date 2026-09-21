# 🚀 CoreHub (Working Title)

> **The Privacy-First, Local-AI Powered Personal Life Operating System.**
> A modular, self-hosted platform connecting your tasks, knowledge, habits, financials, and metrics with zero cloud lock-in and 100% offline intelligence.

---

## 📌 Vision & Philosophy

* **100% Data Ownership:** Your personal life data stays on your server in transparent, human-readable formats (JSON/SQLite/Markdown).
* **Local-First AI Integration:** Harness LLMs and local embeddings via local runtimes (e.g., Ollama) without sending sensitive personal data to cloud APIs. Option to add cloud API keys later for low-spec servers.
* **Extreme Modular Architecture:** No monolithic bloat. The core engine strictly handles auth, storage, event routing, and the plugin runtime. Everything else is an addon.
* **Self-Hosted Core & Multi-Device Access:** Host the main server yourself on your home hardware, then access everything seamlessly via a web interface or mobile app.
* **Multi-User System:** Native user management allowing different people to host isolated profiles on a single server instance.
* **Custom Layout System:** Define custom box-based dashboards—choose what widgets display, where they get data, and what runs in the background.

---

## 🏗️ System Architecture

The project follows an event-driven, micro-kernel architecture to keep the core binary light and fast.


```

┌─────────────────────────────────────────────────────────────┐
│                       CLIENT LAYER                          │
│            Web App (Svelte/React) / Mobile App              │
└──────────────────────────────┬──────────────────────────────┘
                   │ WebSocket / REST API |
┌──────────────────────────────▼──────────────────────────────┐
│                        CORE SERVER                          │
│  ┌──────────────┐   ┌──────────────┐   ┌─────────────────┐  │
│  │ User / Auth  │   │ Layout Engine│   │ Event Bus       │  │
│  │ (Multi-user) │   │ (Dashboard)  │   │ (Background jobs)│  │
│  └──────────────┘   └──────────────┘   └─────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                  PLUGIN RUNTIME / SDK                 │  │
│  └───────────────────────────────────────────────────────┘  │
└──────────────┬──────────────────────┬───────────────────────┘
               │                      │
┌──────────────▼────────────┐ ┌───────▼───────────────────────┐
│     STORAGE / DATA        │ │        LOCAL AI ENGINE        │
│ (SQLite / Local Files)    │ │   (Ollama / Local Embeddings) │
└───────────────────────────┘ └───────────────────────────────┘

```

---

## 🧩 Plugin & Addon Framework

Everything in CoreHub is a module. Plugins provide two main components:
1. **Background Drivers / Services:** Fetch external data (stock APIs, Google Fit / Apple Health sync, open banking), run scheduled tasks, or process AI jobs.
2. **UI Widgets (Boxes):** Dynamic frontend components registered into the global grid layout engine.

### Manifest Spec Example (`plugin.json`)

```json
{
  "id": "com.corehub.calorie-tracker",
  "name": "Calorie & Health Sync",
  "version": "1.0.0",
  "description": "Tracks daily nutrition and syncs steps via Google Fit / Apple Health",
  "permissions": [
    "network:external",
    "storage:local",
    "ai:llm_query"
  ],
  "widgets": [
    {
      "id": "macro-summary-card",
      "name": "Daily Macros",
      "defaultGridSize": { "w": 2, "h": 2 },
      "entry": "web/dist/macro-card.js"
    }
  ],
  "backgroundTasks": [
    {
      "name": "sync_fitness_data",
      "schedule": "0 */1 * * *"
    }
  ]
}

```

---

## 🎯 Target Core Modules (MVP Roadmap)

### 1. 🧠 Knowledge & Tasks (Phase 1)

* **Markdown Wiki:** Fast, local note-taking with bi-directional linking.
* **Task Engine:** Fast keyboard-first todo list and Kanban board.
* **Local AI Integration:** Semantic search over local notes using vector embeddings. Auto-summarizes daily completed tasks.

### 2. 📊 Health & Fitness (Phase 2)

* **Calorie & Fitness Tracker:** Integrations for Google Fit / Apple Health APIs or manual quick-log entries.
* **AI Macro Assistant:** Input plain text (*"2 eggs and an avocado"*), and the local LLM automatically parses and logs the macros.

### 3. 📈 Finance & Assets (Phase 3)

* **Stock & Portfolio Tracker:** Real-time stock ticker, ETF updates, and asset overview widgets.
* **Banking Integrations:** Custom modular adapter for importing CSV transactions or hooking into open-banking APIs safely.

---

## 🤖 AI Layer Configuration

* **Primary Provider:** Local Ollama runtime / OpenAI-compatible local endpoint (100% offline).
* **Fallback Provider (Optional):** User-configured cloud API keys (e.g., Anthropic Claude, OpenAI) for setups where the server hardware lacks GPU acceleration.
* **Privacy Controls:** Strict per-plugin AI permissions. Addons must request explicit user authorization before accessing local context or external LLMs.

---

## 🔒 Security & Multi-User Architecture

* **Multi-Tenancy:** Isolated storage schemas and dashboard configs per user account.
* **Role-Based Access:** Admin dashboard for system resource management, docker status, and plugin installation permissions.
* **API Key Management:** Encrypted secret storage for external API credentials (OAuth tokens, finance keys).

---

## 🛠️ Recommended Tech Stack

| Layer | Recommended Choice | Rationale |
| --- | --- | --- |
| **Backend Core** | Go or Rust | Fast, compiles to a single binary, minimal RAM footprint on home servers. |
| **Database** | SQLite + Vector Extension | Zero config, self-contained in local files. |
| **Frontend Web** | SvelteKit or React + Tailwind | Fast layout rendering with smooth drag-and-drop box placement. |
| **Mobile App** | React Native / Expo | Cross-platform mobile client sharing UI components with web. |
| **Local AI Engine** | Ollama API / local embeddings | Easy zero-config local AI deployment via Docker or native binary. |

---

## 🚦 Getting Started for Contributors

1. Clone repository: `git clone https://github.com/your-username/corehub.git`
2. Run core engine in dev mode: `docker-compose up -d`
3. Load sample plugin: `./corehub plugin install ./examples/sample-todo-plugin`

```

```
