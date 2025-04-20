CREATE TABLE IF NOT EXISTS level_requirements (
    id SERIAL PRIMARY KEY,
    level INT UNIQUE NOT NULL,
    name VARCHAR(100),
    min_points INT NOT NULL,
    max_points INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ANALYZE level_requirements;