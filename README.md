# Sensory LogClient

This repository contains the `logClient` package, a logging client designed for the Sensory ecosystem. It provides functionality to send structured log entries over a TCP connection to a logging server and persist them directly into a MongoDB database.

## Description

The `logClient` package acts as a robust client for applications requiring centralized logging within the Sensory platform. It enables applications to:
*   Establish a TCP connection to a designated log server.
*   Authenticate and retrieve project metadata using an API key from a database.
*   Construct detailed log entries including log level, timestamp, project and service identifiers, the log message, and runtime information (file, line, function).
*   Send these structured logs over the established TCP connection.
*   Concurrently save log entries directly into a MongoDB database for persistence and analysis.
*   Provide convenience methods for common log levels (INFO, ERROR, FATAL).
*   Offer a "no-operation" client for environments where logging should be disabled or mocked.

**What this tool can be used for:**
This client is essential for applications built within the Sensory ecosystem that need to:
*   Centralize log collection for monitoring and debugging purposes.
*   Persist application logs for auditing, analysis, and historical reference.
*   Integrate seamlessly with the Sensory platform's logging infrastructure.
*   Provide immediate insight into application behavior and potential issues through structured and contextualized log data.

## Features

*   **TCP-based Logging**: Sends structured log entries over a TCP connection to a configured address.
*   **Project Metadata Integration**: Automatically retrieves project-specific details (name, service, ID, user ID) from a database using an API key.
*   **Structured Log Entries**: Generates log entries containing `Level`, `Timestamp`, `Project`, `Service`, `Message`, and `Runtime` details (file, line, function).
*   **Database Persistence**: Automatically inserts log entries into a MongoDB collection for long-term storage.
*   **Log Level Utilities**: Dedicated methods for `INFO`, `ERROR`, and `FATAL` level logging.
*   **No-Operation Client**: Supports creating a dummy client that prints logs to standard output without sending them over TCP or persisting to the database.
*   **Runtime Context Capture**: Automatically captures the file, line number, and function name where the log originated.

## Installation

This package is part of the broader Sensory project. To use it, you will need to ensure your Go environment can access the required dependencies.

```bash
go get go.mongodb.org/mongo-driver/v2/bson
# The following sensory-specific packages are assumed to be available
# within your Go module's context, likely from the same monorepo.
# github.com/EsanSamuel/sensory/db
# github.com/EsanSamuel/sensory/helpers
# github.com/EsanSamuel/sensory/models
```

## Usage

### Initializing the Log Client

To create a new logging client, you need an API key and the TCP address of the log server.

```go
package main

import (
	"log"
	"github.com/EsanSamuel/sensory/LogClient" // Assuming correct module path
)

func main() {
	apiKey := "YOUR_API_KEY"
	logServerAddr := "localhost:8080" // Example address

	client, err := logClient.New(apiKey, logServerAddr)
	if err != nil {
		log.Fatalf("Failed to create log client: %v", err)
	}
	defer client.conn.Close() // Ensure the connection is closed
	
	// Use the client to send logs
	client.INFO("Application started successfully!")
}
```

### Creating a No-Op Client

For testing or disabling logging, you can create a no-operation client.

```go
package main

import (
	"github.com/EsanSamuel/sensory/LogClient"
	"fmt"
)

func main() {
	noOpClient := logClient.NewNoOp()

	// Logs sent via noOpClient will print to console but not be sent or persisted.
	noOpClient.INFO("This is a no-op info message.")
	noOpClient.ERROR("This is a no-op error message.")
	noOpClient.Send("DEBUG", "This is a no-op debug message sent via Send method.")
	
	fmt.Println("No-op client operations complete.")
}
```

### Sending Log Entries

You can send logs using specific level methods or the generic `Send` method.

```go
package main

import (
	"log"
	"github.com/EsanSamuel/sensory/LogClient"
	"fmt"
)

func main() {
	apiKey := "YOUR_API_KEY"
	logServerAddr := "localhost:8080"

	client, err := logClient.New(apiKey, logServerAddr)
	if err != nil {
		log.Fatalf("Failed to create log client: %v", err)
	}
	defer client.conn.Close()

	// Using convenience methods
	client.INFO("User 'John Doe' logged in.")
	client.ERROR(fmt.Errorf("database connection failed: %w", err))

	// Using the generic Send method
	client.Send("WARN", "High memory usage detected.")

	// FATAL logs will exit the application
	// client.FATAL("Critical system error, shutting down.")
}
```

## Folder Structure Explanation

```
Sensory/
└── LogClient/
    ├── client.go    # Core client logic: connection, project metadata retrieval,
    |                # log entry construction, TCP transmission, and database persistence.
    └── logger.go    # Provides convenience methods (INFO, ERROR, FATAL) for simpler logging.
```

## Technologies

*   **Go**: The primary programming language.
*   **TCP Networking**: Utilized for establishing connections and sending log data.
*   **MongoDB**: Used for persistent storage of log entries (via `go.mongodb.org/mongo-driver`).
*   **JSON Encoding**: For marshaling log entries into a structured format before transmission.

## License

Not specified.