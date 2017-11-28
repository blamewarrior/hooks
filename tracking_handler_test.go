package main_test

import (
	"context"
	"errors"
	"testing"

	main "github.com/blamewarrior/hooks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoriesServiceMock struct {
	mock.Mock
}

func (m *RepositoriesServiceMock) Track(ctx context.Context, repoFullName, callbackURL string) error {

	args := m.Called(ctx, repoFullName, callbackURL)
	return args.Error(0)

}

func (m *RepositoriesServiceMock) Untrack(ctx context.Context, repoFullName, callbackURL string) error {

	args := m.Called(ctx, repoFullName, callbackURL)
	return args.Error(0)

}

func TestTrackingHandler_DoAction(t *testing.T) {
	reposService := new(RepositoriesServiceMock)

	reposService.On(
		"Track",
		context.Background(),
		"blamewarrior/hooks",
		"https://blamewarrior.com/blamewarrior/hooks/webhook",
	).Return(nil)

	reposService.On(
		"Untrack",
		context.Background(),
		"blamewarrior/hooks",
		"https://blamewarrior.com/blamewarrior/hooks/webhook",
	).Return(nil)

	handler := main.NewTrackingHandler("blamewarrior.com")

	suits := []struct {
		Action string
		Err    error
	}{
		{
			"track",
			nil,
		},
		{
			"untrack",
			nil,
		},

		{
			"custom",
			errors.New("Unsupported action custom"),
		},
	}

	for _, suits := range suits {
		err := handler.DoAction(reposService, "blamewarrior/hooks", suits.Action)
		assert.Equal(t, suits.Err, err)
	}

}
