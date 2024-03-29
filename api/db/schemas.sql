CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(128) UNIQUE NOT NULL,
    age INTEGER NOT NULL,
    email TEXT UNIQUE NOT NULL,
    valid_email BOOLEAN NOT NULL,
    password TEXT NOT NULL,
    roles TEXT[] NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS "tokens" (
    user_id UUID PRIMARY KEY,
    access_token TEXT DEFAULT '' NOT NULL,
    refresh_token TEXT DEFAULT '' NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS "login_attempts" (
    user_id UUID PRIMARY KEY,
    attempts INTEGER DEFAULT 0 NOT NULL,
    last_failed_login_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS "email_validations" (
    user_id UUID PRIMARY KEY,
    verification_code VARCHAR(6) NOT NULL,
    expiration_time INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
)
