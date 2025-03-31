package postgres

import (
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCommentDB_CreateComment(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	type args struct {
		postID   uuid.UUID
		parentID *uuid.UUID
		content  string
	}
	testsTable := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			name: "OK",
			args: args{
				postID:   uuid.UUID{},
				parentID: nil,
				content:  "test content",
			},
			expectedError: nil,
		},
	}

	for _, test := range testsTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mock.ExpectBegin()

			mock.ExpectQuery("SELECT allow_comments FROM %s WHERE id = $1").
				WithArgs(test.args.postID).
				WillReturnRows(pgxmock.NewRows([]string{"allowComments"}).
					AddRow(true))

			mock.ExpectQuery("INSERT INTO comments").
				WithArgs(test.args.postID, test.args.parentID, test.args.content).
				WillReturnRows(pgxmock.NewRows([]string{"id", "createdAt"}).
					AddRow(1, time.Time{}))

			mock.ExpectCommit()

			assert.NoError(t, err)
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestCommentDB_GetComments(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	type args struct {
		postID   uuid.UUID
		parentID *uuid.UUID
		limit    int
		offset   int
	}
	testsTable := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			name: "OK",
			args: args{
				postID:   uuid.UUID{},
				parentID: nil,
				limit:    1,
				offset:   0,
			},
			expectedError: nil,
		},
	}

	for _, test := range testsTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mock.ExpectBegin()

			mock.ExpectQuery("SELECT id, post_id, parent_id, content, created_at FROM comments WHERE parent_id = $1 and post_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4").
				WithArgs(test.args.postID, test.args.parentID, test.args.limit, test.args.offset).
				WillReturnRows(pgxmock.NewRows([]string{"id", "post_id", "parent_id", "content", "created_at"}).
					AddRow(1, 1, 1, "test", time.Time{}))

			mock.ExpectCommit()

			assert.NoError(t, err)

		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestCommentDB_AllowComments(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	type args struct {
		postID        uuid.UUID
		allowComments bool
	}
	testsTable := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			name: "OK",
			args: args{
				postID:        uuid.UUID{},
				allowComments: true,
			},
			expectedError: nil,
		},
	}

	for _, test := range testsTable {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mock.ExpectBegin()

			mock.ExpectQuery("UPDATE %s SET allow_comments = $1 WHERE id = $2").
				WithArgs(test.args.allowComments, test.args.postID).
				WillReturnRows(pgxmock.NewRows([]string{"allowComments"}).
					AddRow(true))

			mock.ExpectCommit()

			assert.NoError(t, err)
		})

		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
