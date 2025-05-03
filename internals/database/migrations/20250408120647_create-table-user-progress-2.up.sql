CREATE TABLE IF NOT EXISTS user_readings (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reading_id INTEGER NOT NULL REFERENCES readings(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    attempt INTEGER DEFAULT 1 NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index
CREATE INDEX IF NOT EXISTS idx_user_readings_user_id_reading_id
ON user_readings (user_id, reading_id);

CREATE INDEX IF NOT EXISTS idx_user_readings_user_id_unit_id
ON user_readings (user_id, unit_id);

ANALYZE user_readings;


CREATE TABLE IF NOT EXISTS user_evaluations (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    evaluation_id INTEGER NOT NULL REFERENCES evaluations(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    attempt INTEGER DEFAULT 1 NOT NULL,
    percentage_grade INTEGER DEFAULT 0 NOT NULL,
    time_duration INTEGER DEFAULT 0 NOT NULL,
    point INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE INDEX IF NOT EXISTS idx_eval_user_eval ON user_evaluations (user_id, evaluation_id);
-- CREATE INDEX IF NOT EXISTS idx_eval_user_unit ON user_evaluations (user_id, unit_id);
-- ANALYZE user_evaluations;


CREATE TABLE IF NOT EXISTS user_exams (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exam_id INTEGER NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    attempt INTEGER NOT NULL DEFAULT 1,
    percentage_grade INTEGER NOT NULL DEFAULT 0,
    time_duration INTEGER NOT NULL DEFAULT 0,
    point INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- CREATE INDEX IF NOT EXISTS idx_exam_user_exam ON user_exams (user_id, exam_id);
-- CREATE INDEX IF NOT EXISTS idx_exam_user_unit ON user_exams (user_id, unit_id);
-- ANALYZE user_exams;