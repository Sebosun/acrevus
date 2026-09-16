-- name: LinkArticleToUser :one
INSERT INTO user_articles (user_id, article_id)
VALUES ($1, $2)
RETURNING *;

