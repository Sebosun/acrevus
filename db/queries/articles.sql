-- name: CreateArticle :one
INSERT INTO articles (author, title, html, url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetArticle :one
SELECT * FROM articles
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetArticleByURL :one
SELECT * FROM articles
WHERE url = $1 AND deleted_at IS NULL;

-- name: GetUserArticles :many
SELECT articles.*
FROM articles
JOIN user_articles ON user_articles.article_id = articles.id
WHERE user_articles.user_id = $1
    AND articles.deleted_at IS NULL;

