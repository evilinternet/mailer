package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	GroqAPIKey  string
	SMTPHost    string
	SMTPPort    string
	SMTPLogin   string
	SMTPKey     string
	SenderName  string
	SenderEmail string
	Subject     string
}

func Load(path string) (*Config, error) {
	env := make(map[string]string)

	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimPrefix(line, "\xef\xbb\xbf")
			line = strings.TrimRight(line, "\r")
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, `"'`)
				env[key] = val
			}
		}
	}

	for _, key := range []string{
		"GROQ_API_KEY", "BREVO_SMTP_HOST", "BREVO_SMTP_PORT",
		"BREVO_LOGIN", "BREVO_SMTP_KEY", "SENDER_NAME",
		"SENDER_EMAIL", "SUBJECT",
	} {
		if val := os.Getenv(key); val != "" {
			env[key] = val
		}
	}

	return &Config{
		GroqAPIKey:  env["GROQ_API_KEY"],
		SMTPHost:    env["BREVO_SMTP_HOST"],
		SMTPPort:    env["BREVO_SMTP_PORT"],
		SMTPLogin:   env["BREVO_LOGIN"],
		SMTPKey:     env["BREVO_SMTP_KEY"],
		SenderName:  env["SENDER_NAME"],
		SenderEmail: env["SENDER_EMAIL"],
		Subject:     env["SUBJECT"],
	}, nil
}
