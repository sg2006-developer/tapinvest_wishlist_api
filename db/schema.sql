CREATE TABLE master_data (
    bond_name TEXT NOT NULL,
    yield DECIMAL(6,2) NOT NULL,
    payout_frequency TEXT NOT NULL,
    maturity_date DATE NOT NULL,
    min_investment DECIMAL(15,2) NOT NULL,
    rating TEXT NOT NULL,
    isin VARCHAR(20) NOT NULL,
    logo_url TEXT,
    detail_url TEXT,
    tenure DECIMAL(6,2) NOT NULL,
    PRIMARY KEY (isin)
);

CREATE TABLE wish_lists (
    wish_list_id INTEGER GENERATED ALWAYS AS IDENTITY,
    wish_list_name VARCHAR(255) NOT NULL,
    PRIMARY KEY (wish_list_id)
);

CREATE TABLE wish_isin (
    wish_list_id INT NOT NULL,
    isin VARCHAR(20) NOT NULL,

    PRIMARY KEY (wish_list_id, isin),

    CONSTRAINT fk_wishlist
        FOREIGN KEY (wish_list_id)
        REFERENCES wish_lists(wish_list_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_isin
        FOREIGN KEY (isin)
        REFERENCES master_data(isin)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
