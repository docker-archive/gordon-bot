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
	repo    string = "docker/docker"
	baseurl string = "https://leeroy.dockerproject.com/"
)

type PullRequest struct {
	Number  int    `json:"number"`
	Repo    string `json:"repo"`
	Context string `json:"context"`
}

func sendRequest(pr PullRequest, url string) (err error) {
	data, err := json.Marshal(pr)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("LEEROY_USERNAME"), os.Getenv("LEEROY_PASS"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("Requesting %s for PR %d for %s returned status code: %d", url, pr.Number, repo, resp.StatusCode)
	}

	return nil
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

	if err := sendRequest(pr, baseurl+"build/retry"); err != nil {
		return "", err
	}

	return fmt.Sprintf("Rebuilding PR %d at https://github.com/%s/pull/%d", pr.Number, repo, pr.Number), nil
}

func customBuild(command *bot.Cmd, context string) (msg string, err error) {
	if len(command.Args) < 1 {
		return fmt.Sprintf("Not enough args. try: !%s 9437 docker/docker", context), nil
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
		Number:  num,
		Repo:    repo,
		Context: context,
	}

	if err := sendRequest(pr, baseurl+"build/custom"); err != nil {
		return "", err
	}

	return fmt.Sprintf("Building PR %d on %s at https://github.com/%s/pull/%d", pr.Number, context, repo, pr.Number), nil
}

func lxcBuild(command *bot.Cmd) (msg string, err error) {
	return customBuild(command, "lxc")
}

func windowsBuild(command *bot.Cmd) (msg string, err error) {
	return customBuild(command, "windows")
}

func init() {
	bot.RegisterCommand(
		"rebuild",
		"Rebuilds a PR number on Jenkins.",
		"9437",
		rebuild)

	bot.RegisterCommand(
		"lxc",
		"Build a PR on an lxc box with the lxc driver",
		"9345",
		lxcBuild,
	)

	bot.RegisterCommand(
		"windows",
		"Build a PR on a windows box",
		"9345",
		windowsBuild,
	)
}
