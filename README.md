# TapInvest Wishlist API

A production-ready Go REST API for managing bond wishlists, built with Gin, pgxpool, and PostgreSQL.

## Prerequisites

- Go 1.20+
- PostgreSQL
- Postman (for testing)

## Setup

1. **Database Setup**
   Ensure PostgreSQL is running. Create a database named `tapinvest_api`.
   ```bash
   psql -U postgres -c "CREATE DATABASE tapinvest_api;"
   ```
   Execute the schema file to create the necessary tables:
   ```bash
   psql -U postgres -d tapinvest_api -f db/schema.sql
   ```
   *Optional: Insert some sample bonds into `master_data`.*

2. **Environment Configuration**
   Copy the example environment file and adjust if necessary:
   ```bash
   cp .env.example .env
   ```

3. **Install Dependencies**
   ```bash
   go mod tidy
   ```

4. **Run the API**
   ```bash
   go run main.go
   ```
   The API will start on `http://localhost:8080`.

## Testing with Postman

1. Open Postman.
2. Click **Import** and select the `postman_collection.json` file located in this repository.
3. You will have access to all endpoints grouped into folders. Test the endpoints sequentially:
   - Health Check
   - Get All Bonds
   - Wishlists (Create, Get All, Get Detail, Rename, Delete)
   - Wishlist Items (Add Bond, Remove Bond)

## Architecture Details

- **Gin Framework**: Used for routing and handling HTTP requests.
- **pgxpool**: Used for robust connection pooling and executing parameterized SQL queries to prevent SQL injection.
- **Thin Handlers**: Controllers (`handlers/`) focus strictly on HTTP request/response parsing and validation, delegating data access logic to `repository/`.
- **Validation**: Strict boundary rules are enforced:
  - Max 5 wishlists allowed.
  - Max 10 bonds per wishlist.
  - Wishlist names must be unique and <= 25 characters.
  - Duplicate bonds inside the same wishlist are rejected.
  - When a wishlist is deleted, its mappings are automatically removed via `ON DELETE CASCADE`.
