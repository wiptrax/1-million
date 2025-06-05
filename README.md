# Golang Server for 1 Million DB Inserts

This project provides a Golang server designed to benchmark and test database performance by performing 1 million insert operations. It's ideal for simulating high load and evaluating system behavior under stress.

## Features

-   **HTTP Server**: Triggers 1 million inserts via an API endpoint.
-   **Bulk Inserts**: Optimized for high-volume database inserts.
-   **Database Compatibility**: Easily adaptable to PostgreSQL, MySQL, SQLite, etc.
-   **Metrics**: Tracks insertion performance.
-   **Parallelism**: Supports concurrent inserts for improved throughput.

## Prerequisites

-   Go 1.16+
-   A running database (e.g., PostgreSQL, MySQL, SQLite).
-   Required Go database driver (e.g., `github.com/jackc/pgx/v5` for PostgreSQL).

### Example Database Schema (PostgreSQL)

```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    insert_time_milli BIGINT NOT NULL
);
```

## Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/wiptrax/1-million.git
    cd 1-million
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

## Configuration

Adjust the following in `main.go`:

-   **`dbConn`**: Your database connection string.
-   **`insertBatchSize`**: Number of records inserted per batch.
-   **`concurrency`**: Number of concurrent insert operations.

## Running the Server

1.  **Build the application:**
    ```bash
    go build -o db-insert-server
    ```

2.  **Run the server:**
    ```bash
    ./db-insert-server
    ```
    The server will run on `http://localhost:8080` by default.

3.  **Start Insert Operation:**
    Trigger the inserts by making a GET request to the `/insert` endpoint:
    ```bash
    curl http://localhost:8080/insert
    ```

**Example Output:**

```json
{
  "status": "success",
  "message": "Successfully inserted 1000000 records into the database."
}
```

## Testing the Inserts

Use tools like [k6](https://k6.io/) to simulate load.

**Example `load-test.js` for k6:**

```javascript
import http from 'k6/http';
import { check } from 'k6';

export default function () {
  let res = http.get('http://localhost:8080/insert');
  check(res, {
    'is status 200': (r) => r.status === 200,
  });
}
```

**Run k6 test:**

```bash
k6 run load-test.js
```

## Troubleshooting

-   **Database connection issues**: Verify `dbConn` and database status.
-   **Slow insertions**: Optimize schema, adjust batch size, or concurrency.
-   **Insufficient resources**: Increase server resources or distribute load.

## License

This project is licensed under the [MIT License](LICENSE).

## Folder Structure

```
.
├── cmd
│   └── server
│       └── main.go        # Main server entry point
├── internal
│   └── db
│   |   └── db.go      # DB client and insert logic
|   └──sqlc
|       ├── db.go
|       ├── model.go
|       └── user.sql.go
├── go.mod
└── README.md
```