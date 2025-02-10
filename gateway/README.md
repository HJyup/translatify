# Gateway Service

> [!NOTE]  
> This service is under construction and is not yet available for public use.

## Overview
The **Gateway Service** is a central API gateway that routes requests to various microservices in the **Translatify** ecosystem. It provides authentication, request forwarding, WebSocket support, and API gateway functionalities while integrating with **Consul** for service discovery and **Clerk** for authentication.

## Features
- **API Gateway:** Routes requests to backend microservices.
- **Authentication & Security:** Uses **Clerk** for user authentication.
- **Service Discovery:** Registers with **Consul** to discover microservices dynamically.
- **WebSocket Support:** Provides real-time communication capabilities.
- **Observability & Tracing:** Integrated with **OpenTelemetry & Jaeger** for monitoring.
- **gRPC Gateway:** Converts gRPC requests to REST endpoints.

## Installation & Setup
### Prerequisites
- Docker
- Go (v1.22+)
- Consul
- Clerk API Key

### Steps to Run the Service
1. Clone the repository:
   ```sh
   git clone https://github.com/HJyup/translatify
   cd gateway
   ```

2. Build and run the service using Docker:
   ```sh
   docker build -t gateway-service .
   docker run -p 8080:8080 --env-file .env gateway-service
   ```

3. Alternatively, run locally with Go:
   ```sh
   go run cmd/main.go
   ```

## API Usage
### **Authentication & Routing**
The gateway validates authentication using **Clerk** before forwarding API requests to appropriate services.

### **Example Routes**
- `GET /chat` → Forwards requests to the Chat Service.

### **WebSocket Support**
The service supports WebSockets for real-time updates. Clients can establish WebSocket connections for event-driven messaging.

## Architecture
1. Client makes a **REST API request** to the Gateway.
2. The **Gateway authenticates** the user via Clerk.
3. If valid, it **routes the request** to the appropriate backend microservice.
4. The response is returned via the Gateway API.
5. If using WebSockets, real-time updates are streamed to clients.

## Contributing
Contributions are welcome! Follow best practices and submit PRs with test cases.

