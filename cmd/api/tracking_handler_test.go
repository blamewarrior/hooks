/*
   Copyright (C) 2017 The BlameWarrior Authors.
   This file is a part of BlameWarrior service.
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main_test

import (
	"errors"
	"testing"

	main "github.com/blamewarrior/hooks/cmd/api"
	"github.com/blamewarrior/hooks/github"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoriesServiceMock struct {
	mock.Mock
}

func (m *RepositoriesServiceMock) Track(ctx github.Context, repoFullName, callbackURL string) error {

	args := m.Called(ctx, repoFullName, callbackURL)
	return args.Error(0)

}

func (m *RepositoriesServiceMock) Untrack(ctx github.Context, repoFullName, callbackURL string) error {

	args := m.Called(ctx, repoFullName, callbackURL)
	return args.Error(0)

}

func TestTrackingHandler_DoAction(t *testing.T) {
	reposService := new(RepositoriesServiceMock)

	reposService.On(
		"Track",
		github.Context{},
		"blamewarrior/hooks",
		"https://blamewarrior.com/blamewarrior/hooks/webhook",
	).Return(nil)

	reposService.On(
		"Untrack",
		github.Context{},
		"blamewarrior/hooks",
		"https://blamewarrior.com/blamewarrior/hooks/webhook",
	).Return(nil)

	handler := main.NewTrackingHandler("blamewarrior.com", reposService)

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
		err := handler.DoAction("blamewarrior/hooks", suits.Action)
		assert.Equal(t, suits.Err, err)
	}

}
