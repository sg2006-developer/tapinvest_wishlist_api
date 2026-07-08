# TapInvest Wishlist Backend - Implementation Plan

## Objective

Build a production-style Go backend for the Wishlist feature using the
existing PostgreSQL schema **without changing any table names**.

Tables: - `master_data` - `wish_lists` - `wish_isin`

Use: - Gin - pgx - PostgreSQL - godotenv

------------------------------------------------------------------------

# Required APIs

## Bonds

### GET /bonds

Returns all bonds from `master_data`.

Future ready: - search - pagination

------------------------------------------------------------------------

## Wishlists

### POST /wishlists

Rules: - Max 5 wishlists - Name required - Max 25 characters - Reject
duplicate names - Trim spaces - Return 201

------------------------------------------------------------------------

### GET /wishlists

Return all wishlists.

Also return:

-   wishlistId
-   wishlistName
-   bondCount

Use SQL COUNT() with JOIN.

------------------------------------------------------------------------

### GET /wishlists/{wishlistId}

Return

wishlist information

plus

all bonds inside wishlist.

Join

wish_lists

wish_isin

master_data

------------------------------------------------------------------------

### PUT /wishlists/{wishlistId}

Rename wishlist.

Validation

-   unique
-   \<=25 chars
-   non-empty

------------------------------------------------------------------------

### DELETE /wishlists/{wishlistId}

Delete wishlist.

Foreign key cascade removes mappings.

Return success.

------------------------------------------------------------------------

# Wishlist Items

### POST /wishlists/{wishlistId}/items

Request

{ "isin":"..." }

Validation order

1.  Wishlist exists
2.  Bond exists
3.  Bond already exists?
4.  Wishlist has 10 bonds?
5.  Insert

Return 201.

------------------------------------------------------------------------

### DELETE /wishlists/{wishlistId}/items/{bondId}

Delete mapping only.

Do not delete master_data.

------------------------------------------------------------------------

# Architecture

Client

↓

Gin Router

↓

Handler

↓

Repository

↓

PostgreSQL

Handlers: - Parse request - Validate basic input - Call repository -
Return response

Repository: - SQL only - Parameterized queries (\$1, \$2)

Config: - PostgreSQL connection

Models: - Bond - Wishlist - WishIsin

------------------------------------------------------------------------

# Folder Structure

tapinvest_api/

-   main.go
-   go.mod
-   .env

config/ - database.go

routes/ - routes.go

handlers/ - health_handler.go - bond_handler.go - wishlist_handler.go

repository/ - bond_repository.go - wishlist_repository.go

models/ - bond.go - wishlist.go - wishisin.go

utils/ - response.go - validation.go

------------------------------------------------------------------------

# Coding Standards

-   Thin handlers
-   SQL only in repository
-   No hardcoded credentials
-   Environment variables
-   Parameterized SQL
-   Proper HTTP status codes
-   Logging
-   Consistent JSON responses

------------------------------------------------------------------------

# Response Wrapper

Use one helper.

Success

{ "success": true, "message": "...", "data": {} }

Failure

{ "success": false, "message": "...", "error": {} }

This can be changed later if a different response format is provided.

------------------------------------------------------------------------

# Validation Rules

Wishlist

-   max 5
-   unique
-   \<=25 chars
-   trim whitespace
-   reject empty

Wishlist Items

-   max 10 bonds
-   duplicate not allowed
-   bond may exist in multiple wishlists

------------------------------------------------------------------------

# SQL Queries

Need repository methods for

-   GetAllBonds
-   CreateWishlist
-   GetWishlists
-   GetWishlistDetail
-   RenameWishlist
-   DeleteWishlist
-   AddBond
-   RemoveBond

Always use parameterized SQL.

------------------------------------------------------------------------

# Future Features

Keep code extensible for

-   search
-   pagination
-   sold out bonds
-   expired bonds
-   authentication
-   unit tests

------------------------------------------------------------------------

# Development Order

1.  Database connection
2.  Health API
3.  GET /bonds
4.  GET /wishlists
5.  POST /wishlists
6.  GET wishlist detail
7.  PUT wishlist
8.  DELETE wishlist
9.  POST wishlist item
10. DELETE wishlist item
11. Validation
12. Response helper
13. README
14. Postman collection
