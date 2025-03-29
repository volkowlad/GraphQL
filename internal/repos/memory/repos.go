package memory

import (
	"TestOzon/internal/handler/graph/model"
	"context"
	"github.com/google/uuid"
)

type Post interface {
	CreatePost(ctx context.Context, title, context string, allowComments bool) (*model.Post, error)
	GetPostByID(ctx context.Context, id uuid.UUID, limit, offset int) (*model.Post, error)
	GetPosts(ctx context.Context) ([]*model.Post, error)
}

type Comment interface {
	CreateComment(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, content string) (*model.Comment, error)
	GetComments(ctx context.Context, postId uuid.UUID, parentID *uuid.UUID, limit, offset int) ([]*model.Comment, error)
	AllowComments(ctx context.Context, postID uuid.UUID, allow bool) (bool, error)
}

type Repos struct {
	Post
	Comment
}

func NewReposMemory(db *StoreMemory) *Repos {
	return &Repos{
		Post:    NewPostMem(db),
		Comment: NewCommentsMem(db),
	}
}
