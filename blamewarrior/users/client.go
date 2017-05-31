package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Token string `json:"token"`
}

type Client struct {
	BaseURL string
	c       *http.Client
}

func (client *Client) GetTokenFor(nickname string) (token string, err error) {

	resp, err := client.c.Get(client.BaseURL + "/users/" + nickname)

	if err != nil {
		return "", fmt.Errorf("impossible to get data for %s: %s", nickname, err)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("cannot read response body when getting data for %s: %s", nickname, err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got unsuccessful response for %s, status %d: %s", nickname, resp.StatusCode, string(b))
	}

	user := new(User)

	err = json.Unmarshal(b, &user)

	if err != nil {
		return "", fmt.Errorf("cannot unmarshal responded json from users service: %s", err)
	}

	token = user.Token

	if token == "" {
		return "", fmt.Errorf("token for %s user cannot be empty", nickname)
	}

	return token, nil
}

func NewClient() *Client {
	client := &Client{
		BaseURL: "https://blamewarrior.com",
		c:       http.DefaultClient,
	}

	return client
}
