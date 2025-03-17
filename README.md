Overview

This is a scalable, efficient, and modular RESTful CRUD API built in Go to manage a library system. It allows:

    Users: Register, log in, and get JWT tokens for authentication.
    Books: Create, read, update, and delete books with caching for reads.
    Deployment: Runs in Docker with a mounted codebase, connecting to a remote MySQL database and Redis for caching.

The app is optimized for performance with connection pooling, indexed queries, async logging, and unit-tested services.
Architecture

The system uses a layered, modular design:

    Layers:
        Presentation: HTTP handlers (controllers) and routes (routes).
        Service: Business logic (services) for users and books.
        Repository: Database access (repositories) with GORM.
        Middleware: Security (JWT) and rate limiting.
        Database: Remote MySQL with connection pooling.
        Cache: Redis for fast reads.
        Logger: Async logging with goroutines.
    Scalability: Supports multiple instances, caching, and rate limiting.
    Efficiency: Pagination, indexed fields, and non-blocking logs.
    Modularity: Dependency injection via constructors (New* functions).

Data Flow

    Request: Client hits an endpoint (e.g., GET /books).
    Middleware: Checks JWT (if protected) and rate limits.
    Controller: Parses request and calls service.
    Service: Checks cache, runs logic, and interacts with repository.
    Repository: Queries MySQL (or cache updates).
    Response: JSON sent back, logged asynchronously.

File Structure
text
library-api/
├── main.go
├── config/
│ └── config.go
├── controllers/
│ └── user.go
│ └── book.go
├── services/
│ └── user.go
│ └── user_test.go
│ └── book.go
│ └── book_test.go
├── repositories/
│ └── user.go
│ └── book.go
├── models/
│ └── models.go
├── middleware/
│ └── auth.go
│ └── rate_limit.go
├── database/
│ └── db.go
├── cache/
│ └── redis.go
├── logger/
│ └── logger.go
├── Dockerfile
└── docker-compose.yml
Key Components
Database Tables (Remote MySQL)

    users:
        Columns: id (BIGINT, PK), created_at, updated_at, deleted_at, username (VARCHAR, unique index), password (VARCHAR, hashed).
        Purpose: User authentication.
    books:
        Columns: id (BIGINT, PK), created_at, updated_at, deleted_at, title (VARCHAR, index), author (VARCHAR, index).
        Purpose: Book management.

Endpoints

    Public:
        POST /users: Create a user.
        POST /login: Login and get JWT.
        GET /books: List books (cached).
        GET /books/{id}: Get a book.
    Protected (JWT):
        POST /books: Create a book.
        PUT /books/{id}: Update a book.
        DELETE /books/{id}: Delete a book.

Technologies

    Go: Core language.
    GORM: MySQL ORM.
    Gorilla Mux: Routing.
    JWT: Authentication (jwt-go).
    Bcrypt: Password hashing.
    Redis: Caching (go-redis).
    Testify: Unit testing.
    Docker: Containerization.

Endpoints

    Public:
        GET /books: List all books.
        GET /books/{id}: Get a specific book.
        POST /users: Create a new user (signup).
        POST /login: Log in and get a JWT token.
    Protected (JWT required):
        POST /books: Add a new book.
        PUT /books/{id}: Update a book.
        DELETE /books/{id}: Delete a book.

Technologies

    Go: Language for the API.
    GORM: ORM for MySQL interaction.
    Gorilla Mux: Router for handling HTTP requests.
    JWT: Token-based authentication (jwt-go).
    Bcrypt: Password hashing (golang.org/x/crypto/bcrypt).
    Docker: Containerizes the app.

Setup Instructions

    Remote MySQL:
        Ensure your remote MySQL server has a library database with users and books tables (GORM creates them via AutoMigrate if missing).
        Update DB_HOST, DB_PORT, DB_USER, DB_PASSWORD in docker-compose.yml with your remote details.
    Build and Run:
    bash

docker-compose up --build
Test Endpoints:

    Signup: POST /users
    bash

curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"username": "admin", "password": "mypassword123"}'
Login: POST /login
bash
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username": "admin", "password": "mypassword123"}'
Create Book (Protected): POST /books
bash

        curl -X POST http://localhost:8080/books -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d '{"title": "1984", "author": "George Orwell"}'

Notes

    Security: Use a .env file for production to avoid hard-coding secrets.
    Development: Mounted code (./:/app) allows live edits; restart with docker-compose restart app to apply changes (or use air for auto-reload).
    Production: Remove volumes and use CMD ["./main"] in Dockerfile after building.

endpoint docs

1. POST /users (Create a User)

   Description: Registers a new user by adding them to the users table with a hashed password.
   Public: No JWT required.
   Request:
   bash

curl -X POST http://localhost:8080/users \
 -H "Content-Type: application/json" \
 -d '{"username": "john_doe", "password": "mypassword123",,"email":"test@test.com"}'
Response (201 Created):
json

    {
        "id": 1,
        "created_at": "2025-02-26T12:00:00Z",
        "updated_at": "2025-02-26T12:00:00Z",
        "deleted_at": null,
        "username": "john_doe",
        "password": "$2a$10$..."  // Hashed password
    }
    Notes: Run this first to create a user for testing login.

2. POST /login (Login and Get JWT Token)

   Description: Verifies username/password and returns a JWT token for authentication.
   Public: No JWT required.
   Request:
   bash

curl -X POST http://localhost:8080/login \
 -H "Content-Type: application/json" \
 -d '{"username": "john_doe", "password": "mypassword123"}'
Response (200 OK):
json

    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    Notes: Copy the token value (e.g., eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...) and use it in the Authorization header for protected endpoints below. Replace <token> with this value in later commands.

3. GET /books (List All Books)

   Description: Retrieves all books from the books table.
   Public: No JWT required.
   Request:
   bash

curl -X GET http://localhost:8080/books \
 -H "Content-Type: application/json"
Response (200 OK):
json

    [
        {
            "id": 1,
            "created_at": "2025-02-26T12:05:00Z",
            "updated_at": "2025-02-26T12:05:00Z",
            "deleted_at": null,
            "title": "1984",
            "author": "George Orwell"
        },
        // More books if they exist
    ]
    Notes: Returns an empty array ([]) if no books are added yet.

4. GET /books/{id} (Get a Specific Book)

   Description: Retrieves a single book by its ID.
   Public: No JWT required.
   Request (e.g., ID = 1):
   bash

curl -X GET http://localhost:8080/books/1 \
 -H "Content-Type: application/json"
Response (200 OK):
json

    {
        "id": 1,
        "created_at": "2025-02-26T12:05:00Z",
        "updated_at": "2025-02-26T12:05:00Z",
        "deleted_at": null,
        "title": "1984",
        "author": "George Orwell"
    }
    Notes: Returns an empty object ({}) if the ID doesn’t exist.

5. POST /books (Create a Book)

   Description: Adds a new book to the books table.
   Protected: Requires JWT token.
   Request:
   bash

curl -X POST http://localhost:8080/books \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer <token>" \
 -d '{"title": "1984", "author": "George Orwell"}'
Response (200 OK):
json
{
"id": 1,
"created_at": "2025-02-26T12:05:00Z",
"updated_at": "2025-02-26T12:05:00Z",
"deleted_at": null,
"title": "1984",
"author": "George Orwell"
}
Error (401 Unauthorized without token):
json

    "Missing token"
    Notes: Replace <token> with the token from /login.

6. PUT /books/{id} (Update a Book)

   Description: Updates an existing book by ID.
   Protected: Requires JWT token.
   Request (e.g., ID = 1):
   bash

curl -X PUT http://localhost:8080/books/1 \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer <token>" \
 -d '{"title": "1984 Updated", "author": "George Orwell"}'
Response (200 OK):
json
{
"id": 1,
"created_at": "2025-02-26T12:05:00Z",
"updated_at": "2025-02-26T12:10:00Z",
"deleted_at": null,
"title": "1984 Updated",
"author": "George Orwell"
}
Error (401 Unauthorized without token):
json

    "Missing token"

7. DELETE /books/{id} (Delete a Book)

   Description: Soft-deletes a book by ID (sets deleted_at).
   Protected: Requires JWT token.
   Request (e.g., ID = 1):
   bash

curl -X DELETE http://localhost:8080/books/1 \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer <token>"
Response (200 OK):
json
"Book deleted"
Error (401 Unauthorized without token):
json

    "Missing token"
    Notes: GORM soft-deletes (marks deleted_at), so the row stays in the database but is hidden from normal queries.

How to Use These

    Start the App:
    bash

docker-compose up --build

    Ensure your remote MySQL is configured in docker-compose.yml with DB_HOST, DB_PORT, etc.

Step-by-Step Testing:

    Create a User: Run POST /users first.
    Get a Token: Use POST /login with the same credentials.
    Test Protected Routes: Use the token in the Authorization: Bearer <token> header for POST /books, PUT /books/{id}, and DELETE /books/{id}.
    Test Public Routes: GET /books and GET /books/{id} work without a token.

Example Sequence:
bash

    # Create user
    curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"username": "john_doe", "password": "mypassword123"}'

    # Login to get token
    curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username": "john_doe", "password": "mypassword123"}'
    # Copy the token (e.g., "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

    # Create a book with token
    curl -X POST http://localhost:8080/books -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." -d '{"title": "1984", "author": "George Orwell"}'

    # List books
    curl -X GET http://localhost:8080/books -H "Content-Type: application/json"

Notes

    Token: Replace <token> with the actual JWT from /login. It’s long and starts with eyJ....
    Remote MySQL: Ensure your DB_* environment variables in docker-compose.yml match your remote MySQL setup.
    Errors: If you get 401 Unauthorized, check the token. If connection refused, verify MySQL details.

curl -X POST http://localhost:8080/bulk-emails -H "Authorization: Bearer token"
