# OTP-Based Login Service

A backend service built in Golang that implements OTP-based login and registration, JWT authentication, and basic user management features. The project is built with a clean, layered architecture (controllers, services, repository) and is fully containerized with Docker.

## Features

* **OTP Login/Registration**: Secure, passwordless login using One-Time Passwords.
* **JWT Authentication**: Protected endpoints use JSON Web Tokens for authentication.
* **In-Memory Rate Limiting**: Prevents OTP spam by limiting requests (3 requests per 10 minutes per phone number).
* **User Management**: REST endpoints to list users and retrieve a single user by ID or phone number.
* **Pagination & Search**: The user list endpoint supports pagination and searching by phone number.
* **Swagger API Documentation**: Fully documented API with an interactive Swagger UI.
* **Containerized**: Includes a `Dockerfile` and `docker-compose.yml` for easy setup and deployment.

## Technology Stack

* **Language**: Go
* **Web Framework**: Gin
* **Database**: SQLite
* **Containerization**: Docker & Docker Compose
* **API Documentation**: Swagger (OpenAPI)

## API Documentation

This project uses Swagger for interactive API documentation. Once the application is running, you can access the full documentation and test the endpoints at:

[**http://localhost:8080/swagger/index.html**](http://localhost:8080/swagger/index.html)


---
## Running the Application

You can run the application either locally with Go or using Docker.

### How to Run Locally

#### Prerequisites
* Go (version 1.21 or later) installed.
* A C compiler (like GCC/MinGW on Windows) for the `go-sqlite3` driver.

#### Instructions
1.  **Clone the repository**:
    ```bash
    git clone https://github.com/AryaTabani/otp-login-service.git
    cd otp-login-service
    ```

2.  **Install dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Set Environment Variables** (Optional):
    The application uses a default JWT secret key for development. You can set your own for production.
    ```powershell
    # In PowerShell
    $env:JWT_SECRET_KEY="your-super-secret-key"
    
    # In Bash (Linux/macOS)
    export JWT_SECRET_KEY="your-super-secret-key"
    ```

4.  **Run the application**:
    ```bash
    go run main.go
    ```
    The server will start and be accessible at `http://localhost:8080`.

### How to Run with Docker

#### Prerequisites
* Docker and Docker Compose installed.

#### Instructions
1.  **Clone the repository**:
    ```bash
    git clone https://github.com/AryaTabani/otp-login-service.git
    cd otp-login-service
    ```

2.  **Build and run the container**:
    From the project's root directory, run:
    ```bash
    docker compose up --build
    ```
    This command will build the Go application image, start the service, and make it available at `http://localhost:8080`.

---
## Database Choice Justification

For this project, **SQLite** was chosen as the database for the following reasons:

* **Simplicity**: As a serverless, file-based database, SQLite requires no separate server process or complex configuration. This makes the application setup incredibly simple and fast.
* **Portability**: The entire database is contained in a single file within the project, making it highly portable and easy to share or deploy.
* **Ease of Use**: It uses standard SQL and is well-supported by Go's native `database/sql` package, which simplifies development and reduces the need for external drivers or libraries (beyond the CGO-based driver).

For the scope of this project, SQLite provides all the necessary functionality for user and OTP management without the overhead of a more complex client-server database like PostgreSQL or MySQL.

---
## Example API Requests & Responses

Here are some examples of how to interact with the API.

### 1. Request an OTP

```bash
curl -X POST http://localhost:8080/auth/request-otp \
-H "Content-Type: application/json" \
-d '{
    "phone_number": "09140000"
}'
```
**Success Response:** (The OTP will be printed in your server console)
```json
{
    "success": true,
    "message": "OTP generated and printed to console. It will expire in 2 minutes."
}
```

### 2. Verify OTP and Login/Register

_Check your application console for the generated OTP first._

```bash
curl -X POST http://localhost:8080/auth/verify-otp \
-H "Content-Type: application/json" \
-d '{
    "phone_number": "09140000",
    "otp": "123456"
}'
```
**Success Response:** (This token is used for protected endpoints)
```json
{
    "success": true,
    "message": "User logged in successfully",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

### 3. List Users (Protected)

```bash
curl -X GET "http://localhost:8080/users?page=1&limit=5" \
-H "Authorization: Bearer <your_jwt_token>"
```
**Success Response:**
```json
{
    "success": true,
    "data": {
        "limit": 5,
        "page": 1,
        "total": 1,
        "users": [
            {
                "id": 1,
                "phone_number": "09140000",
                "created_at": "2025-09-15T19:15:36.123Z"
            }
        ]
    }
}
```
