# Sensory

## Description
Sensory is a Go-based RESTful API service designed to provide core functionalities for user and project management, API key generation, and comprehensive logging. It offers a robust backend solution for managing user registrations, project creation, securing projects with API keys, and handling operational logs, including retrieval and filtering capabilities.

## Features
*   **User Management**:
    *   Register new users.
    *   Retrieve details for a specific user.
*   **Project Management**:
    *   Create new projects.
    *   Generate unique API keys for existing projects.
*   **Logging System**:
    *   Retrieve all logs associated with a specific user.
    *   Filter user logs based on log levels.
    *   Retrieve a specific log entry by its ID.
*   **Welcome Endpoint**: A basic endpoint to confirm the API is operational.

## Installation

To set up and run Sensory locally, follow these steps:

1.  **Install Go**:
    Ensure you have Go version `1.25.5` or later installed. You can download it from [golang.org/dl](https://golang.org/dl/).

2.  **Clone the Repository**:
    Since the module path is `github.com/EsanSamuel/sensory`, it's expected that you would clone the repository from its source.
    ```bash
    git clone https://github.com/EsanSamuel/sensory.git
    cd sensory
    ```

3.  **Download Dependencies**:
    Navigate to the project directory and download the required Go modules.
    ```bash
    go mod tidy
    ```

4.  **Environment Variables**:
    Sensory uses `godotenv` to load environment variables from a `.env` file. Create a file named `.env` in the root of the project directory.
    ```
    # Example .env content (adjust as needed based on your application's requirements)
    # PORT=8000
    # MONGO_URI="mongodb://localhost:27017"
    # JWT_SECRET="your_jwt_secret_key"
    ```
    *Note: The actual required environment variables are not specified in the provided data, but the presence of `godotenv.Load(".env")` indicates their use.*

## Usage

To start the Sensory API server:

1.  **Run the Application**:
    From the project's root directory, execute the `main.go` file.
    ```bash
    go run main.go
    ```

2.  **Access the API**:
    The server will start on port `8000`. You can access the API endpoints at `http://localhost:8000`.

    **Example Endpoints:**
    *   `GET http://localhost:8000/hello`
    *   `POST http://localhost:8000/register`
    *   `GET http://localhost:8000/user/:userId`
    *   `POST http://localhost:8000/project`
    *   `POST http://localhost:8000/project/api_key/:projectId`
    *   `GET http://localhost:8000/logs/:userId`
    *   `GET http://localhost:8000/logs/level/:userId`
    *   `GET http://localhost:8000/log/:logId`

## Folder Structure Explanation

Based on the imports in `main.go`, the project likely adheres to a modular structure:

```
sensory/
├── main.go                       # The main application entry point.
├── go.mod                        # Go module definition file.
├── go.sum                        # Go module checksums.
├── .env                          # Environment variables file (if created).
├── LogServer/                    # Package containing log server functionalities.
│   └── (e.g., logserver.go)
└── controllers/                  # Package containing HTTP handler functions for API routes.
    └── (e.g., user_controller.go, project_controller.go, log_controller.go)
```

## Technologies

*   **Go**: The primary programming language (version `1.25.5`).
*   **Gin**: A high-performance HTTP web framework (`github.com/gin-gonic/gin`).
*   **Godotenv**: For loading environment variables from `.env` files (`github.com/joho/godotenv`).
*   **Go MongoDB Driver (indirect)**: Likely used for database interactions (`go.mongodb.org/mongo-driver/v2`).
*   **GORM (indirect)**: An ORM library for Go, suggesting potential relational database integration (`gorm.io/gorm`).

## Design and built by Esan Samuel