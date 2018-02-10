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
package blamewarrior

import (
	"fmt"
)

type Validator struct {
	messages []string
}

func (v *Validator) MustNotBeEmpty(value string, msgArgs ...interface{}) bool {
	if value != "" {
		return true
	}

	v.populateMessages("must not be empty", msgArgs...)

	return false
}

func (v *Validator) MustNotBeZero(value int, msgArgs ...interface{}) bool {
	if value != 0 {
		return true
	}

	v.populateMessages("must not be zero", msgArgs...)

	return false
}

func (v *Validator) MustNotBeZeroLength(value []interface{}, msgArgs ...interface{}) bool {
	if len(value) != 0 {
		return true
	}

	v.populateMessages("must not be zero length", msgArgs...)

	return false
}

func (v *Validator) populateMessages(defaultMsg string, msgArgs ...interface{}) {
	var msg string

	if len(msgArgs) == 0 || msgArgs == nil {
		msg = fmt.Sprintf(defaultMsg)
	}

	if len(msgArgs) == 1 {
		msg = msgArgs[0].(string)
	}
	if len(msgArgs) > 1 {
		msg = fmt.Sprintf(msgArgs[0].(string), msgArgs[1:]...)
	}

	v.messages = append(v.messages, msg)

}

func (v *Validator) ErrorMessages() []string {
	return v.messages
}

func (v *Validator) IsValid() bool {
	return len(v.messages) == 0
}
