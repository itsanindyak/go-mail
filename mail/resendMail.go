package mail

import (
	"context"
	"fmt"

	"github.com/itsanindyak/email-campaign/pkg/templates"
	"github.com/resend/resend-go/v3"
)

type ResendService struct {
    APIKey string
	client *resend.Client
}

func NewResendService(apiKey string) *ResendService {
	return &ResendService{
		client: resend.NewClient(apiKey),
	}
}


// Send sends a welcome email to the specified recipient using the Resend API.
//
// It requires the RESEND_API_KEY environment variable to be set.
// The email is sent using a predefined template with the recipient's name and email as variables.
// Returns an error if the API key is missing or if the email fails to send.
func (r *ResendService)sendMail(ctx context.Context, recipient string, template templates.Template) error {


	html,subject,err := templates.Render(template)

	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{recipient},
		Subject: subject,
		Html: html,
	}

	_, err = r.client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
