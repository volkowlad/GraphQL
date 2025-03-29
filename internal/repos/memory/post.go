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

func (p *PostMem) CreatePost(ctx context.Context, title, content string, allowComments bool) (*model.Post, error) {
	p.db.mu.Lock()
	defer p.db.mu.Unlock()

	_ = ctx

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

func (p *PostMem) GetPostByID(ctx context.Context, id uuid.UUID, limit, offset int) (*model.Post, error) {
	p.db.mu.Lock()
	defer p.db.mu.Unlock()

	_ = ctx

	post, exist := p.db.Posts[id]
	if !exist {
		slog.Error("post not found")
		return nil, errors.New("post not found")
	}

	return post, nil
}

func (p *PostMem) GetPosts(ctx context.Context) ([]*model.Post, error) {
	p.db.mu.Lock()
	defer p.db.mu.Unlock()

	_ = ctx

	var posts []*model.Post
	for _, post := range p.db.Posts {
		posts = append(posts, post)
	}

	return posts, nil
}
