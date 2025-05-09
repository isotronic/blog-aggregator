-- +goose Up
CREATE TABLE posts(
  id UUID PRIMARY KEY,
  title TEXT,
  url TEXT UNIQUE NOT NULL,
  description TEXT,
  published_at TIMESTAMP NOT NULL,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE posts;