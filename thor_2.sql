-- ==========================================
-- Wish Database Schema
-- ==========================================

CREATE DATABASE wish;

CREATE USER wish_user WITH PASSWORD 'wish_password';

GRANT ALL PRIVILEGES ON DATABASE wish TO wish_user;


-- ==========================================
-- Master Bond Data
-- ==========================================

CREATE TABLE IF NOT EXISTS master_data (
    bond_name TEXT NOT NULL,
    yield DECIMAL(6,2),
    payout_frequency TEXT,
    maturity_date DATE,
    min_investment BIGINT,
    rating VARCHAR(20),
    isin VARCHAR(20) PRIMARY KEY,
    logo_url TEXT,
    detail_url TEXT,
    tenure DECIMAL(6,2) NOT NULL
);

-- Search by Bond Name
CREATE INDEX IF NOT EXISTS idx_master_bond_name
ON master_data(bond_name);


-- ==========================================
-- User Wish Lists
-- ==========================================

CREATE TABLE IF NOT EXISTS wish_lists (
    wish_list_id UUID PRIMARY KEY,
    wish_list_name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- ==========================================
-- Wishlist ↔ Bond Mapping
-- ==========================================

CREATE TABLE IF NOT EXISTS wish_isin (
    wish_list_id UUID NOT NULL,
    isin VARCHAR(20) NOT NULL,

    color VARCHAR(20),
    position INT NOT NULL DEFAULT 0,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (wish_list_id, isin),

    CONSTRAINT fk_wishlist
        FOREIGN KEY (wish_list_id)
        REFERENCES wish_lists(wish_list_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_master_data
        FOREIGN KEY (isin)
        REFERENCES master_data(isin)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- ======================================================
-- SAMPLE CRUD QUERIES (Reference Only)
-- These queries are provided as examples for backend/API
-- implementation and are intentionally commented out.
-- ======================================================


-- ======================================================
-- Insert a Bond
-- ======================================================

-- INSERT INTO master_data (
--     bond_name,
--     yield,
--     payout_frequency,
--     maturity_date,
--     min_investment,
--     rating,
--     isin,
--     logo_url,
--     detail_url,
--     tenure
-- )
-- VALUES (
--     'HDFC Bond',
--     8.25,
--     'Monthly',
--     '2030-12-31',
--     10000,
--     'AAA',
--     'INE123A01010',
--     'https://example.com/logo.png',
--     'https://example.com/bonds/INE123A01010',
--     4.52
-- );


-- ======================================================
-- Create a Wishlist
-- ======================================================

-- INSERT INTO wish_lists (
--     wish_list_id,
--     wish_list_name
-- )
-- VALUES (
--     '550e8400-e29b-41d4-a716-446655440000',
--     'Long Term Investments'
-- );


-- ======================================================
-- Add Bond to Wishlist
-- ======================================================

-- INSERT INTO wish_isin (
--     wish_list_id,
--     isin,
--     color,
--     position,
--     is_pinned
-- )
-- VALUES (
--     '550e8400-e29b-41d4-a716-446655440000',
--     'INE123A01010',
--     'green',
--     1,
--     FALSE
-- );


-- ======================================================
-- Remove Bond from Wishlist
-- ======================================================

-- DELETE FROM wish_isin
-- WHERE wish_list_id = '550e8400-e29b-41d4-a716-446655440000'
--   AND isin = 'INE123A01010';


-- ======================================================
-- Rename Wishlist
-- ======================================================

-- UPDATE wish_lists
-- SET wish_list_name = 'High Yield Bonds'
-- WHERE wish_list_id = '550e8400-e29b-41d4-a716-446655440000';


-- ======================================================
-- Duplicate Wishlist
-- ======================================================

-- INSERT INTO wish_lists (
--     wish_list_id,
--     wish_list_name
-- )
-- VALUES (
--     '550e8400-e29b-41d4-a716-446655440001',
--     'Long Term Investments Copy'
-- );

-- INSERT INTO wish_isin (
--     wish_list_id,
--     isin,
--     color,
--     position,
--     is_pinned
-- )
-- SELECT
--     '550e8400-e29b-41d4-a716-446655440001',
--     isin,
--     color,
--     position,
--     is_pinned
-- FROM wish_isin
-- WHERE wish_list_id = '550e8400-e29b-41d4-a716-446655440000';


-- ======================================================
-- Move Bond Between Wishlists
-- ======================================================

-- DELETE FROM wish_isin
-- WHERE wish_list_id = '550e8400-e29b-41d4-a716-446655440000'
--   AND isin = 'INE123A01010';

-- INSERT INTO wish_isin (
--     wish_list_id,
--     isin,
--     color,
--     position,
--     is_pinned
-- )
-- VALUES (
--     '550e8400-e29b-41d4-a716-446655440001',
--     'INE123A01010',
--     'green',
--     1,
--     FALSE
-- );


-- ======================================================
-- Search Bond by ISIN
-- ======================================================

-- SELECT *
-- FROM master_data
-- WHERE isin = 'INE123A01010';


-- ======================================================
-- Search Bond by Name
-- ======================================================

-- SELECT *
-- FROM master_data
-- WHERE bond_name ILIKE 'HDFC%';


-- ======================================================
-- Get All Bonds in a Wishlist
-- ======================================================

-- SELECT
--     m.*,
--     wi.color,
--     wi.position,
--     wi.is_pinned
-- FROM wish_isin wi
-- JOIN master_data m
--     ON wi.isin = m.isin
-- WHERE wi.wish_list_id = '550e8400-e29b-41d4-a716-446655440000'
-- ORDER BY wi.position;


-- ======================================================
-- Sort Bonds by Yield
-- ======================================================

-- SELECT *
-- FROM master_data
-- ORDER BY yield DESC;


-- ======================================================
-- Sort Bonds by Rating
-- ======================================================

-- SELECT *
-- FROM master_data
-- ORDER BY rating;


-- ======================================================
-- Sort Bonds by Tenure
-- ======================================================

-- SELECT *
-- FROM master_data
-- ORDER BY tenure DESC;


-- ======================================================
-- Sort Bonds by Minimum Investment
-- ======================================================

-- SELECT *
-- FROM master_data
-- ORDER BY min_investment ASC;


-- ======================================================
-- Sort Bonds by Payout Frequency
-- ======================================================

-- SELECT *
-- FROM master_data
-- ORDER BY payout_frequency;