-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id            UUID      PRIMARY KEY,
    email         TEXT      UNIQUE NOT NULL,
    password_hash TEXT      NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP NOT NULL DEFAULT now()
);

-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id            UUID      PRIMARY KEY,
    user_id       UUID      NOT NULL,
    title         TEXT      NOT NULL,
    description   TEXT,
    due_date      DATE,
    is_completed  BOOLEAN   NOT NULL DEFAULT false,
    created_at    TIMESTAMP NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP NOT NULL DEFAULT now()
);

-- Create index on user_id for better query performance
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- Function to automatically update the updated_at column
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for users table
CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Trigger for todos table
CREATE TRIGGER set_timestamp_todos
BEFORE UPDATE ON todos
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
