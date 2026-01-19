package mail

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/itsanindyak/email-campaign/types"
)

// MailSend sends an email to the specified recipient using SMTP.
// It uses the SMTP_URL environment variable for the mail server (defaults to localhost:1025)
// and SENDER_EMAIL for the sender address (defaults to iankoley04@gmail.com).
// Returns an error if the email fails to send.
func MailSend(recipient types.Recipient) error {
	smtpURL := os.Getenv("SMTP_URL")
	if smtpURL == "" {
		smtpURL = "localhost:1025"
	}

	senderEmail := os.Getenv("SENDER_EMAIL")
	if senderEmail == "" {
		senderEmail = "iankoley04@gmail.com"
	}

	formattedMsg := fmt.Sprintf("To: %s\r\nSubject: Test Email\r\n\r\n%s\r\n", recipient.Email, "Just testing our email campaign\r\nname: "+recipient.Name)

	msg := []byte(formattedMsg)

	err := smtp.SendMail(smtpURL, nil, senderEmail, []string{recipient.Email}, msg)

	if err != nil {
		return err
	}
	return nil
}
