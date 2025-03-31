package graph

import (
	"TestOzon/internal/handler/graph/model"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

var commentChannels = make(map[uuid.UUID][]chan *model.Comment)

func sendNewComment(postID uuid.UUID, comment *model.Comment) {
	if subs, ok := commentChannels[postID]; ok {
		for _, ch := range subs {
			ch <- comment
		}
	}
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, allowComments bool) (*model.Post, error) {
	slog.Error("CreatePost called")
	return r.services.CreatePost(ctx, title, content, allowComments)
}

// AddComment is the resolver for the addComment field.
func (r *mutationResolver) AddComment(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, content string) (*model.Comment, error) {
	comment, err := r.services.CreateComment(ctx, postID, parentID, content)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to add comment to post with ID %v and parent ID %v. Reason: %s", postID, parentID, err.Error()))
		return nil, err
	}

	slog.Info(fmt.Sprintf("successfully added comment to post with ID %v and parent ID %v", postID, parentID))

	sendNewComment(postID, comment)

	slog.Info(fmt.Sprintf("subscribes got new comment from post with ID %v", postID))

	return comment, nil
}

// AllowComments is the resolver for the allowComments field.
func (r *mutationResolver) AllowComments(ctx context.Context, postID uuid.UUID, allowComments bool) (bool, error) {
	slog.Info(fmt.Sprintf("allowComments is %t on post with ID %v", allowComments, postID))
	return r.services.AllowComments(ctx, postID, allowComments)
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	slog.Info("GetPosts called")
	return r.services.GetPosts(ctx)
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, uuid uuid.UUID, limit int32, offset int32) (*model.Post, error) {
	limitInt := int(limit)
	offsetInt := int(offset)
	slog.Info(fmt.Sprintf("GetPost called, uuid: %v", uuid))
	return r.services.GetPostByID(ctx, uuid, limitInt, offsetInt)
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, limit int32, offset int32) ([]*model.Comment, error) {
	limitInt := int(limit)
	offsetInt := int(offset)
	slog.Info(fmt.Sprintf("GetComment called, postID: %v", postID))
	return r.services.GetComments(ctx, postID, parentID, limitInt, offsetInt)
}

// NewComment is the resolver for the newComment field.
func (r *subscriptionResolver) NewComment(_ context.Context, postID uuid.UUID) (<-chan *model.Comment, error) {
	commentChan := make(chan *model.Comment, 1)

	if _, ok := commentChannels[postID]; !ok {
		commentChannels[postID] = []chan *model.Comment{}
	}

	commentChannels[postID] = append(commentChannels[postID], commentChan)

	slog.Info(fmt.Sprintf("new subscription %v", postID))

	return commentChan, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
