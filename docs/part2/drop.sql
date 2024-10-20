-- Disconnect all users from the database

\c bee_management


SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = 'bee_management';

-- Drop the database
DROP DATABASE IF EXISTS bee_management;

-- Confirm deletion
SELECT 'Database bee_management dropped successfully!' AS result;
