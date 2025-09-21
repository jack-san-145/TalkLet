# Talk Let

**Talk Let** is a self-hosted, web-based communication platform designed specifically for academic institutions.  
It enables secure, structured, and real-time communication between students and professors, overcoming the limitations of generic messaging apps like WhatsApp or Telegram.  

The system is built with **Golang** on the backend, **HTMX + Alpine.js** on the frontend, and uses **PostgreSQL** for database management.  
Media files are stored securely in **MinIO object storage**, ensuring full data privacy within institutional infrastructure.  

---

## 🚀 Features

- **Authentication**
  - Student login using **Roll Number and Password**
  - Session-based authentication (no third-party OAuth)

- **Messaging**
  - Real-time **1-to-1 chat**
  - Department/class-based **group chats**
  - Powered by **WebSockets** for instant delivery

- **Group Management**
  - Professors can create groups by uploading **Excel sheets**
  - Automatic student mapping into their respective groups

- **Isolated Department Communication**
  - Each department has its own private chat environment
  - Prevents cross-department message leakage
  - Ensures academic boundaries and structured interaction

- **Media Sharing**
  - Share images, videos, and documents
  - Files stored securely using **MinIO**

- **Frontend**
  - Built with **HTMX** and **Alpine.js** for a lightweight, fast UI
  - Fully responsive (works on desktop and mobile browsers)
  - Dark mode support and emoji picker integration

- **Backend**
  - Developed in **Go**, optimized for concurrency handling
  - **PostgreSQL** as the relational database
  - **Redis** for caching and session management

- **Hosting**
  - Fully **self-hosted** on institutional servers
  - Provides complete control over data and privacy

- **Future Scope**
  - End-to-End Encryption (AES + RSA hybrid model)
  - Advanced admin dashboards for professors
  - Analytics on group activity
  - Cloud backup and file versioning

---

## 📖 Summary

**Talk Let** empowers colleges with their own secure, scalable, and modern messaging platform.  
It combines **real-time chat**, **media sharing**, and **department-based group management** under one self-hosted system.  
By providing **isolated department communication** and **full data ownership**, Talk Let ensures privacy, reliability, and adaptability for academic institutions.

