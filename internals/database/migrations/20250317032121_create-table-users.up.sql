CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL CHECK (LENGTH(user_name) >= 3 AND LENGTH(user_name) <= 50),
    email VARCHAR(255) UNIQUE NOT NULL CHECK (POSITION('@' IN email) > 1),
    password VARCHAR(250),
    google_id VARCHAR(255) UNIQUE,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('owner', 'user', 'teacher', 'treasurer', 'admin')),
    security_question TEXT NOT NULL,
    security_answer VARCHAR(255) NOT NULL,
    donation_name VARCHAR(100),
    original_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);