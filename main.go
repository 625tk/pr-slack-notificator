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
	"time"
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
	githubToken := os.Getenv("GIT_PR_RELEASE_TOKEN")
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	sChannel := os.Getenv("SLACK_CHANNEL")
	slackToken := os.Getenv("SLACK_API_TOKEN")

	prNumber, err := strconv.Atoi(os.Getenv("PR_NUMBER"))
	if err != nil {
		log.Fatal("failed to parse env PR_NUMBER, ", err)
	}

	ctx := context.Background()
	message, err := getPRBody(ctx, repo, githubToken, prNumber)
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
	if webhookURL != "" {
		postViaWebhook(bytes.NewReader(b), webhookURL)
	} else if slackToken != "" {
		postViaAPI(bytes.NewReader(b), slackToken)
	}

}

func postViaWebhook(reader *bytes.Reader, url string) {
	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	a, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}

func postViaAPI(reader *bytes.Reader, token string) {
	url := "https://slack.com/api/chat.postMessage"

	r, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Fatal(err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	_, err = http.DefaultClient.Do(r)
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
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	f := time.Now().In(jst).Format("2006-01-02 15:04")

	replaced := strings.Replace(sp[1], "[TIME]", f, -1)

	prURL := fmt.Sprintf("\n ref: https://github.com/%s/pull/%d \n", repository, prNumber)

	return "```" + replaced + prURL + "```", nil

}
