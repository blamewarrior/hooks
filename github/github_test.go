package github_test

import (
	"github.com/stretchr/testify/mock"
)

type tokenServiceMock struct {
	mock.Mock
}

func (tsMock *tokenServiceMock) GetToken() (string, error) {
	args := tsMock.Called()
	return args.String(0), args.Error(1)

}
