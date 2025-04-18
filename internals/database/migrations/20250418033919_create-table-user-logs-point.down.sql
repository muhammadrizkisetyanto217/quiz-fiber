-- +migrate Down
DROP INDEX IF EXISTS idx_user_point_logs_user_source;
DROP INDEX IF EXISTS idx_user_point_logs_user_id;
DROP TABLE IF EXISTS user_point_logs;