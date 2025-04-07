CREATE TABLE IF NOT EXISTS question_saved (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    source_type_id INTEGER NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    question_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat query pencarian
CREATE INDEX IF NOT EXISTS idx_question_saved_user ON question_saved(user_id);
CREATE INDEX IF NOT EXISTS idx_question_saved_question ON question_saved(question_id);


CREATE TABLE IF NOT EXISTS question_mistakes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_type_id INTEGER NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    question_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Index untuk performa pencarian
CREATE INDEX IF NOT EXISTS idx_question_mistakes_user_id ON question_mistakes(user_id);
CREATE INDEX IF NOT EXISTS idx_question_mistakes_question_id ON question_mistakes(question_id);



CREATE TABLE IF NOT EXISTS user_questions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    question_id INT NOT NULL,
    selected_answer TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    source_type_id INT NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    source_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing
CREATE INDEX idx_user_questions_user_id ON user_questions (user_id);
CREATE INDEX idx_user_questions_question_id ON user_questions (question_id);
CREATE INDEX idx_user_questions_source_type_id ON user_questions (source_type_id);
CREATE INDEX idx_user_questions_source_id ON user_questions (source_id);
