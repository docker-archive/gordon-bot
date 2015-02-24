package leeroy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fabioxgn/go-bot"
)

var (
	RepoPrefix string = "docker/"
	BaseUrl    string = "https://leeroy.dockerproject.com/"
)

type PullRequest struct {
	Number  int    `json:"number"`
	Repo    string `json:"repo"`
	Context string `json:"context"`
}

func parsePullRequest(arg string) (pr PullRequest, err error) {
	// parse for the repo
	// split on #
	nameArgs := strings.SplitN(arg, "#", 2)
	if len(nameArgs) <= 1 {
		return pr, fmt.Errorf("%s did not include #", arg)
	}

	pr.Repo = RepoPrefix + nameArgs[0]

	// parse the second arguement for a /
	// for if its a custom build
	buildArgs := strings.SplitN(nameArgs[1], "/", 2)
	if len(buildArgs) == 2 {
		pr.Context = buildArgs[1]
	}

	// parse as int
	pr.Number, err = strconv.Atoi(buildArgs[0])
	if err != nil {
		return pr, err
	}

	return pr, nil
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
		return fmt.Errorf("Requesting %s for PR %d for %s returned status code: %d. Make sure the repo allows builds.", url, pr.Number, pr.Repo, resp.StatusCode)
	}
	return nil
}

func rebuild(command *bot.Cmd) (msg string, err error) {
	tryString := "Try !rebuild libcontainer#234."
	if len(command.Args) < 1 {
		return "", fmt.Errorf("Not enough args. %s", tryString)
	}

	pr, err := parsePullRequest(command.Args[0])
	if err != nil {
		return "", fmt.Errorf("Error parsing pull request: %v. %s", err, tryString)
	}

	endpoint := "build/retry"
	if pr.Context != "" {
		endpoint = "build/custom"
	}

	if err := sendRequest(pr, BaseUrl+endpoint); err != nil {
		return "", err
	}

	if pr.Context != "" {
		return fmt.Sprintf("Building PR %d on %s at https://github.com/%s/pull/%d", pr.Number, pr.Context, pr.Repo, pr.Number), nil
	}

	return fmt.Sprintf("Rebuilding PR %d at https://github.com/%s/pull/%d", pr.Number, pr.Repo, pr.Number), nil
}

func init() {
	bot.RegisterCommand(
		"rebuild",
		"Rebuilds a PR number on Jenkins.",
		"libcontainer#9437",
		rebuild)
}
