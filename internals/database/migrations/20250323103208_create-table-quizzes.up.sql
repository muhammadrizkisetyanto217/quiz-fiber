-- Buat tabel section_quizzes
CREATE TABLE IF NOT EXISTS section_quizzes (
    id SERIAL PRIMARY KEY,
    name_quizzes VARCHAR(50) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    materials_quizzes TEXT NOT NULL,
    icon_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by INT REFERENCES users(id) ON DELETE CASCADE,
    unit_id INT REFERENCES units(id) ON DELETE CASCADE
);