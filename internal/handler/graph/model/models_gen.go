// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID  `json:"id"`
	PostID    uuid.UUID  `json:"postId"`
	ParentID  *uuid.UUID `json:"parentId,omitempty"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	Replies   []*Comment `json:"replies"`
}

type Mutation struct {
}

type Post struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	AllowComments bool       `json:"allowComments"`
	CreatedAt     time.Time  `json:"createdAt"`
	Comments      []*Comment `json:"comments"`
}

type Query struct {
}

type Subscription struct {
}
