# Chat Service

> [!NOTE]  
> This service is under construction and is not yet available for public use.

## Overview
The **Chat Service** is a microservice responsible for handling one-to-one text messaging between users. It supports real-time messaging, message retrieval, and integrates with RabbitMQ for asynchronous message handling. The service uses **gRPC** for communication, **PostgreSQL** for message storage, and **Consul** for service discovery.

## Features
- **One-to-One Messaging:** Users can send direct messages to each other.
- **Real-Time Updates:** Subscribe to a server-streaming endpoint for live message updates.
- **Message Retrieval:** Retrieve specific messages or list conversation history.
- **RabbitMQ Integration:** Supports asynchronous message processing.
- **Database Support:** Uses PostgreSQL for message storage.
- **Service Discovery:** Registers with Consul for service lookup.
- **Tracing & Observability:** Integrated with Jaeger for distributed tracing.

## Installation & Setup
### Prerequisites
- Docker
- Go (v1.22+)
- PostgreSQL
- RabbitMQ
- Consul

### Steps to Run the Service
1. Clone the repository:
   ```sh
   git clone https://github.com/HJyup/translatify
   cd chat
   ```

2. Build and run the service using Docker:
   ```sh
   docker build -t chat-service .
   docker run -p 8080:8080 --env-file .env chat-service
   ```

3. Alternatively, run locally with Go:
   ```sh
   go run cmd/main.go
   ```

## API Usage
### **gRPC API (Example)**
The gRPC server runs on `0.0.0.0:8080` and exposes endpoints for sending and retrieving messages:
```proto
rpc SendMessage (SendMessageRequest) returns (SendMessageResponse);
rpc GetMessage (GetMessageRequest) returns (GetMessageResponse);
rpc StreamMessages (StreamMessagesRequest) returns (stream Message);
```

## Architecture
1. A user sends a message via the **gRPC API**.
2. The message is **stored in PostgreSQL**.
3. **RabbitMQ** handles message processing and notifications.
4. The recipient can **stream messages in real-time**.
5. The service **registers with Consul** for discovery.

## Contributing
Contributions are welcome! Please follow best practices and ensure tests are included before submitting PRs.
