-- CREATE TABLE IF NOT EXISTS users (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(500) NOT NULL,
--     email VARCHAR(500) NOT NULL,
--     password VARCHAR(500) NOT NULL,
--     refresh_token,
--     access_token,
-- );

-- category (id, name, created_at, updated_at)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

-- transaction (id, name, category_id, amount, account_id, created_at, updated_at)
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    account_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
)

CREATE INDEX idx_transactions_category_id ON transactions(category_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);

INSERT INTO categories (name) 
('Food & Dining'),
('Transportation'),
('Shopping'),
('Entertainment'),
('Bills & Utilities'),
('Healthcare'),
('Income'),
('Investments');