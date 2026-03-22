package main

import (
	"fmt"
	"time"

	"dailyMailer/config"
	"dailyMailer/llm"
	"dailyMailer/mailer"
)

func main() {
	cfg, err := config.Load("config/config.env")
	if err != nil {
		fmt.Println("❌ Failed to load config:", err)
		return
	}

	fmt.Println("🔁 Running daily mailer job at", time.Now().Format("2006-01-02 15:04:05"))

	statement, err := llm.GenerateStatement(cfg.GroqAPIKey)
	if err != nil {
		fmt.Println("❌ Groq error:", err)
		return
	}
	fmt.Println("📝 Statement:", statement)

	recipients, err := mailer.LoadRecipients("recipients/list.txt")
	if err != nil {
		fmt.Println("❌ Failed to load recipients:", err)
		return
	}
	fmt.Printf("📬 Sending to %d recipient(s)...\n", len(recipients))

	mailer.SendAll(cfg, statement, recipients)
}