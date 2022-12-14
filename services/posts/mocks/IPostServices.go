// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	dto "discusiin/dto"
	models "discusiin/models"

	mock "github.com/stretchr/testify/mock"
)

// IPostServices is an autogenerated mock type for the IPostServices type
type IPostServices struct {
	mock.Mock
}

// CreatePost provides a mock function with given fields: post, name, token
func (_m *IPostServices) CreatePost(post models.Post, name string, token dto.Token) error {
	ret := _m.Called(post, name, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Post, string, dto.Token) error); ok {
		r0 = rf(post, name, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePost provides a mock function with given fields: id, token
func (_m *IPostServices) DeletePost(id int, token dto.Token) error {
	ret := _m.Called(id, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, dto.Token) error); ok {
		r0 = rf(id, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllPostByLike provides a mock function with given fields: page, search
func (_m *IPostServices) GetAllPostByLike(page int, search string) ([]dto.PublicPost, int, error) {
	ret := _m.Called(page, search)

	var r0 []dto.PublicPost
	if rf, ok := ret.Get(0).(func(int, string) []dto.PublicPost); ok {
		r0 = rf(page, search)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.PublicPost)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(int, string) int); ok {
		r1 = rf(page, search)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int, string) error); ok {
		r2 = rf(page, search)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetPost provides a mock function with given fields: id
func (_m *IPostServices) GetPost(id int) (dto.PublicPost, error) {
	ret := _m.Called(id)

	var r0 dto.PublicPost
	if rf, ok := ret.Get(0).(func(int) dto.PublicPost); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(dto.PublicPost)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPosts provides a mock function with given fields: name, page, search
func (_m *IPostServices) GetPosts(name string, page int, search string) ([]dto.PublicPost, int, error) {
	ret := _m.Called(name, page, search)

	var r0 []dto.PublicPost
	if rf, ok := ret.Get(0).(func(string, int, string) []dto.PublicPost); ok {
		r0 = rf(name, page, search)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.PublicPost)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(string, int, string) int); ok {
		r1 = rf(name, page, search)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, int, string) error); ok {
		r2 = rf(name, page, search)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetPostsByTopicByLike provides a mock function with given fields: name, page
func (_m *IPostServices) GetPostsByTopicByLike(name string, page int) ([]dto.PublicPost, int, error) {
	ret := _m.Called(name, page)

	var r0 []dto.PublicPost
	if rf, ok := ret.Get(0).(func(string, int) []dto.PublicPost); ok {
		r0 = rf(name, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.PublicPost)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(string, int) int); ok {
		r1 = rf(name, page)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, int) error); ok {
		r2 = rf(name, page)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetRecentPost provides a mock function with given fields: page, search
func (_m *IPostServices) GetRecentPost(page int, search string) ([]dto.PublicPost, int, error) {
	ret := _m.Called(page, search)

	var r0 []dto.PublicPost
	if rf, ok := ret.Get(0).(func(int, string) []dto.PublicPost); ok {
		r0 = rf(page, search)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.PublicPost)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(int, string) int); ok {
		r1 = rf(page, search)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(int, string) error); ok {
		r2 = rf(page, search)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SuspendPost provides a mock function with given fields: token, postId
func (_m *IPostServices) SuspendPost(token dto.Token, postId int) error {
	ret := _m.Called(token, postId)

	var r0 error
	if rf, ok := ret.Get(0).(func(dto.Token, int) error); ok {
		r0 = rf(token, postId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePost provides a mock function with given fields: newPost, id, token
func (_m *IPostServices) UpdatePost(newPost models.Post, id int, token dto.Token) error {
	ret := _m.Called(newPost, id, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Post, int, dto.Token) error); ok {
		r0 = rf(newPost, id, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewIPostServices interface {
	mock.TestingT
	Cleanup(func())
}

// NewIPostServices creates a new instance of IPostServices. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIPostServices(t mockConstructorTestingTNewIPostServices) *IPostServices {
	mock := &IPostServices{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
