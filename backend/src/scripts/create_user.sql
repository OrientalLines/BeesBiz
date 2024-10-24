-- Create a new user named 'user' with password 'postgres'
CREATE USER "user" WITH PASSWORD 'postgres';

-- Grant all privileges on all tables in the public schema to the new user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "user";

-- Allow the new user to create databases
ALTER USER "user" CREATEDB;

-- Connect to the 'postgres' database (or your specific database if different)
\c postgres

-- Grant all privileges on the current database to the new user
GRANT ALL PRIVILEGES ON DATABASE postgres TO "user";

-- If you want this user to have superuser privileges (be cautious with this)
-- ALTER USER "user" WITH SUPERUSER;

-- Confirm the user was created
\du

-- Plain script

-- Create a new user named 'user' with password 'postgres'
CREATE USER "user" WITH PASSWORD 'postgres';

-- Grant all privileges on all tables in the public schema to the new user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "user";

-- Allow the new user to create databases
ALTER USER "user" CREATEDB;

-- Grant all privileges on the current database to the new user
GRANT ALL PRIVILEGES ON DATABASE postgres TO "user";
