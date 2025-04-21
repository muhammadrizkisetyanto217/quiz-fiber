CREATE TABLE IF NOT EXISTS kajian_attendances (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    address TEXT,
    access_time TIMESTAMP NOT NULL,
    topic TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Optional indexes for optimization
CREATE INDEX idx_kajian_attendances_user_id ON kajian_attendances (user_id);
CREATE INDEX idx_kajian_attendances_access_time ON kajian_attendances (access_time);
