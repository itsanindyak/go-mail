package mail

import (
	"context"
	"sync"

	"github.com/itsanindyak/email-campaign/pkg/templates"
)

var (
	once    sync.Once
	mailSvc MailService
	initErr error
)

func initMail() {
	mailSvc, initErr = NewMailService()
}

func Send(
	ctx context.Context,
	recipient string,
	template templates.Template,
) error {

	once.Do(initMail)

	if initErr != nil {
		return initErr
	}

	return mailSvc.sendMail(
		ctx,
		recipient,
		template,
	)
}