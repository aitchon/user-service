# User Service API

This project provides a REST API for managing user data. It supports operations for creating, updating, retrieving, and deleting users. The API is built using the Echo web framework and interacts with an SQLite database.

## Prerequisites

Before running the application, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.22 or higher)
- [SQLite](https://www.sqlite.org/download.html) for local database storage (or any database you're configuring in your project)
- [Swagger](https://swagger.io/tools/swagger-ui/) (optional for API documentation)

## Project Setup

### 1. Clone the Repository

Clone the repository to your local machine using the following command:

```bash
git clone https://github.com/aitchon/user-service.git
cd user-service
```

### 2. Install Dependencies
The project uses Go modules to manage dependencies. Run the following command to install the necessary packages:

```bash
go mod tidy
```

### 3. Set Up the Database
Ensure that you have a SQLite database set up for the application. If you're using SQLite locally, the database will be created automatically when you run the application.

If needed, you can create a new SQLite database file:
```bash
touch user_db/user_data.db
```

### 4. Build the Application
To build the application, use the following command:
```bash
go build -o user-service ./cmd
```
This will compile the Go code and create an executable named user-service.

### 5. Run the Application
To start the application, run:
```bash
./user-service
```
By default, the server will run on http://localhost:3002. Swagger docs can be viewed here: http://localhost:3002/swagger/index.html

### 6. API Endpoints
The application provides the following endpoints:

- POST /users - Create a new user.
- GET /users/{id} - Retrieve a user by ID.
- PUT /users/{id} - Update a user by ID.
- DELETE /users/{id} - Delete a user by ID.

The request body should be in JSON format. Here's an example:

Example Request: POST /users
```json
{
  "user_name": "john_doe",
  "email": "john_doe@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "status": "A",
  "department": "IT"
}
```
Example Response: 200 OK
```json
{
  "message": "User created successfully"
}
```

### 7. Swagger Documentation
To generate Swagger API documentation, follow these steps:

Install the Swagger tool:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate the Swagger documentation:

```bash
swag init --output ./docs --parseDepth 2 --parseDependency --dir ./cmd,./controllers
```

After generating the documentation, visit http://localhost:8080/swagger/index.html to view the interactive API documentation.

### 8. Running Tests
To run the tests, use the following command:

```bash
gingko ./tests
```
This will run the unit tests for the application.