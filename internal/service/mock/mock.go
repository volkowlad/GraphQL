// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	model "TestOzon/internal/handler/graph/model"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// AllowComments mocks base method.
func (m *MockDatabase) AllowComments(ctx context.Context, postID uuid.UUID, allow bool) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllowComments", ctx, postID, allow)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllowComments indicates an expected call of AllowComments.
func (mr *MockDatabaseMockRecorder) AllowComments(ctx, postID, allow interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllowComments", reflect.TypeOf((*MockDatabase)(nil).AllowComments), ctx, postID, allow)
}

// CreateComment mocks base method.
func (m *MockDatabase) CreateComment(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, content string) (*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComment", ctx, postID, parentID, content)
	ret0, _ := ret[0].(*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComment indicates an expected call of CreateComment.
func (mr *MockDatabaseMockRecorder) CreateComment(ctx, postID, parentID, content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComment", reflect.TypeOf((*MockDatabase)(nil).CreateComment), ctx, postID, parentID, content)
}

// CreatePost mocks base method.
func (m *MockDatabase) CreatePost(ctx context.Context, title, context string, allowComments bool) (*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", ctx, title, context, allowComments)
	ret0, _ := ret[0].(*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockDatabaseMockRecorder) CreatePost(ctx, title, context, allowComments interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockDatabase)(nil).CreatePost), ctx, title, context, allowComments)
}

// GetComments mocks base method.
func (m *MockDatabase) GetComments(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComments", ctx, postID, parentID, limit, offset)
	ret0, _ := ret[0].([]*model.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetComments indicates an expected call of GetComments.
func (mr *MockDatabaseMockRecorder) GetComments(ctx, postID, parentID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComments", reflect.TypeOf((*MockDatabase)(nil).GetComments), ctx, postID, parentID, limit, offset)
}

// GetPostByID mocks base method.
func (m *MockDatabase) GetPostByID(ctx context.Context, id uuid.UUID, limit, offset int) (*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostByID", ctx, id, limit, offset)
	ret0, _ := ret[0].(*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostByID indicates an expected call of GetPostByID.
func (mr *MockDatabaseMockRecorder) GetPostByID(ctx, id, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostByID", reflect.TypeOf((*MockDatabase)(nil).GetPostByID), ctx, id, limit, offset)
}

// GetPosts mocks base method.
func (m *MockDatabase) GetPosts(ctx context.Context) ([]*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPosts", ctx)
	ret0, _ := ret[0].([]*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPosts indicates an expected call of GetPosts.
func (mr *MockDatabaseMockRecorder) GetPosts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPosts", reflect.TypeOf((*MockDatabase)(nil).GetPosts), ctx)
}
