package memory

import (
	"TestOzon/internal/handler/graph/model"
	"context"
	"errors"
	"fmt"
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

	if c.db.Posts[postID].AllowComments == false {
		slog.Error("post comments are not allowed")
		return nil, errors.New("post comments are not allowed")
	}

	if c.db.Posts[postID] == nil {
		slog.Error("post not found")
		return nil, errors.New("post not found")
	}

	if content == "" {
		slog.Error("content is empty")
		return nil, errors.New("content is empty")
	}

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

	if c.db.Posts[postId] == nil {
		slog.Error("post not found")
		return nil, errors.New("post not found")
	}

	comments := c.db.Comments[postId]
	if offset >= len(comments) {
		slog.Error(fmt.Sprintf("offset %d more then length of comments", offset))
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

	if results == nil {
		slog.Error(fmt.Sprintf("comments witn parents_id %v not found", parentId))
		return []*model.Comment{}, errors.New("no comments found")
	}

	return results, nil
}

func (c *CommentsMem) AllowComments(_ context.Context, postID uuid.UUID, allow bool) (bool, error) {
	c.db.mu.Lock()
	defer c.db.mu.Unlock()

	post, exists := c.db.Posts[postID]
	if !exists {
		slog.Error(fmt.Sprintf("post %v not exists", postID))
		return false, errors.New("post not exists")
	}

	post.AllowComments = allow

	return allow, nil
}
