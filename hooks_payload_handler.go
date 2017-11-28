package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blamewarrior/hooks/blamewarrior"
)

type HooksPayloadHandler struct{}

func (handler *HooksPayloadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := handler.handlePayload(w, req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\t%s\t%v\t%s", "POST", req.RequestURI, http.StatusInternalServerError, err)
	}
}

func (handler *HooksPayloadHandler) handlePayload(w http.ResponseWriter, req *http.Request) error {
	username := req.URL.Query().Get(":username")
	repo := req.URL.Query().Get(":repo")

	fullName := fmt.Sprintf("%s/%s", username, repo)

	respBytes, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))

	if err != nil {
		return err
	}

	if err := req.Body.Close(); err != nil {
		return err
	}

	event := req.Header.Get("X-GitHub-Event")

	switch event {
	case "pull_request":
		fmt.Println(fullName)
		pullRequest := new(blamewarrior.PullRequest)
		err = json.Unmarshal(respBytes, &pullRequest)

		if err != nil {
			return err
		}
	}
	return nil
}
