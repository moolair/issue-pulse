# issue-pulse

> A lightweight CLI that watches a GitHub repository for new issues and sends real-time Slack alerts — with label filtering and dry-run support.

Built to mirror the kind of support tooling used in technical support workflows (GitHub Issues → team notification → faster response time).

---

## Why I Built This

In technical support roles, knowing about new issues the moment they're filed is the difference between a 5-minute response and a 50-minute one. `issue-pulse` automates that signal — no dashboards to check, no manual polling.

---

## Demo

```
$ go run . --config config.json --dry-run

👀 Watching moolair/demo-repo every 60s
🏷  Label filter: [bug customer-reported]
🔍 dry-run mode — no Slack messages will be sent

[dry-run] Would send to Slack:
{
  "text": ":bell: *New GitHub Issue #42*",
  "attachments": [{
    "color": "#e74c3c",
    "title": "Login fails for SSO users on Firefox",
    "title_link": "https://github.com/moolair/demo-repo/issues/42",
    "fields": [
      { "title": "Author", "value": "customer-abc", "short": true },
      { "title": "Labels", "value": "bug, customer-reported", "short": true }
    ],
    "footer": "issue-pulse"
  }]
}
```

---

## Setup

### 1. Clone & install dependencies

```bash
git clone https://github.com/moolair/issue-pulse.git
cd issue-pulse
go mod tidy
```

### 2. Configure

```bash
cp config.json.example config.json
```

Edit `config.json`:

```yaml
github:
  token: "ghp_YOUR_PERSONAL_ACCESS_TOKEN"   # needs repo:read scope
  owner: "your-github-username"
  repo: "your-repo-name"
  labels:                                    # leave empty to watch ALL issues
    - "bug"
    - "customer-reported"

slack:
  webhook_url: "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

poll_interval_seconds: 60
```

**GitHub token:** [Settings → Developer settings → Personal access tokens](https://github.com/settings/tokens) — enable `repo` (read) scope.

**Slack webhook:** [Create an Incoming Webhook](https://api.slack.com/messaging/webhooks) for your workspace.

### 3. Run

```bash
# Live mode
go run . --config config.json

# Dry-run (prints to stdout, no Slack messages sent)
go run . --config config.json --dry-run

# Build binary
go build -o issue-pulse .
./issue-pulse --config config.json
```

---

## How It Works

```
┌─────────────┐     poll every N seconds     ┌──────────────────┐
│  GitHub API  │ ──────────────────────────▶ │   issue-pulse    │
│  (Issues)   │                              │                  │
└─────────────┘                              │  • label filter  │
                                             │  • dedup by time │
                                             └────────┬─────────┘
                                                      │ new issue found
                                                      ▼
                                             ┌──────────────────┐
                                             │   Slack Webhook  │
                                             │  (rich message)  │
                                             └──────────────────┘
```

1. Polls GitHub Issues API every `poll_interval_seconds`
2. Filters by labels (if configured)
3. Deduplicates using a `lastChecked` timestamp — no repeat alerts
4. Sends a rich Slack message with issue title, author, labels, and direct link

---

## Project Structure

```
issue-pulse/
├── main.go              # CLI entrypoint, poll loop
├── config/
│   └── config.go        # YAML config loading
├── github/
│   └── client.go        # GitHub Issues API client + label filtering
├── slack/
│   └── notifier.go      # Slack Incoming Webhook sender
├── config.json.example  # Config template
└── README.md
```

---

## Future Ideas

- [ ] Persist `lastChecked` to disk (survive restarts)
- [ ] Support multiple repos in one config
- [ ] Add Jira ticket auto-creation on new issue
- [ ] Docker image for easy deployment

---

## License

MIT
