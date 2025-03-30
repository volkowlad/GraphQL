package memory

import (
	"TestOzon/internal/handler/graph/model"
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type PostMem struct {
	db *StoreMemory
}

func NewPostMem(db *StoreMemory) *PostMem {
	return &PostMem{db: db}
}

func (p *PostMem) CreatePost(_ context.Context, title, content string, allowComments bool) (*model.Post, error) {
	p.db.mu.Lock()
	defer p.db.mu.Unlock()

	post := &model.Post{
		ID:            uuid.New(),
		Title:         title,
		Content:       content,
		AllowComments: allowComments,
		CreatedAt:     time.Now(),
	}

	p.db.Posts[post.ID] = post

	return post, nil
}

func (p *PostMem) GetPostByID(_ context.Context, id uuid.UUID, limit, offset int) (*model.Post, error) {
	p.db.mu.Lock()
	defer p.db.mu.Unlock()

	post, exist := p.db.Posts[id]
	if !exist {
		slog.Error("post not found")
		return nil, errors.New("post not found")
	}

	for _, comment := range p.db.Comments[id] {
		if len(comment.ParentID) == 0 {
			post.Comments = append(post.Comments, comment)
		}

	}

	return post, nil
}

func (p *PostMem) GetPosts(_ context.Context) ([]*model.Post, error) {
	p.db.mu.Lock()
	defer p.db.mu.Unlock()

	var posts []*model.Post
	for _, post := range p.db.Posts {
		posts = append(posts, post)
	}

	return posts, nil
}
