CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    source_type_id INT NOT NULL, -- 1=Quiz, 2=Evaluation, 3=Exam
    source_id INT NOT NULL,      -- quizzes_id / evaluation_id / exam_id
    question_text VARCHAR(200) NOT NULL,
    question_answer TEXT[] NOT NULL,
    question_correct VARCHAR(50) NOT NULL,
    tooltips_id INT[],  -- Optional, digunakan jika source_type_id = 1
    status VARCHAR(10) DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'archived')),
    paragraph_help TEXT NOT NULL,
    explain_question TEXT NOT NULL,
    answer_text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tooltips (
    id SERIAL PRIMARY KEY,
    keyword VARCHAR(100) UNIQUE NOT NULL,
    description_short VARCHAR(200) NOT NULL,
    description_long TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS readings (
	id SERIAL PRIMARY KEY,
	title VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'archived')),
	description_long TEXT NOT NULL,
	tooltips_id INT[],
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP,
	unit_id INT,
	created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- Ini akan membuat index GIN pada kolom array tooltips_id di tabel readings
CREATE INDEX idx_readings_tooltips_id ON readings USING GIN (tooltips_id);

