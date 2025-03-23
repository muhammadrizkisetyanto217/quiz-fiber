CREATE TABLE IF NOT EXISTS difficulties (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description_short VARCHAR(200),
    description_long VARCHAR(3000),
    total_categories INT,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description_short VARCHAR(100),
    description_long VARCHAR(2000),
    total_subcategories INT,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    difficulty_id INT REFERENCES difficulties(id)
);


CREATE TABLE IF NOT EXISTS subcategories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    description_long VARCHAR(2000),
    total_themes_or_levels INT,
    update_news JSONB,
    image_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    categories_id INT REFERENCES categories(id)
);


CREATE TABLE IF NOT EXISTS themes_or_levels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    description_short VARCHAR(100),
    description_long VARCHAR(2000),
    total_unit INT,
    update_news JSONB,
    image_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    subcategories_id INT REFERENCES subcategories(id)
);


CREATE TABLE IF NOT EXISTS units (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    description_short VARCHAR(200) NOT NULL,
    description_overview TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    themes_or_level_id INT REFERENCES themes_or_levels(id) ON DELETE CASCADE,
    created_by INT REFERENCES users(id) ON DELETE CASCADE
);





