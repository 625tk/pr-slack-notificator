package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type GithubResponse struct {
	Body string `json:"body"`
}

type SlackPayload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func main() {
	repo := os.Getenv("REPOSITORY")
	token := os.Getenv("GIT_PR_RELEASE_TOKEN")
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	sChannel := os.Getenv("SLACK_CHANNEL")

	prNumber, err := strconv.Atoi(os.Getenv("PR_NUMBER"))
	if err != nil {
		log.Fatal("failed to parse env PR_NUMBER, ", err)
	}

	ctx := context.Background()
	message, err := getPRBody(ctx, repo, token, prNumber)
	if err != nil {
		log.Fatal(err)
	}

	p := SlackPayload{
		Channel: sChannel,
		Text:    message,
	}

	b, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	_, err = http.Post(webhookURL, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
}

func getPRBody(ctx context.Context, repository, token string, prNumber int) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d", repository, prNumber)
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url,
		nil,
	)
	if err != nil {
		log.Println("failed to create request")
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("failed to request")
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("failed to read body")
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code != ok %d, %s", res.StatusCode, string(b))
	}

	r := GithubResponse{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		log.Println("failed to unmarshal json")
		return "", err
	}

	fmt.Println(r)
	sp := strings.Split(r.Body, "```")
	if len(sp) < 2 {
		return "", errors.New("failed to split quotes")
	}
	return "```" + sp[1] + "```", nil
}
