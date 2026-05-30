package mail

import (
	"context"

	"github.com/itsanindyak/email-campaign/pkg/templates"
)

type MailService interface {
	sendMail(ctx context.Context, recipient string, template templates.Template) error
}