-- +goose Up
CREATE TABLE posts (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title TEXT,
  url TEXT UNIQUE,
  description TEXT,
  published_at TEXT,
  feed_id INT REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
