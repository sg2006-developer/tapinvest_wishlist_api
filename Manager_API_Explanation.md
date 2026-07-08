# Architecture & Implementation Overview
*A high-level summary of the TapInvest backend architecture and design decisions.*

## 1. What is the tech stack?
The backend is built in **Go (Golang)** using the **Gin Web Framework**. We chose Gin because it's extremely fast, lightweight, and perfect for building high-performance REST APIs. 

For the database, we are using **PostgreSQL** connected via the **pgxpool** driver. Instead of using a bulky ORM like GORM, we use `pgx` directly to write highly optimized, raw SQL queries. This gives us maximum performance and control, especially for complex sorting operations.

## 2. How is the code organized? (The Flow)
The application follows a standard **Controller-Repository** architectural pattern:
1. **`routes/routes.go`**: This is the "Traffic Cop." It maps incoming URLs (like `GET /api/v1/bond`) to specific handler functions.
2. **`handlers/`**: These are the "Translators." They take the HTTP request, parse the JSON or query parameters, call the Repository for data, and then format the response back into JSON for the frontend.
3. **`repository/`**: This is the "Database layer." All SQL queries live here. Handlers do not talk to the database directly; they only ask the Repository to fetch or save data.
4. **`models/`**: These are the "Blueprints." They define the shapes of our Go structs (like `Bond` or `Wishlist`) so the backend knows what the data looks like.

## 3. Why did we change Wishlist IDs to UUIDs?
Previously, wishlists used simple integer IDs (1, 2, 3...). We migrated these to **UUIDs** (Universally Unique Identifiers) for security and scalability.
- **Security**: Integer IDs make it easy for malicious users to scrape data by simply guessing IDs (e.g., `/wishlist/4`, `/wishlist/5`). UUIDs are cryptographically random and impossible to guess.
- **Scalability**: If we ever decide to use a distributed database, UUIDs prevent ID collisions across multiple servers.

## 4. How does the Fuzzy Search work?
We implemented a true fuzzy search for bonds (`/api/v1/bond/search`). 
Instead of relying on basic substring matching, we enabled the **`pg_trgm`** (Trigram) extension inside PostgreSQL. 
- **Why?** It breaks words down into 3-letter chunks (trigrams) to measure string similarity. If a user mistypes "HDFX" instead of "HDFC", the system calculates a similarity score and still returns the correct bond. This dramatically improves the User Experience (UX).

## 5. How does Sorting & Pinning work?
We do all sorting directly in the database using SQL `ORDER BY` clauses rather than sorting lists in memory (which is much faster).
- When a user asks to sort a wishlist by "Yield", the database query looks like this: `ORDER BY is_pinned DESC, yield DESC, position ASC`.
- **The Rule**: `is_pinned DESC` is *always* injected as the first sorting priority. This guarantees that pinned items always float to the top of the frontend UI, regardless of what other sorting mode the user clicks on.

## 6. How does Drag-and-Drop Reordering work?
Instead of the frontend making 10 separate API calls when a user drags a bond to a new position, we built a **Bulk Reorder Endpoint** (`PATCH /api/v1/wishlist/:id/reorder`).
- The frontend sends one array of ISINs in their new desired order.
- The backend wraps the entire update in a **Database Transaction**. It updates the `position` of all bonds simultaneously.
- **Why?** If the server crashes halfway through updating positions, the transaction rolls back. This guarantees the list order never gets corrupted or ends up with duplicate positions.
