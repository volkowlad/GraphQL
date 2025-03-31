package graph

import (
	"TestOzon/internal/handler/graph/model"
	"github.com/google/uuid"
)

func fetchCommentsByPostID(comments []*model.Comment) ([][]*model.Comment, error) {
	commentsMap := make(map[*uuid.UUID][]*model.Comment)
	for _, comment := range comments {
		commentsMap[comment.ParentID] = append(commentsMap[comment.ParentID], comment)
	}

	result := make([][]*model.Comment, 0, len(comments))
	for i, key := range comments {
		if comment, ok := commentsMap[key.ParentID]; ok {
			result[i] = comment
		} else {
			result[i] = []*model.Comment{}
		}
	}

	return result, nil
}
