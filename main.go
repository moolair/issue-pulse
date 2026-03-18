package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/moolair/issue-pulse/config"
	gh "github.com/moolair/issue-pulse/github"
	"github.com/moolair/issue-pulse/slack"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	dryRun := flag.Bool("dry-run", false, "print alerts without sending to Slack")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ghClient := gh.NewClient(cfg.GitHub.Token)
	notifier := slack.NewNotifier(cfg.Slack.WebhookURL, *dryRun)

	fmt.Printf("👀 Watching %s/%s every %ds\n",
		cfg.GitHub.Owner, cfg.GitHub.Repo, cfg.PollIntervalSeconds)
	if len(cfg.GitHub.Labels) > 0 {
		fmt.Printf("🏷  Label filter: %v\n", cfg.GitHub.Labels)
	}
	if *dryRun {
		fmt.Println("🔍 dry-run mode — no Slack messages will be sent")
	}

	// Track the last time we checked so we don't re-alert old issues
	lastChecked := time.Now().UTC()

	for {
		time.Sleep(time.Duration(cfg.PollIntervalSeconds) * time.Second)

		issues, err := ghClient.FetchOpenIssues(cfg.GitHub.Owner, cfg.GitHub.Repo, lastChecked)
		if err != nil {
			log.Printf("error fetching issues: %v", err)
			continue
		}

		now := time.Now().UTC()

		for _, issue := range issues {
			// Only alert on genuinely new issues (created after our last check)
			if !issue.CreatedAt.After(lastChecked) {
				continue
			}
			if !gh.HasMatchingLabel(issue, cfg.GitHub.Labels) {
				continue
			}

			labels := make([]string, len(issue.Labels))
			for i, l := range issue.Labels {
				labels[i] = l.Name
			}

			log.Printf("new issue #%d: %s", issue.Number, issue.Title)

			if err := notifier.Send(
				issue.Number,
				issue.Title,
				issue.HTMLURL,
				issue.User.Login,
				labels,
				issue.CreatedAt,
			); err != nil {
				log.Printf("failed to send Slack alert: %v", err)
			}
		}

		lastChecked = now
	}
}
