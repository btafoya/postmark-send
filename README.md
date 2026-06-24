# postmark-send

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/btafoya/postmark-send)](go.mod)
[![Postmark](https://img.shields.io/badge/Postmark-API-blue)](https://postmarkapp.com)
[![Built with Go](https://img.shields.io/badge/Built%20with-Go-00ADD8.svg)](https://go.dev)

CLI tool for sending email via [Postmark](https://www.postmarkapp.com/?via=tafoya). Built to be called by AI agents — all output is newline-delimited JSON, errors included, exit codes are reliable.

---

## Requirements

- Go 1.21+
- A Postmark account with a Server API Token

## Install

The fastest way is with `go install`:

```bash
go install github.com/btafoya/postmark-send@v0.1.0
```

The binary is installed to `$GOPATH/bin` or `~/go/bin` by default. Make sure that directory is on your `PATH`.

Or install the latest untagged commit:

```bash
GOPROXY=direct go install github.com/btafoya/postmark-send@latest
```

Or clone and build from source:

```bash
git clone https://github.com/btafoya/postmark-send
cd postmark-send
go build -o postmark-send .
```

## Setup

Run the interactive setup to save your API token to `~/.bashrc`:

```bash
./postmark-send --setup
# Postmark API token: <paste token>
# Saved. Run: source ~/.bashrc
source ~/.bashrc
```

Or set it manually:

```bash
export POSTMARK_API_TOKEN=your_server_token
```

---

## Usage

### Flags (human / scripted use)

```bash
./postmark-send \
  --from sender@example.com \
  --to recipient@example.com \
  --subject "Hello" \
  --text "Plain text body" \
  --html "<b>HTML body</b>" \
  --tag onboarding
```

`--from`, `--to`, and `--subject` are required. At least one of `--text` or `--html` is required. `--tag` is optional.

### JSON stdin (AI agent use)

Pass a single JSON object on stdin with `--json`:

```bash
echo '{
  "from": "sender@example.com",
  "to": "recipient@example.com",
  "subject": "Hello",
  "text": "Plain text body",
  "html": "<b>HTML body</b>",
  "tag": "onboarding"
}' | ./postmark-send --json
```

---

## Output

All output is JSON on stdout. **Never parse stderr** — errors are on stdout as `{"error":"..."}`.

**Success:**
```json
{"message_id":"abc-123-def","to":"recipient@example.com"}
```

**Error:**
```json
{"error":"--from, --to, and --subject are required"}
```

Exit code `0` = success. Exit code `1` = error (check the `error` field).

---

## For AI Agents

Use the `--json` flag and pipe a JSON object on stdin. Check the output for an `error` field before assuming success.

### Send an email

```bash
echo '{"from":"you@example.com","to":"them@example.com","subject":"Hello","text":"Hi there"}' \
  | ./postmark-send --json
```

### Check the result

```bash
result=$(echo '{"from":"you@example.com","to":"them@example.com","subject":"Hello","text":"Hi there"}' \
  | ./postmark-send --json)

if echo "$result" | grep -q '"error"'; then
  echo "Failed: $result"
else
  echo "Sent. Message ID: $(echo "$result" | grep -o '"message_id":"[^"]*"')"
fi
```

### With HTML body

```bash
echo '{
  "from": "you@example.com",
  "to": "them@example.com",
  "subject": "Hello",
  "text": "Hi there",
  "html": "<p>Hi there</p>",
  "tag": "welcome"
}' | ./postmark-send --json
```

### Expected output

Success:
```json
{"message_id":"abc-123","to":"them@example.com"}
```

Failure:
```json
{"error":"POSTMARK_API_TOKEN not set — run with --setup or: export POSTMARK_API_TOKEN=your_token"}
```

Exit code `0` = sent. Exit code `1` = not sent.

---

## Flags reference

| Flag | Description |
|------|-------------|
| `--from` | Sender email address |
| `--to` | Recipient email address |
| `--subject` | Email subject line |
| `--text` | Plain text body |
| `--html` | HTML body |
| `--tag` | Postmark message tag (optional) |
| `--json` | Read all fields from JSON on stdin instead of flags |
| `--setup` | Interactive setup: writes `POSTMARK_API_TOKEN` to `~/.bashrc` |

## JSON schema (stdin mode)

```json
{
  "from":    "string (required)",
  "to":      "string (required)",
  "subject": "string (required)",
  "text":    "string (required if no html)",
  "html":    "string (required if no text)",
  "tag":     "string (optional)"
}
```

---

## About the Author

**Brian Tafoya** — Senior Software Engineer | Systems Architect | Project Manager

Results-driven Senior Software Engineer with over 30 years of experience architecting, developing, and maintaining robust, scalable systems. Proven expertise across full-stack web development, multiple programming languages, and both Linux and Microsoft ecosystems, including development with .NET and administration of Windows Server environments. Extensive background in systems and network administration, with a strong focus on security, performance, and operational reliability. A veteran of the Internet industry since 1993, delivering end-to-end solutions to complex technical challenges with precision and efficiency.

[github.com/btafoya](https://github.com/btafoya/) · [briantafoya.com](https://briantafoya.com/) - [tafoyaventures.com](https://tafoyaventures.com/)

---

## License

MIT — see [LICENSE](LICENSE).

Built on [mattevans/postmark-go](https://github.com/mattevans/postmark-go).
