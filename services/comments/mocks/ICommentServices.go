// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	dto "discusiin/dto"
	models "discusiin/models"

	mock "github.com/stretchr/testify/mock"
)

// ICommentServices is an autogenerated mock type for the ICommentServices type
type ICommentServices struct {
	mock.Mock
}

// CreateComment provides a mock function with given fields: comment, postID, token
func (_m *ICommentServices) CreateComment(comment models.Comment, postID int, token dto.Token) error {
	ret := _m.Called(comment, postID, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Comment, int, dto.Token) error); ok {
		r0 = rf(comment, postID, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteComment provides a mock function with given fields: commentID, token
func (_m *ICommentServices) DeleteComment(commentID int, token dto.Token) error {
	ret := _m.Called(commentID, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, dto.Token) error); ok {
		r0 = rf(commentID, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllComments provides a mock function with given fields: id
func (_m *ICommentServices) GetAllComments(id int) ([]dto.PublicComment, error) {
	ret := _m.Called(id)

	var r0 []dto.PublicComment
	if rf, ok := ret.Get(0).(func(int) []dto.PublicComment); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.PublicComment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateComment provides a mock function with given fields: newComment, token
func (_m *ICommentServices) UpdateComment(newComment models.Comment, token dto.Token) error {
	ret := _m.Called(newComment, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Comment, dto.Token) error); ok {
		r0 = rf(newComment, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewICommentServices interface {
	mock.TestingT
	Cleanup(func())
}

// NewICommentServices creates a new instance of ICommentServices. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewICommentServices(t mockConstructorTestingTNewICommentServices) *ICommentServices {
	mock := &ICommentServices{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
