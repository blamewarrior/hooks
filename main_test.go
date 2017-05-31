package main_test

import (
	"github.com/stretchr/testify/mock"
)

type GitServiceMock struct {
	mock.Mock
}

func (m *GitServiceMock) Track(number int) (bool, error) {

	args := m.Called(number)
	return args.Bool(0), args.Error(1)

}
