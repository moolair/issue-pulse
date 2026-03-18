package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type payload struct {
	Text        string       `json:"text"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	TitleLink string `json:"title_link"`
	Fields []field `json:"fields"`
	Footer string  `json:"footer"`
	Ts     int64   `json:"ts"`
}

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Notifier struct {
	webhookURL string
	dryRun     bool
	httpClient *http.Client
}

func NewNotifier(webhookURL string, dryRun bool) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		dryRun:     dryRun,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (n *Notifier) Send(issueNumber int, title, url, author string, labels []string, createdAt time.Time) error {
	labelStr := "none"
	if len(labels) > 0 {
		labelStr = ""
		for i, l := range labels {
			if i > 0 {
				labelStr += ", "
			}
			labelStr += l
		}
	}

	p := payload{
		Text: fmt.Sprintf(":bell: *New GitHub Issue #%d*", issueNumber),
		Attachments: []attachment{
			{
				Color:     "#e74c3c",
				Title:     title,
				TitleLink: url,
				Fields: []field{
					{Title: "Author", Value: author, Short: true},
					{Title: "Labels", Value: labelStr, Short: true},
				},
				Footer: "issue-pulse",
				Ts:     createdAt.Unix(),
			},
		},
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	if n.dryRun {
		fmt.Printf("[dry-run] Would send to Slack:\n%s\n\n", string(body))
		return nil
	}

	resp, err := n.httpClient.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned %d", resp.StatusCode)
	}
	return nil
}
