-- Disable foreign keys temporarily to avoid conflicts
SET session_replication_role = replica;

TRUNCATE TABLE users RESTART IDENTITY CASCADE;

SET session_replication_role = DEFAULT;
