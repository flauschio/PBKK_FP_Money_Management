-- DROP TABLE IF EXISTS scheduled_transactions CASCADE;
-- DROP TABLE IF EXISTS budgets CASCADE;
-- DROP TABLE IF EXISTS transactions CASCADE;
-- DROP TABLE IF EXISTS categories CASCADE;
-- DROP TABLE IF EXISTS accounts CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;

-- user (id, name, email, password,  refresh_token, access_token)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    email VARCHAR(500) UNIQUE NOT NULL,
    password VARCHAR(500) NOT NULL,
    refresh_token TEXT,
    access_token TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

-- category (id, name, created_at, updated_at)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

-- account (id, bank_name, amount, created_at, updated_at)
CREATE TABLE IF NOT EXISTS scheduled_transaction(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

-- transaction (id, name, category_id, amount, account_id, created_at, updated_at)
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    account_id INTEGER,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- user_id INTERGER REFERENCES user(id) ON DELETE SET NULL,
    -- account_id INTERGER REFERENCES scheduled_transaction(id) ON DELETE SET NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
);

-- budget (id, category_id, amount, criteria (monthly, annual), updated_at)
CREATE TABLE IF NOT EXISTS budget (
    id SERIAL PRIMARY KEY,
    amount DECIMAL(12, 2) NOT NULL,
    criteria ENUM('monthly', 'annually') NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- user_id INTERGER REFERENCES user(id) ON DELETE SET NULL,
    -- account_id INTERGER REFERENCES scheduled_transaction(id) ON DELETE SET NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
);

-- scheduled transaction (id, name, amount, repetition (monthly, anually, 6 months, 3 months), repeat_at, created_at, updated_at)
CREATE TABLE IF NOT EXISTS scheduled_transaction(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    repetition ENUM('monthly', 'anually', '6 months', '3 months') NOT NULL,
    repeat_at DATETIME NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- user_id INTERGER REFERENCES user(id) ON DELETE SET NULL,
    -- account_id INTERGER REFERENCES scheduled_transaction(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

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

INSERT INTO accounts (bank_name, amount) VALUES
('Chase Checking', 5420.50),
('Savings Account', 12800.00),
('Credit Card', -850.25);