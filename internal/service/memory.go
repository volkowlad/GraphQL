package service

import (
	"TestOzon/internal/handler/graph/model"
	"TestOzon/internal/repos/memory"
	"context"
	"github.com/google/uuid"
)

type MemoryService struct {
	repoPost    memory.Post
	repoComment memory.Comment
}

func NewMemoryService(repoPost memory.Post, repoComment memory.Comment) *PostgresService {
	return &PostgresService{repoPost: repoPost, repoComment: repoComment}
}

func (s *MemoryService) CreatePost(ctx context.Context, title, context string, allowComments bool) (model.Post, error) {
	return s.repoPost.CreatePost(ctx, title, context, allowComments)
}

func (s *MemoryService) GetPostByID(ctx context.Context, id uuid.UUID, limit, offset int) (model.Post, error) {
	return s.repoPost.GetPostByID(ctx, id, limit, offset)
}

func (s *MemoryService) GetPosts(ctx context.Context) ([]model.Post, error) {
	return s.repoPost.GetPosts(ctx)
}

func (s *MemoryService) CreateComment(ctx context.Context, postID, parentID uuid.UUID, content string) (model.Comment, error) {
	return s.repoComment.CreateComment(ctx, postID, parentID, content)
}

func (s *MemoryService) GetComments(ctx context.Context, id uuid.UUID, limit, offset int) ([]model.Comment, error) {
	return s.repoComment.GetComments(ctx, id, limit, offset)
}

func (s *MemoryService) AllowComments(ctx context.Context, postID uuid.UUID, allow bool) (bool, error) {
	return s.repoComment.AllowComments(ctx, postID, allow)
}
