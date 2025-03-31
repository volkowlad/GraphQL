package service

import (
	"TestOzon/internal/handler/graph/model"
	"TestOzon/internal/repos/memory"
	"TestOzon/internal/repos/postgres"
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mock/mock.go

type Database interface {
	CreatePost(ctx context.Context, title, context string, allowComments bool) (*model.Post, error)
	GetPostByID(ctx context.Context, id uuid.UUID, limit, offset int) (*model.Post, error)
	GetPosts(ctx context.Context) ([]*model.Post, error)
	CreateComment(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, content string) (*model.Comment, error)
	GetComments(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, limit, offset int) ([]*model.Comment, error)
	AllowComments(ctx context.Context, postID uuid.UUID, allow bool) (bool, error)
}

type Service struct {
	Database
}

func NewServiceP(repo *postgres.Repos) *Service {
	return &Service{
		Database: NewPostgresService(repo.Post, repo.Comment),
	}
}

func NewServiceM(repo *memory.Repos) *Service {
	return &Service{
		Database: NewMemoryService(repo.Post, repo.Comment),
	}
}
