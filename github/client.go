package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	HTMLURL   string    `json:"html_url"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	Labels    []struct {
		Name string `json:"name"`
	} `json:"labels"`
	User struct {
		Login string `json:"login"`
	} `json:"user"`
}

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchOpenIssues returns all open issues created after `since`
func (c *Client) FetchOpenIssues(owner, repo string, since time.Time) ([]Issue, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/issues?state=open&since=%s&per_page=50",
		owner, repo, since.UTC().Format(time.RFC3339),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned %d", resp.StatusCode)
	}

	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return issues, nil
}

// HasMatchingLabel returns true if the issue has any of the target labels,
// or if no labels are configured (match everything).
func HasMatchingLabel(issue Issue, targetLabels []string) bool {
	if len(targetLabels) == 0 {
		return true
	}
	for _, tl := range targetLabels {
		for _, il := range issue.Labels {
			if il.Name == tl {
				return true
			}
		}
	}
	return false
}
