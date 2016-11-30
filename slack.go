package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SlackMessage struct {
	Channel   string `json:"channel" schema:"channel"`
	Username  string `json:"username" schema:"username"`
	IconEmoji string `json:"icon_emoji,omitempty" schema:"icon_emoji"`
	IconURL   string `json:"icon_url,omitempty" schema:"icon_url"`
	Text      string `json:"text" schema:"text"`
}

type Slack struct {
	url string
}

func (s *Slack) Send(message SlackMessage) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Response: %d %s", resp.StatusCode, string(b))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("Slack returned status code %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
