CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,

    email      VARCHAR(255) UNIQUE NOT NULL,

    password   VARCHAR(255) NOT NULL,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS listings (
    id          SERIAL PRIMARY KEY,

    user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    title       VARCHAR(255) NOT NULL,
    description TEXT,

    price       BIGINT NOT NULL,

    area        FLOAT NOT NULL,

    rooms       INT DEFAULT 0,
    floor       INT DEFAULT 0,
    total_floors INT DEFAULT 0,
    city        VARCHAR(100) NOT NULL,
    address     VARCHAR(255),
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,

    deal_type VARCHAR(20) NOT NULL CHECK (deal_type IN ('sale','rent')),

    property_type VARCHAR(20) NOT NULL
    CHECK (property_type IN ('apartment','house','commercial','land')),

    status VARCHAR(20) NOT NULL DEFAULT 'active'
    CHECK (status IN ('active','inactive','sold','rented')),
    views_count INT DEFAULT 0,

    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_listings_city    ON listings(city);
CREATE INDEX IF NOT EXISTS idx_listings_type    ON listings(type);
CREATE INDEX IF NOT EXISTS idx_listings_price   ON listings(price);
CREATE INDEX IF NOT EXISTS idx_listings_user_id ON listings(user_id);

CREATE TABLE IF NOT EXISTS favorites (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    listing_id INT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,

    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(user_id, listing_id))
