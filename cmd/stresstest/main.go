package main

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

func main() {
	addTeamRequest := Request{
		Method: "POST",
		URL:    "http://localhost:8080/team/add",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `
			{
			"team_name": "payments",
			"members": [
				{
				"user_id": "u1",
				"username": "Alice",
				"is_active": true
				},
				{
				"user_id": "u2",
				"username": "Bob",
				"is_active": true
				}
			]
			}
		`,
	}

	getTeamRequest := Request{
		Method: "GET",
		URL:    "http://localhost:8080/team/get?team_name=payments",
	}

	setIsActiveRequest := Request{
		Method: "POST",
		URL:    "http://localhost:8080/users/setIsActive",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `
			{
			"user_id": "u2",
			"is_active": false
			}
		`,
	}

	getReviewRequest := Request{
		Method: "GET",
		URL:    "http://localhost:8080/users/getReview?user_id=u1",
	}

	prCreateRequest := Request{
		Method: "POST",
		URL:    "http://localhost:8080/pullRequest/create",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `
			{
			"pull_request_id": "pr-1001",
			"pull_request_name": "Add search",
			"author_id": "u1"
			}
		`,
	}

	prMergeRequest := Request{
		Method: "POST",
		URL:    "http://localhost:8080/pullRequest/merge",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `
			{
			"pull_request_id": "pr-1001"
			}
		`,
	}

	prReassignRequest := Request{
		Method: "POST",
		URL:    "http://localhost:8080/pullRequest/reassign",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `
			{
			"pull_request_id": "pr-1001",
			"old_reviewer_id": "u2"
			}
		`,
	}

	var requests []Request = []Request{
		addTeamRequest, getTeamRequest,
		setIsActiveRequest, getReviewRequest,
		prCreateRequest, prMergeRequest, prReassignRequest,
	}

	targeter := func(t *vegeta.Target) error {
		request := requests[rand.IntN(len(requests))]

		t.Method = request.Method
		t.URL = request.URL
		if request.Method == "POST" {
			t.Body, _ = json.Marshal(request.Body)
			t.Header = http.Header{
				"Content-Type": []string{"application/json"},
			}
		}

		return nil
	}

	rate := vegeta.Rate{Freq: 5, Per: time.Second}
	duration := 30 * time.Second

	attacker := vegeta.NewAttacker(
		vegeta.Timeout(10*time.Second),
		vegeta.Workers(20),
		vegeta.MaxBody(1024*1024),
	)

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	report := vegeta.NewTextReporter(&metrics)
	report.Report(os.Stdout)
}
