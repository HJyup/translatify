# Translation Service

> [!NOTE]  
> This service is under construction and is not yet available for public use.

## Overview
The **Translatify Translation Service** is a Go-based microservice that provides AI-powered translations using **OpenAI's GPT-4**. It integrates with **RabbitMQ** for asynchronous message handling, utilizes **gRPC** for inter-service communication, and supports **service discovery** via Consul.

## Features
- **AI-Powered Translations**: Utilizes OpenAI's GPT-4 to provide high-quality translations.
- **Message Queue Processing**: Listens to RabbitMQ for translation requests.
- **Efficient Caching**: Implements an in-memory cache to reduce redundant API calls.

## Installation & Setup
### Prerequisites
- Docker
- Go (v1.22+)
- RabbitMQ
- OpenAI API Key
- Consul

### Steps to Run the Service
1. Clone the repository:
   ```sh
   git clone https://github.com/HJyup/translatify
   cd translation
   ```

2. Build and run the service using Docker:
   ```sh
   docker build -t translatify-translation .
   docker run -p 8080:8080 --env-file .env translatify-translation
   ```

3. Alternatively, run locally with Go:
   ```sh
   go run cmd/main.go
   ```

## API Usage
### **RabbitMQ Message Handling**
The service listens for messages in RabbitMQ on the `MessageSentEvent` queue. Messages should be in JSON format:
```json
{
  "content": "Hello, how are you?",
  "sourceLang": "en",
  "targetLang": "fr"
}
```
The response will contain the translated content.

## Architecture
1. A message arrives in **RabbitMQ**.
2. The **consumer** extracts text and calls the **translator**.
3. The **translator** checks cache; if found, it returns instantly.
4. If not cached, it queries **OpenAI** for translation.
5. The result is stored in cache and returned to the requester.

## Contributing
Feel free to submit issues or pull requests. Make sure to follow best practices and test your changes before submitting.

