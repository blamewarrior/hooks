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
package blamewarrior_test

import (
	"testing"

	"github.com/blamewarrior/hooks/blamewarrior"
	"github.com/stretchr/testify/assert"
)

func TestMustNotBeEmpty(t *testing.T) {
	examples := map[string]struct {
		Opts    []interface{}
		Message string
	}{
		"empty": {
			Opts:    []interface{}{},
			Message: "must not be empty",
		},
		"only 1 message": {
			Opts:    []interface{}{"title must not be empty"},
			Message: "title must not be empty",
		},
		"parameterized args": {
			Opts:    []interface{}{"title %s", "must not be empty"},
			Message: "title must not be empty",
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			v := new(blamewarrior.Validator)

			args := example.Opts

			v.MustNotBeEmpty("", args...)

			result := v.ErrorMessages()[0]

			assert.Equal(t, example.Message, result)
		})
	}
}

func TestMustNotBeZero(t *testing.T) {
	examples := map[string]struct {
		Opts    []interface{}
		Message string
	}{
		"empty": {
			Opts:    []interface{}{},
			Message: "must not be zero",
		},
		"only 1 message": {
			Opts:    []interface{}{"id must not be zero"},
			Message: "id must not be zero",
		},
		"parameterized args": {
			Opts:    []interface{}{"id %s", "must not be zero"},
			Message: "id must not be zero",
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			v := new(blamewarrior.Validator)

			args := example.Opts

			v.MustNotBeZero(0, args...)

			result := v.ErrorMessages()[0]

			assert.Equal(t, example.Message, result)
		})
	}
}

func TestMustNotBeZeroLength(t *testing.T) {
	examples := map[string]struct {
		Opts    []interface{}
		Message string
	}{
		"empty": {
			Opts:    []interface{}{},
			Message: "must not be zero length",
		},
		"only 1 message": {
			Opts:    []interface{}{"id must not be zero length"},
			Message: "id must not be zero length",
		},
		"parameterized args": {
			Opts:    []interface{}{"id %s", "must not be zero length"},
			Message: "id must not be zero length",
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			v := new(blamewarrior.Validator)

			args := example.Opts

			var emptyList []interface{}

			v.MustNotBeZeroLength(emptyList, args...)

			result := v.ErrorMessages()[0]

			assert.Equal(t, example.Message, result)
		})
	}
}
