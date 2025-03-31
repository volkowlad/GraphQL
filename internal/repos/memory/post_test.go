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
)

func TestPostMem_CreatePost(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockDatabase, title, content string, allowComment bool)

	testTable := []struct {
		name          string
		title         string
		content       string
		allowComments bool
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:          "OK",
			title:         "test",
			content:       "test",
			allowComments: true,
			mockBehaviour: func(s *mock_service.MockDatabase, title, content string, allowComments bool) {
				s.EXPECT().CreatePost(context.Background(), title, content, allowComments).Return(&model.Post{
					Title:         title,
					Content:       content,
					AllowComments: allowComments,
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:          "title empty",
			title:         "",
			content:       "test",
			allowComments: true,
			mockBehaviour: func(s *mock_service.MockDatabase, title, content string, allowComments bool) {
				s.EXPECT().CreatePost(context.Background(), title, content, allowComments).Return(&model.Post{}, errors.New("title must not be empty"))
			},
			expectedError: errors.New("title must not be empty"),
		},

		{
			name:          "content empty",
			title:         "test",
			content:       "",
			allowComments: true,
			mockBehaviour: func(s *mock_service.MockDatabase, title, content string, allowComments bool) {
				s.EXPECT().CreatePost(context.Background(), title, content, allowComments).Return(&model.Post{}, errors.New("content must not be empty"))
			},
			expectedError: errors.New("content must not be empty"),
		},
		{
			name:          "service failure",
			title:         "test",
			content:       "test",
			allowComments: true,
			mockBehaviour: func(s *mock_service.MockDatabase, title, content string, allowComments bool) {
				s.EXPECT().CreatePost(context.Background(), title, content, allowComments).Return(nil, errors.New("service error"))
			},
			expectedError: errors.New("service error"),
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_service.NewMockDatabase(c)
			test.mockBehaviour(create, test.title, test.content, test.allowComments)

			post, err := create.CreatePost(context.Background(), test.title, test.content, test.allowComments)

			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err)
				return
			}

			assert.Equal(t, test.title, post.Title)
			assert.Equal(t, test.content, post.Content)
			assert.Equal(t, test.allowComments, post.AllowComments)
			assert.NoError(t, err, "Unexpected error")
			assert.NotNil(t, post.ID, "Post ID should not be nil")
			assert.NotNil(t, post.CreatedAt, "Post CreatedAt should not be nil")
		})
	}
}

func TestPostMem_GetPostsGetPosts(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockDatabase)

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		expected      []*model.Post
		expectedError error
	}{
		{
			name: "OK",
			mockBehaviour: func(s *mock_service.MockDatabase) {
				s.EXPECT().GetPosts(context.Background()).Return([]*model.Post{{}}, nil)
			},
			expected:      []*model.Post{{}},
			expectedError: nil,
		},
		{
			name: "service failure",
			mockBehaviour: func(s *mock_service.MockDatabase) {
				s.EXPECT().GetPosts(context.Background()).Return(nil, errors.New("service error"))
			},
			expectedError: errors.New("service error"),
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_service.NewMockDatabase(c)
			test.mockBehaviour(create)

			posts, err := create.GetPosts(context.Background())
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err)
				return
			}

			assert.Equal(t, test.expected, posts)
		})
	}
}

func TestPostMem_GetPostByIDPostByID(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockDatabase, id uuid.UUID, limit, offset int)

	testTable := []struct {
		name          string
		id            uuid.UUID
		limit         int
		offset        int
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:   "OK",
			id:     uuid.New(),
			limit:  1,
			offset: 0,
			mockBehaviour: func(s *mock_service.MockDatabase, id uuid.UUID, limit, offset int) {
				s.EXPECT().GetPostByID(context.Background(), id, limit, offset).Return(&model.Post{
					ID:       id,
					Comments: []*model.Comment{},
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:   "ID not found",
			id:     uuid.Nil,
			limit:  1,
			offset: 0,
			mockBehaviour: func(s *mock_service.MockDatabase, id uuid.UUID, limit, offset int) {
				s.EXPECT().GetPostByID(context.Background(), id, limit, offset).Return(&model.Post{}, errors.New("post not found"))
			},
			expectedError: errors.New("post not found"),
		},
		{
			name:   "server failure",
			id:     uuid.New(),
			limit:  1,
			offset: 0,
			mockBehaviour: func(s *mock_service.MockDatabase, id uuid.UUID, limit, offset int) {
				s.EXPECT().GetPostByID(context.Background(), id, limit, offset).Return(&model.Post{}, errors.New("server error"))
			},
			expectedError: errors.New("server error"),
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_service.NewMockDatabase(c)
			test.mockBehaviour(create, test.id, test.limit, test.offset)

			post, err := create.GetPostByID(context.Background(), test.id, test.limit, test.offset)
			if test.expectedError != nil {
				assert.Equal(t, test.expectedError, err)
				return
			}

			assert.Equal(t, test.id, post.ID)
			assert.NotNil(t, post.Comments)
		})
	}
}
