package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	postmark "github.com/mattevans/postmark-go"
)

type emailInput struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
	Tag     string `json:"tag"`
}

type result struct {
	MessageID string `json:"message_id,omitempty"`
	To        string `json:"to,omitempty"`
	Error     string `json:"error,omitempty"`
}

func fatal(msg string) {
	json.NewEncoder(os.Stdout).Encode(result{Error: msg}) //nolint
	os.Exit(1)
}

func main() {
	from := flag.String("from", "", "Sender email address")
	to := flag.String("to", "", "Recipient email address")
	subject := flag.String("subject", "", "Email subject")
	text := flag.String("text", "", "Plain text body")
	html := flag.String("html", "", "HTML body")
	tag := flag.String("tag", "", "Postmark tag (optional)")
	useJSON := flag.Bool("json", false, "Read email fields from stdin as JSON instead of flags")
	setup := flag.Bool("setup", false, "Interactively configure POSTMARK_API_TOKEN in ~/.bashrc")
	flag.Parse()

	if *setup {
		runSetup()
		return
	}

	token := os.Getenv("POSTMARK_API_TOKEN")
	if token == "" {
		fatal("POSTMARK_API_TOKEN not set — run with --setup or: export POSTMARK_API_TOKEN=your_token")
	}

	var in emailInput
	if *useJSON {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fatal("failed to read stdin: " + err.Error())
		}
		if err := json.Unmarshal(data, &in); err != nil {
			fatal("invalid JSON: " + err.Error())
		}
	} else {
		in = emailInput{
			From:    *from,
			To:      *to,
			Subject: *subject,
			Text:    *text,
			HTML:    *html,
			Tag:     *tag,
		}
	}

	if in.From == "" || in.To == "" || in.Subject == "" {
		fatal("--from, --to, and --subject are required (or provide JSON with from/to/subject)")
	}
	if in.Text == "" && in.HTML == "" {
		fatal("provide at least --text or --html body")
	}

	client := postmark.NewClient(
		postmark.WithClient(&http.Client{
			Transport: &postmark.AuthTransport{Token: token},
		}),
	)

	resp, _, err := client.Email.Send(&postmark.Email{
		From:     in.From,
		To:       in.To,
		Subject:  in.Subject,
		TextBody: in.Text,
		HTMLBody: in.HTML,
		Tag:      in.Tag,
	})
	if err != nil {
		fatal(err.Error())
	}

	json.NewEncoder(os.Stdout).Encode(result{ //nolint
		MessageID: resp.MessageID,
		To:        resp.To,
	})
}

func runSetup() {
	fmt.Print("Postmark API token (from https://account.postmarkapp.com/servers → API Tokens): ")
	var token string
	fmt.Scanln(&token)
	if token == "" {
		fmt.Fprintln(os.Stderr, "no token entered, aborting")
		os.Exit(1)
	}

	bashrc := os.Getenv("HOME") + "/.bashrc"
	f, err := os.OpenFile(bashrc, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not open ~/.bashrc:", err)
		os.Exit(1)
	}
	defer f.Close()

	line := fmt.Sprintf("\nexport POSTMARK_API_TOKEN=%q\n", token)
	if _, err := f.WriteString(line); err != nil {
		fmt.Fprintln(os.Stderr, "write failed:", err)
		os.Exit(1)
	}

	fmt.Println("Saved. Run: source ~/.bashrc")
}
