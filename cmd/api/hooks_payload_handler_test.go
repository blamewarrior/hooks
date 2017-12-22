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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	main "github.com/blamewarrior/hooks/cmd/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MediatorMock struct {
	mock.Mock
}

func (m *MediatorMock) Mediate(event string, payload []byte) (err error) {
	args := m.Called(event, payload)
	return args.Error(0)
}

func TestHooksPayloadHandler(t *testing.T) {

	payload := make([]byte, 0)

	mediatorMock := new(MediatorMock)
	mediatorMock.On("Mediate", "pull_request", payload).Return(nil)

	handler := main.NewHooksPayloadHandler(mediatorMock)

	req, err := http.NewRequest(
		"POST",
		"/webhook?:username=blamewarrior_user&:repo=public-repo",
		strings.NewReader(string(payload)),
	)

	require.NoError(t, err)

	req.Header.Add("X-GitHub-Event", "pull_request")

	http.DefaultClient.Do(req)

	require.NoError(t, err)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

}
