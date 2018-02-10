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

package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blamewarrior/hooks"
)

type HooksPayloadHandler struct {
	mediator hooks.Mediator
}

func NewHooksPayloadHandler(mediator hooks.Mediator) *HooksPayloadHandler {
	return &HooksPayloadHandler{mediator}
}

func (handler *HooksPayloadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := handler.handlePayload(w, req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
	}
}

func (handler *HooksPayloadHandler) handlePayload(w http.ResponseWriter, req *http.Request) error {

	respBytes, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))

	if err != nil {
		return err
	}

	if err := req.Body.Close(); err != nil {
		return err
	}

	event := req.Header.Get("X-GitHub-Event")

	err = handler.mediator.Mediate(event, respBytes)

	return err
}
