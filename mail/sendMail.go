package mail

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/itsanindyak/email-campaign/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter           = otel.Meter("email-engine")
	sendDuration, _ = meter.Float64Histogram(
		"smtp_send_duration_seconds",
		metric.WithDescription("Time taken to establish SMTP connection and send"),
		// Adding a description helps your team understand it in Grafana
	)
)

// MailSend sends an email to the specified recipient using SMTP.
// It uses the SMTP_URL environment variable for the mail server (defaults to localhost:1025)
// and SENDER_EMAIL for the sender address (defaults to iankoley04@gmail.com).
// Returns an error if the email fails to send.
func MailSend(ctx context.Context, recipient types.Recipient) error {

	tracer := otel.Tracer("email-engine")

	ctx, span := tracer.Start(ctx,"smtp.Sendmail")

	defer span.End()

	// Tag the trace with searchable data
	span.SetAttributes(attribute.String("recipient.email",recipient.Email))

	smtpURL := os.Getenv("SMTP_URL")
	if smtpURL == "" {
		smtpURL = "localhost:1025"
	}

	senderEmail := os.Getenv("SENDER_EMAIL")
	if senderEmail == "" {
		senderEmail = "iankoley04@gmail.com"
	}

	startTime := time.Now()

	formattedMsg := fmt.Sprintf("To: %s\r\nSubject: Test Email\r\n\r\n%s\r\n", recipient.Email, "Just testing our email campaign\r\nname: "+recipient.Name)

	msg := []byte(formattedMsg)
	time.Sleep(100 * time.Millisecond)

	err := smtp.SendMail(smtpURL, nil, senderEmail, []string{recipient.Email}, msg)

	sendDuration.Record(ctx,time.Since(startTime).Seconds())

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error,"SMTP failure")
		return err
	}

	span.SetStatus(codes.Ok,"Email sent")
	return nil
}
