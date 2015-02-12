package leeroy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/fabioxgn/go-bot"
)

var (
	repo string = "docker/docker"
	url  string = "https://leeroy.dockerproject.com/build/retry"
)

type PullRequest struct {
	Number int    `json:"number"`
	Repo   string `json:"repo"`
}

func rebuild(command *bot.Cmd) (msg string, err error) {
	if len(command.Args) < 1 {
		return "Not enough args. try: !rebuild 9437 docker/docker", nil
	}
	if len(command.Args) >= 2 {
		repo = command.Args[1]
	}

	// convert the string to an int
	num, err := strconv.Atoi(command.Args[0])
	if err != nil {
		return "", fmt.Errorf("converting PR num to int failed: %v", err)
	}

	pr := PullRequest{
		Number: num,
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
	req.SetBasicAuth(os.Getenv("LEEROY_USERNAME"), os.Getenv("LEEROY_PASS"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Sprintf("Rebuilding PR %d for %s returned status code: %d", pr.Number, repo, resp.StatusCode), nil
	}

	return fmt.Sprintf("Rebuilding PR %d at https://github.com/%s/pull/%d", pr.Number, repo, pr.Number), nil
}

func init() {
	bot.RegisterCommand(
		"rebuild",
		"Rebuilds a PR number on Jenkins.",
		"9437",
		rebuild)
}
