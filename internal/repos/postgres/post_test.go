package postgres

import (
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPostDB_CreatePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	type args struct {
		title         string
		content       string
		allowComments bool
	}

	testTable := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			name:          "OK",
			args:          args{title: "test", content: "test", allowComments: true},
			expectedError: nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mock.ExpectBegin()

			mock.ExpectQuery("INSERT INTO posts").
				WithArgs(test.args.title, test.args.content, test.args.allowComments, time.Now()).
				WillReturnRows(pgxmock.NewRows([]string{"id", "title", "content", "allowComments", "createdAt"}).
					AddRow(1, "test", "test", true, time.Now()))

			mock.ExpectCommit()

			assert.NoError(t, err)
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestPostDB_GetPostByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	type args struct {
		id     uuid.UUID
		limit  int
		offset int
	}
	testTable := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			name:          "OK",
			args:          args{id: uuid.UUID{}, limit: 10, offset: 0},
			expectedError: nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mock.ExpectBegin()

			mock.ExpectQuery("SELECT id, title, content, allow_comments, created_at FROM posts WHERE id = $1").
				WithArgs(test.args.id).
				WillReturnRows(pgxmock.NewRows([]string{"id", "title", "content", "allowComments", "createdAt"}).
					AddRow(1, "test", "test", true, time.Time{}))

			mock.ExpectQuery("SELECT id, post_id, parent_id, content, created_at FROM %s WHERE post_id = $1 and parent_id IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3").
				WithArgs(test.args.id, test.args.limit, test.args.offset).
				WillReturnRows(pgxmock.NewRows([]string{"id", "postID", "parentID", "content", "createdAt"}).
					AddRow(1, 1, 1, "test", time.Time{}))

			mock.ExpectCommit()

			assert.NoError(t, err)
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestPostDB_GetPosts(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testTable := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "OK",
			expectedError: nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mock.ExpectBegin()

			mock.ExpectQuery("SELECT id, title, content, allow_comments, created_at FROM posts").
				WillReturnRows(pgxmock.NewRows([]string{"id", "title", "content", "allowComments", "createdAt"}).
					AddRow(1, "test", "test", true, time.Time{}))

			mock.ExpectCommit()

			assert.NoError(t, err)
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
