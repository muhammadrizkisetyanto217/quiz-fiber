-- +migrate Down
DROP INDEX IF EXISTS idx_user_progress_user_id;
DROP TABLE IF EXISTS user_progress;