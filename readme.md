# 🗨️ Talk Let – Lightweight Self-Hosted Chat Platform

**Talk Let** is a lightweight, real-time, self-hosted chat platform built using **Golang**, **WebSockets**, **Alpine.js**, and **PostgreSQL**. Designed with simplicity, privacy, and performance in mind, Talk Let offers an open-source alternative to commercial tools like Microsoft Teams, Slack, and WhatsApp — without the cost, lock-in, or bloat.

## 🚀 Features

### 🔧 Core Functionality
- **Real-time messaging** with WebSockets
- **Session-based authentication** (mobile number + OTP flow)
- **User login/logout** with secure session management
- **One-on-one chat support**
- **Alpine.js-based conditional rendering** (no big frontend framework)

### 🗃️ Storage & Persistence
- **PostgreSQL** for chat and user data
- **MinIO** (S3-compatible) for file and media storage
- **Dockerized** setup for simple deployment

### 🌐 Tech Stack
| Layer      | Tool/Tech            |
|------------|----------------------|
| Backend    | Go (Golang)          |
| Frontend   | HTML, CSS, JS, Alpine.js |
| Realtime   | WebSockets           |
| Auth       | Session-based        |
| Database   | PostgreSQL           |
| Storage    | MinIO (S3)           |
| Container  | Docker               |

---

## 🌍 Why Talk Let?

> A self-hosted, privacy-first chat system — no third-party vendors, no monthly fees.

### ✅ Perfect For:
- 🏫 **Educational Institutions** — college chat systems, campus teams  
- 🧑‍💻 **Developer Teams** — who want self-hosted internal messaging  
- 🌐 **NGOs & FOSS Communities** — where privacy and transparency matter  
- 🏠 **Local Network Usage** — offices, labs, rural intranet, and more
