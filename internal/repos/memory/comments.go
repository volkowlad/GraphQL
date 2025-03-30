package memory

import (
	"TestOzon/internal/handler/graph/model"
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type CommentsMem struct {
	db *StoreMemory
}

func NewCommentsMem(db *StoreMemory) *CommentsMem {
	return &CommentsMem{db: db}
}

func (c *CommentsMem) CreateComment(_ context.Context, postID uuid.UUID, parentID *uuid.UUID, content string) (*model.Comment, error) {
	c.db.mu.Lock()
	defer c.db.mu.Unlock()

	comment := &model.Comment{
		ID:        uuid.New(),
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	c.db.Comments[comment.PostID] = append(c.db.Comments[comment.PostID], comment)

	return comment, nil
}

func (c *CommentsMem) GetComments(_ context.Context, postId uuid.UUID, parentId *uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	c.db.mu.Lock()
	defer c.db.mu.Unlock()

	//_ = ctx

	comments := c.db.Comments[postId]
	if offset >= len(comments) {
		slog.Error("offset more then length of comments", offset)
		return []*model.Comment{}, errors.New("offset out of range")
	}

	end := offset + limit
	if end > len(comments) {
		end = len(comments)
	}

	var results []*model.Comment
	for _, comment := range comments[offset:end] {
		if comment.ParentID == parentId {
			results = append(results, comment)
		}
	}

	return results, nil
}

func (c *CommentsMem) AllowComments(_ context.Context, postID uuid.UUID, allow bool) (bool, error) {
	c.db.mu.Lock()
	defer c.db.mu.Unlock()

	//_ = ctx

	post, exists := c.db.Posts[postID]
	if !exists {
		slog.Error("post not exists", postID)
		return false, errors.New("post not exists")
	}

	post.AllowComments = allow

	return allow, nil
}
