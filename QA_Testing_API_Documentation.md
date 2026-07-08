# QA & Testing Team - API Documentation (v1)

This document outlines the latest backend API specifications, including all endpoints, path changes, and business rules implemented. Please use this reference to update your testing scripts and automation suites.

## Global Changes
1. **Base URL Update**: All routes have been migrated to the `/api/v1` namespace.
2. **UUIDs for Wishlists**: `wishlistId` is no longer an integer. It is now a 36-character UUID string (e.g., `123e4567-e89b-12d3-a456-426614174000`). All endpoints referencing wishlists expect this UUID.
3. **CamelCase Keys**: JSON response and request body keys strictly follow camelCase (e.g., `bondYield`, `minInvestment`, `bondIsin`).

---

## 1. Bond Catalog & Search

### 1.1 Get All Bonds (with sorting)
- **Method**: `GET`
- **Endpoint**: `/api/v1/bond`
- **Query Params**: 
  - `sortBy` (Options: `isin`, `bondYield`, `minInvestment`, `tenure`, `rating`. Default: `isin`)
  - `sortOrder` (Options: `asc`, `desc`. Default: `asc`)
- **Note**: `NULLS LAST` is applied to fields like yield and minInvestment so missing data floats to the bottom.

### 1.2 Search Bonds (Fuzzy Search)
- **Method**: `GET`
- **Endpoint**: `/api/v1/bond/search?q={search_term}`
- **Behavior**: Uses PostgreSQL trigram similarity. It handles exact ISIN matching, partial name matching, and typo-tolerant fuzzy matching (e.g., "HDFX" will match "HDFC").
- **Empty Query**: If `q=` is empty, it returns `[]` immediately.

---

## 2. Wishlist Management

### 2.1 Create Wishlist
- **Method**: `POST`
- **Endpoint**: `/api/v1/wishlist`
- **Body**: `{ "name": "string" }`
- **Rules**: Max 50 chars. Max 5 global wishlists per user. Fails with `409` if name already exists.

### 2.2 Get All Wishlists
- **Method**: `GET`
- **Endpoint**: `/api/v1/wishlist`
- **Behavior**: Returns wishlist summary objects, sorted by `createdAt DESC`.

### 2.3 Get Wishlist Details (with sorting)
- **Method**: `GET`
- **Endpoint**: `/api/v1/wishlist/:wishlistId`
- **Query Params**: 
  - `sortBy` (Options: `manual`, `addedRecently`, `color`, `yield`, `minInvestment`, `tenure`, `rating`. Default: `manual`)
- **Behavior**: Regardless of the `sortBy` parameter, **Pinned bonds are always returned at the top**.

### 2.4 Rename Wishlist
- **Method**: `PATCH`
- **Endpoint**: `/api/v1/wishlist/:wishlistId`
- **Body**: `{ "name": "string" }`

### 2.5 Delete Wishlist
- **Method**: `DELETE`
- **Endpoint**: `/api/v1/wishlist/:wishlistId`
- **Behavior**: Deletes the wishlist and cascade-deletes all its bond relationships.

---

## 3. Bond Items (Inside Wishlist)

### 3.1 Add Bond
- **Method**: `POST`
- **Endpoint**: `/api/v1/wishlist/:wishlistId/bond`
- **Body**: `{ "bondIsin": "IN123..." }`
- **Rules**: Max 10 bonds per wishlist. Returns `409` if bond is already in the wishlist.

### 3.2 Remove Bond
- **Method**: `DELETE`
- **Endpoint**: `/api/v1/wishlist/:wishlistId/bond/:bondIsin`

### 3.3 Set Bond Color
- **Method**: `PATCH`
- **Endpoint**: `/api/v1/wishlist/:wishlistId/bond/:bondIsin/color`
- **Body**: `{ "color": "#FF0000" }` (or `null` to clear)

### 3.4 Set Bond Position (Manual Sorting)
- **Method**: `PATCH`
- **Endpoint**: `/api/v1/wishlist/:wishlistId/bond/:bondIsin/position`
- **Body**: `{ "position": 2 }` (0-indexed integer)

### 3.5 Pin/Unpin Bond
- **Method**: `PATCH`
- **Endpoint**: `/api/v1/wishlist/:wishlistId/bond/:bondIsin/pin`
- **Body**: `{ "isPinned": true }`

### 3.6 Bulk Reorder Bonds
- **Method**: `PATCH`
- **Endpoint**: `/api/v1/wishlist/:wishlistId/reorder`
- **Body**: `{ "bondIsins": ["ISIN_1", "ISIN_2", "ISIN_3"] }`
- **Rules**: Must contain an array of *all* ISINs currently in the wishlist in the new desired order.

---
**Status Codes Summary:**
- `200/201/204`: Success
- `400`: Bad Request (Validation failure)
- `404`: Not Found (Invalid UUID or resource missing)
- `409`: Conflict (Duplicates)
- `422`: Unprocessable (Limits reached - 5 wishlists, 10 bonds)
