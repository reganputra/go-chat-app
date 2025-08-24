# Go Chat App

A real-time chat application built with Go, featuring WebSocket communication, JWT authentication, dual database support, and comprehensive monitoring with the ELK stack.

## Overview

Go Chat App is a modern, scalable real-time messaging application that combines the power of Go's concurrency with WebSocket technology to deliver instant messaging capabilities. The application features a robust architecture with user authentication, message persistence, and comprehensive monitoring and logging.

## Tech Stack

### Backend Framework
- **Go 1.24** - Core programming language
- **Fiber v2** - Fast HTTP web framework with Express-like features
- **WebSocket** - Real-time bidirectional communication

### Authentication & Security
- **JWT (JSON Web Tokens)** - Secure authentication with access and refresh tokens
- **bcrypt** - Password hashing
- **Rate Limiting** - API protection (50 requests per minute per IP)

### Databases
- **MySQL** - User data and session management (via GORM)
- **MongoDB** - Chat message history and persistence

### Monitoring & Observability
- **Elastic APM** - Application Performance Monitoring
- **ELK Stack** - Complete logging and monitoring solution
  - **Elasticsearch** - Search and analytics engine
  - **Logstash** - Log processing pipeline
  - **Kibana** - Data visualization and dashboards
  - **Filebeat** - Log shipping

### Additional Features
- **Docker** - Containerization support
- **HTML Templating** - Server-side rendering
- **Environment Configuration** - Flexible configuration management
- **Structured Logging** - Comprehensive application logging

## Project Structure

```
go-chat-app/
├── app/
│   ├── controllers/           # HTTP request handlers
│   ├── models/               # Data models and validation
│   ├── repositories/         # Data access layer
│   └── websocket/           # WebSocket implementation
├── bootstrap/               # Application initialization
├── elk_stack/              # ELK monitoring configuration
├── logs/                   # Application logs
├── pkg/
│   ├── database/          # Database setup and configuration
│   ├── jwt/               # JWT token management
│   ├── response/          # Standardized API responses
│   └── router/            # HTTP routing and middleware
├── views/                 # HTML templates
├── docker-compose.yaml    # Main application containers
└── main.go               # Application entry point
```

## API Endpoints

### Base URL
- **HTTP API**: `http://localhost:4000/api`
- **WebSocket**: `ws://localhost:8080/message/v1/send`
- **Monitoring Dashboard**: `http://localhost:4000/dashboard`

### Authentication Endpoints

#### Register User
```
POST /api/user/v1/register
Content-Type: application/json

{
    "username": "string (6-20 chars, unique)",
    "password": "string (min 6 chars)",
    "full_name": "string (min 6 chars)"
}
```

#### Login User
```
POST /api/user/v1/login
Content-Type: application/json

{
    "username": "string",
    "password": "string"
}

Response:
{
    "username": "string",
    "full_name": "string",
    "token": "jwt_access_token",
    "refresh_token": "jwt_refresh_token"
}
```

#### Logout User
```
DELETE /api/user/v1/logout
Authorization: Bearer {access_token}
```

#### Refresh Token
```
PUT /api/user/v1/refresh-token
Authorization: Bearer {refresh_token}

Response:
{
    "token": "new_access_token",
    "refresh_token": "new_refresh_token"
}
```

### Message Endpoints

#### Get Message History
```
GET /api/message/v1/history
Authorization: Bearer {access_token}
```

### WebSocket Endpoint

#### Send Real-time Messages
```
WebSocket: ws://localhost:8080/message/v1/send

Message Format:
{
    "from": "username",
    "message": "message content",
    "date": "2025-01-24T09:10:00Z"
}
```

## Database Schema

### MySQL Tables (User Management)

#### Users Table
```sql
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

#### User Sessions Table
```sql
CREATE TABLE user_sessions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    token_expired TIMESTAMP NOT NULL,
    refresh_token_expired TIMESTAMP NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### MongoDB Collections (Message Storage)

#### Chat History Collection
```json
{
    "_id": "ObjectId",
    "from": "username",
    "message": "message content",
    "date": "ISODate"
}
```

## Installation & Setup

### Prerequisites
- Go 1.24 or higher
- Docker and Docker Compose
- MySQL database
- MongoDB database

### Environment Variables
Create a `.env` file in the root directory:

```env
# Application Configuration
APP_HOST=localhost
APP_PORT=4000
APP_PORT_SOCKET=8080

# MySQL Database
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_NAME=go_chat_app

# MongoDB
MONGODB_URI=mongodb://localhost:27017

# JWT Configuration (add your JWT secrets)
JWT_SECRET=your_jwt_secret_key
```

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-chat-app
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up databases**
   - Start MySQL and create database `go_chat_app`
   - Start MongoDB service

4. **Start ELK Stack (Optional)**
   ```bash
   cd elk_stack
   docker-compose up -d
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

### Docker Setup

1. **Start all services**
   ```bash
   docker-compose up -d
   ```

2. **Start ELK monitoring**
   ```bash
   cd elk_stack
   docker-compose up -d
   ```

## WebSocket Real-time Messaging

The application uses WebSocket for real-time communication:

1. **Connect** to `ws://localhost:8080/message/v1/send`
2. **Send messages** in JSON format with `from`, `message`, and `date` fields
3. **Receive messages** broadcasted to all connected clients
4. Messages are automatically **persisted** to MongoDB
5. **APM tracing** is applied to all WebSocket operations

### Client-side WebSocket Example
```javascript
const socket = new WebSocket('ws://localhost:8080/message/v1/send');

// Send message
socket.send(JSON.stringify({
    from: 'username',
    message: 'Hello, World!',
    date: new Date().toISOString()
}));

// Receive messages
socket.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log(`${message.from}: ${message.message}`);
};
```

## Monitoring & Observability

### Application Performance Monitoring
- **Elastic APM** tracks all HTTP requests and WebSocket operations
- **Custom tracing** for database operations and authentication flows
- Performance metrics available at APM server (`http://localhost:8200`)

### ELK Stack Monitoring
- **Elasticsearch**: `http://localhost:9200`
- **Kibana Dashboard**: `http://localhost:5601`
- **Logstash**: Port 5044 for log ingestion
- **Filebeat**: Automatically ships application logs

### Built-in Monitoring
- **Fiber Monitor**: `http://localhost:4000/dashboard`
- **Application Logs**: `./logs/chat_message.log`
- **Structured logging** with request tracing

## API Rate Limiting

All API endpoints are protected with rate limiting:
- **Limit**: 50 requests per minute per IP address
- **Scope**: Applied to all `/api/*` endpoints
- **Reset**: Automatic reset every minute

## Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT Authentication**: Separate access and refresh tokens
- **Token Expiration**: Configurable token lifetimes
- **Session Management**: Secure session storage and cleanup
- **Input Validation**: Comprehensive request validation
- **CORS Protection**: Built-in Fiber security middleware

## Usage Examples

### 1. User Registration and Authentication Flow
```bash
# Register new user
curl -X POST http://localhost:4000/api/user/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","full_name":"Test User"}'

# Login user
curl -X POST http://localhost:4000/api/user/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Get message history (requires token)
curl -X GET http://localhost:4000/api/message/v1/history \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 2. WebSocket Chat Integration
```html
<!DOCTYPE html>
<html>
<head>
    <title>Go Chat App</title>
</head>
<body>
    <div id="messages"></div>
    <input type="text" id="messageInput" placeholder="Type your message...">
    <button onclick="sendMessage()">Send</button>

    <script>
        const socket = new WebSocket('ws://localhost:8080/message/v1/send');
        
        socket.onmessage = function(event) {
            const message = JSON.parse(event.data);
            document.getElementById('messages').innerHTML += 
                `<p><strong>${message.from}:</strong> ${message.message}</p>`;
        };
        
        function sendMessage() {
            const input = document.getElementById('messageInput');
            socket.send(JSON.stringify({
                from: 'currentUser',
                message: input.value,
                date: new Date().toISOString()
            }));
            input.value = '';
        }
    </script>
</body>
</html>
```
