-- Drop indexes terlebih dahulu (opsional, karena drop table akan otomatis hapus index juga)
DROP INDEX IF EXISTS idx_kajian_attendances_user_id;
DROP INDEX IF EXISTS idx_kajian_attendances_access_time;

-- Drop table
DROP TABLE IF EXISTS kajian_attendances;
