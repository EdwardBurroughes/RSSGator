-- +goose Up
CREATE TABLE feed_follows (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_id INT NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  CONSTRAINT unique_user_feed UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;