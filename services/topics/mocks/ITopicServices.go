// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	dto "discusiin/dto"
	models "discusiin/models"

	mock "github.com/stretchr/testify/mock"
)

// ITopicServices is an autogenerated mock type for the ITopicServices type
type ITopicServices struct {
	mock.Mock
}

// CreateTopic provides a mock function with given fields: topic, token
func (_m *ITopicServices) CreateTopic(topic models.Topic, token dto.Token) (models.Topic, error) {
	ret := _m.Called(topic, token)

	var r0 models.Topic
	if rf, ok := ret.Get(0).(func(models.Topic, dto.Token) models.Topic); ok {
		r0 = rf(topic, token)
	} else {
		r0 = ret.Get(0).(models.Topic)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Topic, dto.Token) error); ok {
		r1 = rf(topic, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTopic provides a mock function with given fields: id
func (_m *ITopicServices) GetTopic(id int) (models.Topic, error) {
	ret := _m.Called(id)

	var r0 models.Topic
	if rf, ok := ret.Get(0).(func(int) models.Topic); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(models.Topic)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTopics provides a mock function with given fields:
func (_m *ITopicServices) GetTopics() ([]models.Topic, error) {
	ret := _m.Called()

	var r0 []models.Topic
	if rf, ok := ret.Get(0).(func() []models.Topic); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Topic)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveTopic provides a mock function with given fields: token, id
func (_m *ITopicServices) RemoveTopic(token dto.Token, id int) error {
	ret := _m.Called(token, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(dto.Token, int) error); ok {
		r0 = rf(token, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveTopic provides a mock function with given fields: topic, token
func (_m *ITopicServices) SaveTopic(topic models.Topic, token dto.Token) error {
	ret := _m.Called(topic, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(models.Topic, dto.Token) error); ok {
		r0 = rf(topic, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTopicDescription provides a mock function with given fields: topic, token
func (_m *ITopicServices) UpdateTopicDescription(topic models.Topic, token dto.Token) (models.Topic, error) {
	ret := _m.Called(topic, token)

	var r0 models.Topic
	if rf, ok := ret.Get(0).(func(models.Topic, dto.Token) models.Topic); ok {
		r0 = rf(topic, token)
	} else {
		r0 = ret.Get(0).(models.Topic)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Topic, dto.Token) error); ok {
		r1 = rf(topic, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewITopicServices interface {
	mock.TestingT
	Cleanup(func())
}

// NewITopicServices creates a new instance of ITopicServices. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewITopicServices(t mockConstructorTestingTNewITopicServices) *ITopicServices {
	mock := &ITopicServices{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
