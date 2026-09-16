-- +goose Up
CREATE TABLE articles (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author TEXT,
    title TEXT,
    html TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE user_articles (
  user_id BIGINT NOT NULL REFERENCES "users" (id) ON DELETE CASCADE,
  article_id BIGINT NOT NULL REFERENCES "articles" (id) ON DELETE CASCADE,
  UNIQUE(user_id, article_id)
);

-- +goose Down
DROP TABLE user_articles;
DROP TABLE articles;
