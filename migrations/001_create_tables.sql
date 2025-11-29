DROP TABLE IF EXISTS scheduled_transactions CASCADE;
DROP TABLE IF EXISTS budgets CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- user (id, name, email, password, refresh_token, access_token, created_at, updated_at)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    email VARCHAR(500) UNIQUE NOT NULL,
    password VARCHAR(500) NOT NULL,
    refresh_token TEXT,
    access_token TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- category (id, name, created_at, updated_at)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- account (id, bank_name, amount, created_at, updated_at)
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    bank_name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- transaction (id, name, category_id, amount, account_id, created_at, updated_at)
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- budget (id, category_id, amount, criteria (monthly, annual), created_at, updated_at)
CREATE TABLE IF NOT EXISTS budgets (
    id SERIAL PRIMARY KEY,
    amount DECIMAL(12, 2) NOT NULL,
    criteria VARCHAR(50) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL
);

-- scheduled transaction (id, name, amount, repetition (monthly, annually, 6 months, 3 months), repeat_at, category_id, account_id, created_at, updated_at)
CREATE TABLE IF NOT EXISTS scheduled_transactions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    repetition VARCHAR(50) NOT NULL,
    repeat_at TIMESTAMP NOT NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_category_id ON transactions(category_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);

-- INSERT INTO categories (name) VALUES
-- ('Food & Dining'),
-- ('Transportation'),
-- ('Shopping'),
-- ('Entertainment'),
-- ('Bills & Utilities'),
-- ('Healthcare'),
-- ('Income'),
-- ('Investments');

-- INSERT INTO accounts (bank_name, amount) VALUES
-- ('Chase Checking', 5420.50),
-- ('Savings Account', 12800.00),
-- ('Credit Card', -850.25);
