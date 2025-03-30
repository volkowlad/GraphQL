package postgres

import (
	"TestOzon/internal/handler/graph/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type PostDB struct {
	db *pgxpool.Pool
}

func NewPostDB(db *pgxpool.Pool) *PostDB {
	return &PostDB{db: db}
}

func (p *PostDB) CreatePost(ctx context.Context, title, context string, allowComments bool) (*model.Post, error) {
	var post model.Post

	tx, err := p.db.Begin(ctx)
	if err != nil {
		slog.Error("error to begin creating post", err.Error())
		return &model.Post{}, err
	}

	query := fmt.Sprintf(`
INSERT INTO %s (title, content, allow_comments, created_at)
VALUES ($1, $2, $3, now())
RETURNING id, title, content, allow_comments, created_at`, postTable)

	err = tx.QueryRow(ctx, query, title, context, allowComments).Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("error to create post", err.Error())
		return &model.Post{}, err
	}

	return &post, tx.Commit(ctx)
}

func (p *PostDB) GetPostByID(ctx context.Context, id uuid.UUID, limit, offset int) (*model.Post, error) {
	post := &model.Post{}
	var comments []*model.Comment

	tx, err := p.db.Begin(ctx)
	if err != nil {
		slog.Error("error to begin creating post", err.Error())
		return &model.Post{}, err
	}

	query := fmt.Sprintf(`
SELECT id, title, content, allow_comments, created_at
FROM %s WHERE id = $1`, postTable)

	err = tx.QueryRow(ctx, query, id).Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("error to get post", err.Error())
		return &model.Post{}, err
	}

	queryComments := fmt.Sprintf(`
SELECT id, post_id, parent_id, content, created_at
FROM %s 
WHERE post_id = $1 and parent_id IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3`, commentTable)
	rows, err := tx.Query(ctx, queryComments, id, limit, offset)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("error to comments post", err.Error())
		return &model.Post{}, err
	}

	for rows.Next() {
		comment := &model.Comment{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt); err != nil {
			tx.Rollback(ctx)
			slog.Error("failed to get one comment", err.Error())
			return &model.Post{}, err
		}

		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		tx.Rollback(ctx)
		slog.Error("failed to get comments", err.Error())
		return &model.Post{}, err
	}

	post.Comments = comments
	return post, tx.Commit(ctx)
}

func (p *PostDB) GetPosts(ctx context.Context) ([]*model.Post, error) {
	var posts []*model.Post

	tx, err := p.db.Begin(ctx)
	if err != nil {
		slog.Error("error to begin get posts", err.Error())
		return []*model.Post{}, err
	}

	query := fmt.Sprintf(`
SELECT id, title, content, allow_comments, created_at
FROM %s`, postTable)
	rows, err := tx.Query(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("error to get posts", err.Error())
		return []*model.Post{}, err
	}

	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt); err != nil {
			tx.Rollback(ctx)
			slog.Error("failed to get one post", err.Error())
			return []*model.Post{}, err
		}

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		tx.Rollback(ctx)
		slog.Error("failed to get posts", err.Error())
		return []*model.Post{}, err
	}

	return posts, tx.Commit(ctx)
}
