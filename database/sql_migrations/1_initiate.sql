-- +migrate Up
-- +migrate StatementBegin

DROP TABLE IF EXISTS volunteers, donations, emergency_reports, evacuation_routes, distribution_logs, logistics, refugees, shelters, disasters, users;

-- Tabel Users
CREATE TYPE user_role AS ENUM ('admin', 'donor', 'user');

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role user_role NOT NULL,
    contact VARCHAR(20),
    is_2fa BOOLEAN DEFAULT FALSE,
    otp_code VARCHAR(6),
    otp_expiry TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Disasters
CREATE TYPE disaster_status AS ENUM ('active', 'resolved', 'archived');

CREATE TABLE IF NOT EXISTS disasters (
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    description TEXT,
    status disaster_status NOT NULL DEFAULT 'active',
    reported_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Shelters
CREATE TABLE IF NOT EXISTS shelters (
    id SERIAL PRIMARY KEY,
    disaster_id INT REFERENCES disasters(id),
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    capacity_total INT NOT NULL,
    capacity_remaining INT NOT NULL,
    emergency_needs TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Refugees
CREATE TABLE IF NOT EXISTS refugees (
    id SERIAL PRIMARY KEY,
    disaster_id INT REFERENCES disasters(id),
    name VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    condition VARCHAR(100),
    needs TEXT,
    shelter_id INT REFERENCES shelters(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Logistics
CREATE TYPE logistics_status AS ENUM ('available', 'distributed', 'out_of_stock');

CREATE TABLE IF NOT EXISTS logistics (
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    status logistics_status NOT NULL DEFAULT 'available',
    disaster_id INT REFERENCES disasters(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Distribution Logs
CREATE TABLE IF NOT EXISTS distribution_logs (
    id SERIAL PRIMARY KEY,
    logistic_id INT REFERENCES logistics(id),
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    sender_name VARCHAR(255) NOT NULL,
    recipient_name VARCHAR(255) NOT NULL,
    quantity_sent INT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Evacuation Routes
CREATE TYPE evacuation_status AS ENUM ('safe', 'risky', 'blocked');

CREATE TABLE IF NOT EXISTS evacuation_routes (
    id SERIAL PRIMARY KEY,
    disaster_id INT REFERENCES disasters(id),
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    distance DECIMAL(10,2) NOT NULL,
    route TEXT NOT NULL,
    status evacuation_status NOT NULL DEFAULT 'safe',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Emergency Reports
CREATE TABLE IF NOT EXISTS emergency_reports (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    disaster_id INT REFERENCES disasters(id),
    description TEXT NOT NULL,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Donations
CREATE TYPE donation_status AS ENUM ('pending', 'confirmed', 'rejected');

CREATE TABLE IF NOT EXISTS donations (
    id SERIAL PRIMARY KEY,
    donor_id INT REFERENCES users(id),
    disaster_id INT REFERENCES disasters(id),
    amount DECIMAL(10,2),
    item_name VARCHAR(255),
    status donation_status NOT NULL DEFAULT 'pending',
    donated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel Volunteers
CREATE TYPE volunteer_status AS ENUM ('available', 'on_mission', 'completed');

CREATE TABLE IF NOT EXISTS volunteers (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    disaster_id INT REFERENCES disasters(id),
    skill VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    status volunteer_status NOT NULL DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate StatementEnd