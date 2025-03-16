// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package libsqlssg

import (
	"context"
	"database/sql"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (username, name, password)
VALUES (?, ?, ?) RETURNING id
`

type CreateAuthorParams struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createAuthor, arg.Username, arg.Name, arg.Password)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createPost = `-- name: CreatePost :one
INSERT INTO posts (title, slug, body, metadata, author_id) 
VALUES (?, ?, ?, ?, ?) RETURNING id, title, slug, body, metadata, created_at, updated_at, author_id
`

type CreatePostParams struct {
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Body     string `json:"body"`
	Metadata string `json:"metadata"`
	AuthorID int64  `json:"author_id"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Title,
		arg.Slug,
		arg.Body,
		arg.Metadata,
		arg.AuthorID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Slug,
		&i.Body,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AuthorID,
	)
	return i, err
}

const getAuthorByID = `-- name: GetAuthorByID :one
SELECT id, username, name, password, is_admin FROM authors WHERE id = ?
`

func (q *Queries) GetAuthorByID(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthorByID, id)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Name,
		&i.Password,
		&i.IsAdmin,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, name, password, is_admin FROM authors WHERE username = ?
`

type GetUserRow struct {
	ID       int64        `json:"id"`
	Name     string       `json:"name"`
	Password string       `json:"password"`
	IsAdmin  sql.NullBool `json:"is_admin"`
}

func (q *Queries) GetUser(ctx context.Context, username string) (GetUserRow, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Password,
		&i.IsAdmin,
	)
	return i, err
}

const updateAuthor = `-- name: UpdateAuthor :one
UPDATE authors SET name = ?, password = ?, is_admin = ? WHERE id = ? RETURNING id, username, name, password, is_admin
`

type UpdateAuthorParams struct {
	Name     string       `json:"name"`
	Password string       `json:"password"`
	IsAdmin  sql.NullBool `json:"is_admin"`
	ID       int64        `json:"id"`
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) (Author, error) {
	row := q.db.QueryRowContext(ctx, updateAuthor,
		arg.Name,
		arg.Password,
		arg.IsAdmin,
		arg.ID,
	)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Name,
		&i.Password,
		&i.IsAdmin,
	)
	return i, err
}
