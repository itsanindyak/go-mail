package mail

import (
	"fmt"
	"log"
	"os"

	"github.com/itsanindyak/email-campaign/types"
	"github.com/resend/resend-go/v3"
)

// Send sends a welcome email to the specified recipient using the Resend API.
//
// It requires the RESEND_API_KEY environment variable to be set.
// The email is sent using a predefined template with the recipient's name and email as variables.
// Returns an error if the API key is missing or if the email fails to send.
func Send(recipient types.Recipient) error {
	apiKey := os.Getenv("RESEND_API_KEY")

	if apiKey == "" {
		log.Printf("RESEND_API_KEY environment variable is not set")
		return fmt.Errorf("RESEND_API_KEY not configured")
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@itsak.xyz",
		To:      []string{recipient.Email},
		Subject: fmt.Sprintf("Welcome to %s, %s", "itsak", recipient.Name),
		Template: &resend.EmailTemplate{
			Id: "ce55501d-09c2-4b6f-aa26-7c61092e8f32",
			Variables: map[string]any{
				"NAME": recipient.Name,
				"MAIL": recipient.Email,
			},
		},
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
