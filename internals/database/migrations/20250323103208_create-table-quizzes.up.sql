-- Buat tabel section_quizzes
CREATE TABLE IF NOT EXISTS section_quizzes (
    id SERIAL PRIMARY KEY,
    name_quizzes VARCHAR(50) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    materials_quizzes TEXT NOT NULL,
    icon_url VARCHAR(100),
    total_quizzes INTEGER[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    unit_id INT REFERENCES units(id) ON DELETE CASCADE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS quizzes (
    id SERIAL PRIMARY KEY,
    name_quizzes VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    total_question INTEGER[] NOT NULL DEFAULT '{}',
    icon_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    section_quizzes_id INT REFERENCES section_quizzes(id) ON DELETE CASCADE,
    created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS exams (
    id SERIAL PRIMARY KEY,
    name_exams VARCHAR(50) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    total_question INTEGER[] NOT NULL DEFAULT '{}',
    icon_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    unit_id INT REFERENCES units(id) ON DELETE CASCADE, 
    created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS evaluations (
    id SERIAL PRIMARY KEY,
    name_evaluation VARCHAR(50) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    total_question INTEGER[] NOT NULL DEFAULT '{}',
    icon_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    unit_id INT REFERENCES units(id) ON DELETE CASCADE, 
    created_by UUID REFERENCES users(id) ON DELETE CASCADE
);