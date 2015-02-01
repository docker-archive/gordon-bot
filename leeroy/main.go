package leeroy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fabioxgn/go-bot"
)

var (
	repo string = "docker/docker"
	url  string = "https://leeroy.dockerproject.com/retry"
)

type PullRequest struct {
	Number string `json:"number"`
	Repo   string `json:"repo"`
}

func rebuild(command *bot.Cmd) (msg string, err error) {
	if len(command.Args) < 1 {
		return "Not enough args. try: !rebuild 9437 docker/docker", nil
	}
	if len(command.Args) >= 2 {
		repo = command.Args[1]
	}

	pr := PullRequest{
		Number: command.Args[0],
		Repo:   repo,
	}
	data, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Sprintf("Rebuilding PR %s for %s returned status code: %d", pr.Number, repo, resp.StatusCode), nil
	}

	return fmt.Sprintf("Rebuilding PR %s at https://github.com/%s/pull/%s", pr.Number, repo, pr.Number), nil
}

func init() {
	bot.RegisterCommand(
		"rebuild",
		"Rebuilds a PR number on Jenkins.",
		"9437",
		rebuild)
}
