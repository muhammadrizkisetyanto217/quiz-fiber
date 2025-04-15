-- Migration untuk membuat tabel progress pengguna

CREATE TABLE IF NOT EXISTS user_quizzes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    quiz_id INTEGER NOT NULL,
    attempt INTEGER NOT NULL DEFAULT 1,
    percentage_grade INTEGER NOT NULL DEFAULT 0,
    time_duration INTEGER NOT NULL DEFAULT 0,
    point INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_quizzes_user_id ON user_quizzes (user_id);
CREATE INDEX IF NOT EXISTS idx_user_quizzes_quiz_id ON user_quizzes (quiz_id);

CREATE TABLE IF NOT EXISTS user_section_quizzes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    section_quizzes_id INTEGER NOT NULL,
    complete_quiz INTEGER[] NOT NULL DEFAULT '{}',
    total_quiz INTEGER[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_section_quizzes_user_id ON user_section_quizzes (user_id);
CREATE INDEX IF NOT EXISTS idx_user_section_quizzes_section_id ON user_section_quizzes (section_quizzes_id);

CREATE TABLE IF NOT EXISTS user_unit (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    unit_id INTEGER NOT NULL,
    attempt_reading INTEGER DEFAULT 0 NOT NULL,
    attempt_evaluation INTEGER DEFAULT 0 NOT NULL,
    complete_section_quizzes INTEGER[] NOT NULL DEFAULT '{}',
    total_section_quizzes INTEGER[] NOT NULL DEFAULT '{}',
    grade_exam INTEGER NOT NULL DEFAULT 0,
    grade_result INTEGER NOT NULL DEFAULT 0,
    is_passed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat pencarian berdasarkan user_id dan unit_id
CREATE INDEX IF NOT EXISTS idx_user_unit_user_id_unit_id ON user_unit (user_id, unit_id);

CREATE INDEX IF NOT EXISTS idx_user_unit_user_id ON user_unit (user_id);
CREATE INDEX IF NOT EXISTS idx_user_unit_unit_id ON user_unit (unit_id);
ANALYZE user_unit;

CREATE TABLE IF NOT EXISTS user_themes_or_levels (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    themes_or_levels_id INTEGER NOT NULL,
    complete_unit JSONB NOT NULL DEFAULT '{}'::jsonb,
    total_unit INTEGER[] NOT NULL DEFAULT '{}',
    grade_result INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_themes_user_id ON user_themes_or_levels (user_id);
CREATE INDEX IF NOT EXISTS idx_user_themes_theme_id ON user_themes_or_levels (themes_or_levels_id);

CREATE TABLE IF NOT EXISTS user_subcategory (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    subcategory_id INTEGER NOT NULL,
    complete_themes_or_levels JSONB NOT NULL DEFAULT '{}'::jsonb,
    total_themes_or_levels INTEGER[] NOT NULL DEFAULT '{}',
    grade_result INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_subcategory_user_id ON user_subcategory (user_id);
CREATE INDEX IF NOT EXISTS idx_user_subcategory_subcat_id ON user_subcategory (subcategory_id);

CREATE TABLE IF NOT EXISTS user_category (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    category_id INTEGER NOT NULL,
    complete_category INTEGER[] NOT NULL DEFAULT '{}',
    total_category INTEGER[] NOT NULL DEFAULT '{}',
    grade_result INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_category_user_id ON user_category (user_id);
CREATE INDEX IF NOT EXISTS idx_user_category_cat_id ON user_category (category_id);