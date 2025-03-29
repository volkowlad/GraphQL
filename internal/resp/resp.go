package resp

import (
	"time"
)

type PostResp struct {
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	AllowComments bool           `json:"allowComments"`
	CreatedAt     time.Time      `json:"createdAt"`
	Comments      []*CommentResp `json:"comments"`
}

type CommentResp struct {
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"createdAt"`
	Replies   []*CommentResp `json:"replies"`
}
