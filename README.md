# Translatify

--- 

## Overview
**Translatify** is a **scalable microservices-based translation and messaging system** designed to provide real-time, AI-powered translations, seamless messaging, and an efficient API gateway. It leverages **gRPC, RabbitMQ, OpenTelemetry, and Consul** to enable distributed and high-performance communication between services.

## Features
- **AI-Powered Translations** → Uses **OpenAI GPT-4** for context-aware translations.
- **Real-Time Messaging** → Supports **one-to-one chat** with WebSocket support.
- **API Gateway** → Centralized API management with **authentication and routing**.
- **Service Discovery** → Uses **Consul** for dynamic service registration and discovery.
- **Distributed Tracing** → Implements **Jaeger & OpenTelemetry** for observability.
- **Message Queueing** → Uses **RabbitMQ** for asynchronous processing.
- **gRPC & REST Support** → Exposes services via **gRPC & gRPC-Gateway (REST APIs)**.

## Microservices Architecture
### **1. Translation Service**
- Provides AI-based translations using **OpenAI API**.
- Implements **caching** to optimize API usage.
- Listens to **RabbitMQ** for batch translation requests.

### **2. Chat Service**
- Manages **one-to-one messaging** between users.
- Stores chat history in **PostgreSQL**.
- Supports **WebSockets** for real-time updates.

### **3. Gateway Service**
- Acts as the **entry point** for all API requests.
- Uses **Clerk** for authentication and token validation.
- Routes requests to respective microservices.

### **4. Common Module**
- Provides shared **utilities** for logging, environment handling, JSON processing, and database scanning.
- Implements **service discovery, message queueing, and tracing**.

## Tech Stack
- **Programming Language:** Go
- **API Communication:** gRPC, gRPC-Gateway (REST)
- **Service Discovery:** Consul
- **Messaging Queue:** RabbitMQ
- **Tracing & Monitoring:** OpenTelemetry, Jaeger
- **Database:** PostgreSQL
- **Authentication:** Clerk
- **Containerization:** Docker

## Installation & Setup
### Prerequisites
- Docker
- Go (v1.22+)
- PostgreSQL
- RabbitMQ
- Consul
- Clerk API Key
- OpenAI API Key (for Translation Service)

### Steps to Run the System
1. Clone the repository:
   ```sh
   git clone https://github.com/HJyup/translatify.git
   cd translatify
   ```

2. Set up environment variables:
   ```sh
   cp .env.example .env
   ```
   Edit `.env` with necessary credentials (RabbitMQ, PostgreSQL, Consul, OpenAI API Key, Clerk API Key).

3. Build and run services using Docker:
   ```sh
   docker-compose up --build
   ```

4. Alternatively, run services manually:
   ```sh
   go run services/translation/cmd/main.go
   go run services/chat/cmd/main.go
   go run services/gateway/cmd/main.go
   ```

## Contributing
We welcome contributions! Please submit a PR following our guidelines.

## License
MIT License. See `LICENSE` for details.