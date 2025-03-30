package postgres

import (
	"TestOzon/internal/handler/graph/model"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type CommentDB struct {
	db *pgxpool.Pool
}

func NewCommentDB(db *pgxpool.Pool) *CommentDB {
	return &CommentDB{db: db}
}

func (c *CommentDB) CreateComment(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, content string) (*model.Comment, error) {
	var comment model.Comment

	tx, err := c.db.Begin(ctx)
	if err != nil {
		slog.Error("error to begin creating post", err.Error())
		return &model.Comment{}, err
	}

	var allow bool
	queryAllow := fmt.Sprintf(`
SELECT allow_comments
FROM %s
WHERE id = $1`, postTable)

	err = tx.QueryRow(ctx, queryAllow, postID).Scan(&allow)
	if err != nil {
		slog.Error("error to query post comments", err.Error())
		return &model.Comment{}, err
	}
	if !allow {
		slog.Error("comments not allowed", postID)
		return &model.Comment{}, errors.New("comments not allowed")
	}

	query := fmt.Sprintf(`
INSERT INTO %s (post_id, parent_id, content)
VALUES ($1, $2, $3)
RETURNING id, created_at`, commentTable)

	if parentID == &uuid.Nil {
		parentID = &uuid.UUID{}
	}

	err = tx.QueryRow(ctx, query, postID, parentID, content).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("error to create post", err.Error())
		return &model.Comment{}, err
	}

	comment.PostID = postID
	comment.ParentID = parentID
	comment.Content = content

	return &comment, tx.Commit(ctx)
}

func (c *CommentDB) GetComments(ctx context.Context, postId uuid.UUID, parentId *uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		slog.Error("error to begin creating post", err.Error())
		return []*model.Comment{}, err
	}

	var comments []*model.Comment

	query := fmt.Sprintf(`
SELECT id, post_id, parent_id, content, created_at
FROM %s
WHERE parent_id = $1 and post_id = $2
ORDER BY created_at DESC
LIMIT $2 OFFSET %3`, commentTable)

	rows, err := tx.Query(ctx, query, parentId, postId, limit, offset)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("error to get comments", err.Error())
		return []*model.Comment{}, err
	}

	for rows.Next() {
		comment := &model.Comment{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt); err != nil {
			tx.Rollback(ctx)
			slog.Error("error to get comments", err.Error())
			return []*model.Comment{}, err
		}

		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		tx.Rollback(ctx)
		slog.Error("error to get comments", err.Error())
		return []*model.Comment{}, err
	}

	return comments, tx.Commit(ctx)
}

func (c *CommentDB) AllowComments(ctx context.Context, postID uuid.UUID, allow bool) (bool, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		slog.Error("error to begin creating post", err.Error())
		return false, err
	}

	query := fmt.Sprintf(`
UPDATE %s
SET allow_comments = $1
WHERE id = $2`, postTable)

	_, err = tx.Exec(ctx, query, allow, postID)
	if err != nil {
		tx.Rollback(ctx)
		slog.Error("error to allow comments", err.Error())
		return false, err
	}

	return allow, tx.Commit(ctx)
}
