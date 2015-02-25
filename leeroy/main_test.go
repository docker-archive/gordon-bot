package leeroy

import (
	"strings"
	"testing"
)

type TestPR struct {
	ShouldFail bool
	PR         PullRequest
	ErrorMatch string
}

func TestParsePullRequest(t *testing.T) {
	tests := map[string]TestPR{
		"1234": TestPR{
			ShouldFail: true,
			ErrorMatch: "did not include #",
		},
		"docker#1234": TestPR{
			ShouldFail: false,
			PR: PullRequest{
				Number: 1234,
				Repo:   RepoPrefix + "docker",
			},
		},
		"docker#23/lxc": TestPR{
			ShouldFail: false,
			PR: PullRequest{
				Number:  23,
				Repo:    RepoPrefix + "docker",
				Context: "lxc",
			},
		},
		"libcontainer#456/windows": TestPR{
			ShouldFail: true,
			ErrorMatch: "Custom builds are not currently allowed on docker/libcontainer, only docker/docker",
		},
		"docker123": TestPR{
			ShouldFail: true,
			ErrorMatch: "did not include #",
		},
		"docker": TestPR{
			ShouldFail: true,
			ErrorMatch: "did not include #",
		},
	}

	for testString, testPR := range tests {
		pr, err := parsePullRequest(testString)

		if testPR.ShouldFail && (err == nil || !strings.Contains(err.Error(), testPR.ErrorMatch)) {
			t.Fatalf("Parsing %s should have failed with %s: %v", testString, testPR.ErrorMatch, err)
		}

		if !testPR.ShouldFail && err != nil {
			t.Fatalf("Parsing %s should not have failed: %v", testString, err)
		}

		if !testPR.ShouldFail && pr != testPR.PR {
			t.Fatalf("Parsing %s should have returned %#v, got: %#v", testString, pr, testPR.PR)
		}

	}

}
