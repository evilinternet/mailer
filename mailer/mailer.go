package mailer

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"dailyMailer/config"
)

type emailData struct {
	Statement  string
	SenderName string
}

func LoadRecipients(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var recipients []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			recipients = append(recipients, line)
		}
	}
	return recipients, nil
}

func SendAll(cfg *config.Config, statement string, recipients []string) {
	tmpl, err := template.ParseFiles("templates/email.html")
	if err != nil {
		fmt.Println("❌ Failed to load email template:", err)
		return
	}

	auth := smtp.PlainAuth("", cfg.SMTPLogin, cfg.SMTPKey, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	for _, recipient := range recipients {
		var body bytes.Buffer
		if err := tmpl.Execute(&body, emailData{
			Statement:  statement,
			SenderName: cfg.SenderName,
		}); err != nil {
			fmt.Printf("❌ Template error for %s: %v\n", recipient, err)
			continue
		}

		msg := buildMessage(cfg, recipient, body.String())

		if err := smtp.SendMail(addr, auth, cfg.SMTPLogin, []string{recipient}, []byte(msg)); err != nil {
			fmt.Printf("❌ Failed to send to %s: %v\n", recipient, err)
		} else {
			fmt.Printf("✅ Sent to %s\n", recipient)
		}
	}
}

func buildMessage(cfg *config.Config, to, htmlBody string) string {
	from := fmt.Sprintf("%s <%s>", cfg.SenderName, cfg.SenderEmail)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: %s\r\n", from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", to))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", cfg.Subject))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(htmlBody)
	return sb.String()
}