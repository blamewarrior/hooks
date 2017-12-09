/*
   Copyright (C) 2016 The BlameWarrior Authors.
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

package hooks

// var SendingError = fmt.Errorf("sending error")

// type Mediator interface {
// 	Mediate(payload string) (err error)
// }

// type MediatorService struct {
// 	ConsumerBaseURL string
// 	c               *http.Client

// 	hooks hooks.Hooks
// }

// func (service *PullRequestPublishService) Mediate(payload string) (err error) {

// 	if position, result, err = service.pullRequests.Save(pullRequest); err != nil {
// 		return err
// 	}

// 	err = service.send(result.ValueBytes)

// 	if err = service.pullRequests.Delete(pullRequest); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (service *PullRequestPublisher) send(payload []byte) (err error) {

// 	response, err := service.c.Post(
// 		service.ConsumerBaseURL+"/process_hook",
// 		"application/json",
// 		bytes.NewBuffer(payload),
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if response.StatusCode != http.StatusOK {
// 		return SendingError
// 	}
// }

// func (service *PullRequestPublishService) assignReviewer(pullRequest *PullRequest) error {
// 	return nil
// }
