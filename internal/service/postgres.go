package service

import (
	"TestOzon/internal/handler/graph/model"
	"TestOzon/internal/repos/postgres"
	"context"
	"github.com/google/uuid"
)

type PostgresService struct {
	repoPost    postgres.Post
	repoComment postgres.Comment
}

func NewPostgresService(repoPost postgres.Post, repoComment postgres.Comment) *PostgresService {
	return &PostgresService{repoPost: repoPost, repoComment: repoComment}
}

func (s *PostgresService) CreatePost(ctx context.Context, title, context string, allowComments bool) (*model.Post, error) {
	return s.repoPost.CreatePost(ctx, title, context, allowComments)
}

func (s *PostgresService) GetPostByID(ctx context.Context, id uuid.UUID, limit, offset int) (*model.Post, error) {
	return s.repoPost.GetPostByID(ctx, id, limit, offset)
}

func (s *PostgresService) GetPosts(ctx context.Context) ([]*model.Post, error) {
	return s.repoPost.GetPosts(ctx)
}

func (s *PostgresService) CreateComment(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, content string) (*model.Comment, error) {
	return s.repoComment.CreateComment(ctx, postID, parentID, content)
}

func (s *PostgresService) GetComments(ctx context.Context, postId uuid.UUID, parentId *uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	return s.repoComment.GetComments(ctx, postId, parentId, limit, offset)
}

func (s *PostgresService) AllowComments(ctx context.Context, postID uuid.UUID, allow bool) (bool, error) {
	return s.repoComment.AllowComments(ctx, postID, allow)
}
