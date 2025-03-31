package memory

import (
	"TestOzon/internal/handler/graph/model"
	mock_service "TestOzon/internal/service/mock"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCommentsMem_CreateComment(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockDatabase, postID uuid.UUID, parentID *uuid.UUID, content string)

	testTable := []struct {
		name          string
		postID        uuid.UUID
		parentID      *uuid.UUID
		content       string
		allowComments bool
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:     "OK",
			postID:   uuid.UUID{},
			parentID: &uuid.Nil,
			content:  "test content",
			mockBehaviour: func(s *mock_service.MockDatabase, postID uuid.UUID, parentID *uuid.UUID, content string) {
				s.EXPECT().CreateComment(context.Background(), postID, parentID, content).Return(&model.Comment{
					ID:        uuid.UUID{},
					PostID:    uuid.UUID{},
					ParentID:  &uuid.Nil,
					Content:   "test content",
					CreatedAt: time.Time{},
				}, nil)
			},
		},
		{
			name:          "allowComment false",
			postID:        uuid.New(),
			parentID:      &uuid.Nil,
			content:       "test content",
			allowComments: false,
			mockBehaviour: func(s *mock_service.MockDatabase, postID uuid.UUID, parentID *uuid.UUID, content string) {
				s.EXPECT().CreateComment(context.Background(), postID, parentID, content).Return(&model.Comment{}, errors.New("post comments are not allowed"))
			},
			expectedError: errors.New("post comments are not allowed"),
		},
		{
			name:     "post not found",
			postID:   uuid.New(),
			parentID: &uuid.Nil,
			content:  "test content",
			mockBehaviour: func(s *mock_service.MockDatabase, postID uuid.UUID, parentID *uuid.UUID, content string) {
				s.EXPECT().CreateComment(context.Background(), postID, parentID, content).Return(&model.Comment{}, errors.New("post not found"))
			},
			expectedError: errors.New("post not found"),
		},
		{
			name:     "content is empty",
			postID:   uuid.New(),
			parentID: &uuid.Nil,
			content:  "",
			mockBehaviour: func(s *mock_service.MockDatabase, postID uuid.UUID, parentID *uuid.UUID, content string) {
				s.EXPECT().CreateComment(context.Background(), postID, parentID, content).Return(&model.Comment{}, errors.New("content is empty"))
			},
			expectedError: errors.New("content is empty"),
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			create := mock_service.NewMockDatabase(ctrl)
			test.mockBehaviour(create, test.postID, test.parentID, test.content)

			comment, err := create.CreateComment(context.Background(), test.postID, test.parentID, test.content)
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err)
				return
			}

			assert.Equal(t, test.postID, comment.PostID)
			assert.Equal(t, test.parentID, comment.ParentID)
			assert.Equal(t, test.content, comment.Content)
			assert.NotNil(t, comment.CreatedAt)
			assert.NotNil(t, comment.ID)
		})
	}
}

func TestCommentsMem_GetComments(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockDatabase, postId uuid.UUID, parentId *uuid.UUID, limit, offset int)

	testTable := []struct {
		name          string
		postID        uuid.UUID
		parentID      *uuid.UUID
		limit         int
		offset        int
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:     "OK",
			postID:   uuid.UUID{},
			parentID: &uuid.Nil,
			limit:    1,
			offset:   0,
			mockBehaviour: func(s *mock_service.MockDatabase, postId uuid.UUID, parentId *uuid.UUID, limit, offset int) {
				s.EXPECT().GetComments(context.Background(), postId, parentId, limit, offset).Return([]*model.Comment{
					{
						ID:       uuid.UUID{},
						PostID:   uuid.UUID{},
						ParentID: &uuid.Nil,
						Replies:  []*model.Comment{},
					},
				}, nil)
			},
		},
		{
			name:     "post not found",
			postID:   uuid.New(),
			parentID: &uuid.Nil,
			limit:    1,
			offset:   0,
			mockBehaviour: func(s *mock_service.MockDatabase, postId uuid.UUID, parentId *uuid.UUID, limit, offset int) {
				s.EXPECT().GetComments(context.Background(), postId, parentId, limit, offset).Return([]*model.Comment{}, errors.New("post not found"))
			},
			expectedError: errors.New("post not found"),
		},
		{
			name:     "offset greater than length of commits",
			postID:   uuid.UUID{},
			parentID: &uuid.Nil,
			limit:    1,
			offset:   100,
			mockBehaviour: func(s *mock_service.MockDatabase, postId uuid.UUID, parentId *uuid.UUID, limit, offset int) {
				s.EXPECT().GetComments(context.Background(), postId, parentId, limit, offset).Return([]*model.Comment{}, errors.New("offset out of range"))
			},
			expectedError: errors.New("offset out of range"),
		},
		{
			name:     "server failure",
			postID:   uuid.UUID{},
			parentID: &uuid.Nil,
			limit:    1,
			offset:   0,
			mockBehaviour: func(s *mock_service.MockDatabase, postId uuid.UUID, parentId *uuid.UUID, limit, offset int) {
				s.EXPECT().GetComments(context.Background(), postId, parentId, limit, offset).Return([]*model.Comment{}, errors.New("server failure"))
			},
			expectedError: errors.New("server failure"),
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			create := mock_service.NewMockDatabase(ctrl)
			test.mockBehaviour(create, test.postID, test.parentID, test.limit, test.offset)

			comments, err := create.GetComments(context.Background(), test.postID, test.parentID, test.limit, test.offset)
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err)
				return
			}

			assert.Equal(t, test.limit, len(comments))
			assert.Equal(t, test.offset, len(comments[0].Replies))
			assert.Equal(t, test.postID, comments[0].PostID)
			assert.Equal(t, test.parentID, comments[0].ParentID)
			assert.NotNil(t, comments[0].CreatedAt)
			assert.NotNil(t, comments[0].ID)
		})
	}
}

func TestCommentsMem_AllowCommentsComment(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockDatabase, postID uuid.UUID, allow bool)

	testTable := []struct {
		name          string
		postID        uuid.UUID
		allowComments bool
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:          "OK",
			postID:        uuid.UUID{},
			allowComments: true,
			mockBehaviour: func(s *mock_service.MockDatabase, postID uuid.UUID, allow bool) {
				s.EXPECT().AllowComments(context.Background(), postID, allow).Return(true, nil)
			},
			expectedError: nil,
		},
		{
			name:          "post not exists",
			postID:        uuid.UUID{},
			allowComments: false,
			mockBehaviour: func(s *mock_service.MockDatabase, postID uuid.UUID, allow bool) {
				s.EXPECT().AllowComments(context.Background(), postID, allow).Return(false, errors.New("post not exists"))
			},
			expectedError: errors.New("post not exists"),
		},
		{
			name:          "server failure",
			postID:        uuid.UUID{},
			allowComments: false,
			mockBehaviour: func(s *mock_service.MockDatabase, postID uuid.UUID, allow bool) {
				s.EXPECT().AllowComments(context.Background(), postID, allow).Return(false, errors.New("server failure"))
			},
			expectedError: errors.New("server failure"),
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			create := mock_service.NewMockDatabase(ctrl)
			test.mockBehaviour(create, test.postID, test.allowComments)

			res, err := create.AllowComments(context.Background(), test.postID, test.allowComments)
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err)
				return
			}

			assert.Equal(t, test.allowComments, res)
		})
	}
}
