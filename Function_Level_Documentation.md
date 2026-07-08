# Function-Level Codebase Documentation

This document provides a granular, file-by-file breakdown of the core functions powering the TapInvest backend API. It is designed for developers or technical managers who need to understand exactly what each function in the Handlers and Repositories is doing.

---

## 1. Handlers (`/handlers` directory)
*Handlers are responsible for receiving HTTP requests, validating inputs (like checking if a UUID is valid), parsing JSON bodies, calling the database repository, and sending back the final HTTP JSON response.*

### `wishlist_handler.go`
- **`CreateWishlist(c *gin.Context)`**: Parses the `{ "name": "..." }` body. Returns HTTP 400 if empty. Calls `repo.Create()`. Returns HTTP 409 if the name already exists or HTTP 422 if the user already has 5 wishlists.
- **`GetWishlists(c *gin.Context)`**: A simple GET request that calls `repo.GetAll()` and wraps the resulting list in a standard `{ "data": [...], "success": true }` JSON response.
- **`GetWishlist(c *gin.Context)`**: Extracts the `wishlistId` from the URL and validates it is a correct UUID. Extracts the `sortBy` query parameter (defaulting to "manual"). Calls `repo.GetByID()` and returns the fully populated wishlist and its sorted bonds.
- **`RenameWishlist(c *gin.Context)`**: Validates the UUID, parses the new name from the JSON body, and calls `repo.Update()`.
- **`DeleteWishlist(c *gin.Context)`**: Validates the UUID and calls `repo.Delete()`. Returns a 204 No Content on success.
- **`AddBond(c *gin.Context)`**: Parses the `{ "bondIsin": "..." }` body. Calls `repo.AddBond()`. Returns HTTP 409 if the bond is already in the list, HTTP 422 if the list is full (10 bonds), or HTTP 404 if the bond ISIN doesn't exist.
- **`SetBondColor`, `SetBondPosition`, `SetBondPin`**: Extracts the specific field from the JSON body and calls the respective atomic update function in the repository.
- **`ReorderBonds(c *gin.Context)`**: Parses an array of ISIN strings from the body. Calls `repo.ReorderBonds()`. Returns a 400 Bad Request if the array length doesn't perfectly match the database state.

### `bond_handler.go`
- **`GetBonds(c *gin.Context)`**: Extracts `sortBy` and `sortOrder` from the URL. Defaults to `isin` and `asc`. Passes them to `repo.GetAll()`.
- **`SearchBonds(c *gin.Context)`**: Extracts the `q` (query) string from the URL. If the query is empty, it short-circuits and immediately returns an empty array `[]` without hitting the database. Otherwise, it calls `repo.Search()`.

---

## 2. Repositories (`/repository` directory)
*Repositories are strictly for interacting with PostgreSQL. They contain raw SQL queries executed via the `pgx` driver. They do not know what HTTP is.*

### `wishlist_repository.go`
- **`Create(ctx, name)`**: Runs two SQL `COUNT` queries first: one to ensure the user has <5 wishlists, and one to ensure the new name is unique (using a case-insensitive `LOWER()` check). It then generates a new `uuid.New()` and runs an `INSERT` statement.
- **`GetAll(ctx)`**: Runs a `LEFT JOIN` between `wish_lists` and `wish_isin`. It uses `GROUP BY` to dynamically count how many bonds are inside each wishlist (`bond_count`), returning them ordered by creation date.
- **`GetByID(ctx, id, sortBy)`**: This is the most complex read function. First, it fetches the wishlist metadata. Then it runs a massive `switch` statement against the `sortBy` parameter to generate a dynamic SQL `ORDER BY` clause. Finally, it runs a `JOIN` to pull all bond data from `master_data`, sorting it directly in the database (always placing `is_pinned` at the top).
- **`AddBond(ctx, wishlistID, isin)`**: Runs multiple validation queries:
  1. Does the wishlist exist?
  2. Does the bond exist in `master_data`?
  3. Are there fewer than 10 bonds currently in this list?
  4. Is this ISIN already in this list?
  5. What is the current `MAX(position)`?
  After all checks pass, it runs an `INSERT` into `wish_isin`, assigning it `position = maxPos + 1` so it drops to the bottom of the list.
- **`ReorderBonds(ctx, wishlistID, bondIsins)`**: Highly critical function. It opens a database transaction (`tx, err := r.db.Begin(ctx)`). It counts the rows in the database and ensures it exactly matches `len(bondIsins)`. It then loops over the provided array and executes an `UPDATE wish_isin SET position = i` for every bond. If anything fails, it rolls back. If successful, it commits.

### `bond_repository.go`
- **`GetAll(ctx, sortBy, sortOrder)`**: Uses a `switch` statement to map the user's `sortBy` string to actual database column names (to prevent SQL injection). It dynamically builds the `ORDER BY` clause, ensuring that `NULLS LAST` is appended to numerical fields like Yield so empty data doesn't incorrectly sort to the top.
- **`Search(ctx, query)`**: Executes the fuzzy search. It runs a SQL query utilizing `ILIKE` for exact or partial matches, AND `similarity(bond_name, $2) > 0.2` to utilize the `pg_trgm` trigram extension for typo tolerance. It dynamically sorts the results by their similarity score so the best matches appear first.
