-- NOT NULL because go/sql is not good at handling null values.

CREATE TABLE IF NOT EXISTS snippets (
    id SERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(100) NOT NULL DEFAULT 'Untitled',
    author VARCHAR(100) NOT NULL DEFAULT 'Unknown',
    work VARCHAR(100) NOT NULL DEFAULT 'Unknown',
    content TEXT NOT NULL DEFAULT 'Quote',
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '100 days'
);
CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);