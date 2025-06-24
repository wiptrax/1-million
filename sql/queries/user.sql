-- name: CreateUser :one
INSERT INTO users (insert_time_milli)
VALUES ($1)
RETURNING *;