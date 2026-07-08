# Technical Architecture & Codebase Guide
*A comprehensive deep-dive for technical managers to understand the structure, routing, and design decisions of the TapInvest Wishlist API.*

---

## 1. Directory Structure & File Purposes

The codebase strictly follows a **Clean Architecture** (specifically the Controller-Repository pattern). This ensures a strong separation of concerns: routing is separated from HTTP parsing, which is separated from database logic.

```text
Tap_invest_api2/
├── main.go                  # Application Entry Point
├── .env                     # Environment variables (DB credentials)
├── db/
│   ├── schema.sql           # Original DB schema backup
│   └── new_schema.sql       # Current schema used for the 'wish' database
├── models/                  # Domain Entities & DTOs
│   ├── bond.go
│   ├── wishisin.go
│   └── wishlist.go
├── routes/
│   └── routes.go            # Router & Endpoint mappings
├── handlers/                # HTTP Controllers (Request/Response)
│   ├── bond_handler.go
│   ├── wishlist_handler.go
│   └── health_handler.go
└── repository/              # Data Access Layer (SQL)
    ├── bond_repository.go
    └── wishlist_repository.go
```

### Breakdown of Folders & Files:

#### 1. `main.go`
- **Purpose**: The absolute entry point of the app.
- **What it does**: It loads the `.env` file, initializes the PostgreSQL connection pool using `pgxpool`, passes that pool into the Repositories, injects those Repositories into the Handlers, mounts the `routes.go` router, and starts the HTTP server.

#### 2. `routes/routes.go`
- **Purpose**: The "Switchboard" of the API.
- **What it does**: This is where you can see exactly which URL maps to which function. For example, it defines that `PATCH /api/v1/wishlist/:wishlistId/bond/:bondIsin/color` is handled by `WishlistHandler.SetBondColor`. **(If your manager asks "Where is the API defined?", point them here).**

#### 3. `handlers/` (The HTTP Layer)
- **Purpose**: Input validation and output formatting.
- **What it does**: Files here (`bond_handler.go` and `wishlist_handler.go`) take an incoming HTTP request, parse the JSON body or URL parameters (e.g., verifying a UUID is valid), and call the Repository. They contain **zero SQL**. They just format the repository's result into the `{ "data": ... }` JSON envelope and return standard HTTP status codes (200, 400, 404).

#### 4. `repository/` (The Database Layer)
- **Purpose**: Raw data retrieval and persistence.
- **What it does**: Files here (`bond_repository.go` and `wishlist_repository.go`) contain all the actual SQL queries. Handlers call these files, and these files talk to Postgres. 

#### 5. `models/`
- **Purpose**: Blueprints.
- **What it does**: Contains Go `structs`. These structs dictate the exact JSON `camelCase` keys that the frontend receives (e.g., `bondYield`) and provide the memory structure for scanning SQL rows.

---

## 2. Request Lifecycle (How Data Flows)

If your manager asks **"How does the code work when a request comes in?"**, explain this flow using the *Bulk Reorder* endpoint as an example:

1. **Client Sends Request**: `PATCH /api/v1/wishlist/123e4567.../reorder` with a JSON array of ISINs.
2. **Router (`routes.go`)**: Sees `PATCH .../reorder` and forwards it to `WishlistHandler.ReorderBonds`.
3. **Handler (`wishlist_handler.go`)**: 
   - Validates the `wishlistId` is a real UUID.
   - Binds the JSON body to a Go struct to ensure the array isn't empty.
   - Calls `repository.ReorderBonds(id, isins)`.
4. **Repository (`wishlist_repository.go`)**:
   - Opens a Database Transaction (`BEGIN`).
   - Checks if all ISINs exist and if the count matches the DB exactly.
   - Loops through the array and runs `UPDATE` SQL statements to set the `position` integer.
   - Commits the transaction (`COMMIT`) and returns success to the handler.
5. **Handler**: Sends a `200 OK` back to the client.

---

## 3. Key Technical Decisions to Highlight

When having a technical discussion, these are the architectural wins you should mention:

1. **Why `pgx` instead of an ORM (like GORM)?**
   - We chose to write raw SQL using the `pgxpool` driver because it is significantly faster and gives us absolute control over complex `ORDER BY` clauses. ORMs struggle with complex tie-breaker sorting (like `is_pinned DESC, color ASC NULLS LAST, position ASC`), but raw SQL handles it flawlessly.

2. **Database Transactions (`tx.Begin()`)**
   - For endpoints that modify multiple rows at once (like Bulk Reordering bonds), we wrap the SQL in a transaction. If the server crashes on the 5th bond update out of 10, the database rolls back the entire request, ensuring the wishlist's order is never corrupted.

3. **Trigram Fuzzy Search (`pg_trgm`)**
   - In `bond_repository.go`, the search function doesn't just do basic `LIKE` matching. We enabled the `pg_trgm` extension in Postgres. This breaks search terms into 3-letter combinations (n-grams) to calculate a `similarity()` score. It makes the API typo-tolerant (searching "HDFX" finds "HDFC") and ranks the best matches first.

4. **Security via UUIDs**
   - We migrated wishlist primary keys from auto-incrementing integers (`1, 2, 3`) to UUIDs (`123e4567-e89b...`). This prevents Insecure Direct Object Reference (IDOR) vulnerabilities, as attackers cannot simply guess the next wishlist ID to scrape data.
