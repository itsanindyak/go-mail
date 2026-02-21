package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/itsanindyak/email-campaign/mail"
	"github.com/itsanindyak/email-campaign/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/time/rate"
)

var (
	meter           = otel.Meter("email-engine")
	sentCounter, _  = meter.Int64Counter("email_sent_success_total")
	failCounter, _  = meter.Int64Counter("email_sent_failed_total")
	retryCounter, _ = meter.Int64Counter("email_retry_queued_total")
	dlqCounter, _   = meter.Int64Counter("email_sent_to_dlq_total")
)

// EmailWorker processes emails from the recipient channel using the specified worker ID.
// For each recipient, it attempts to send an email and handles success/failure appropriately:
//   - On success: increments the sent counter metric
//   - On failure with attempts < 3: re-queues to retry channel with incremented attempt count
//   - On failure with attempts >= 3: moves recipient to dead-letter queue
//
// The worker signals completion via the provided WaitGroup when the channel is closed.
func EmailWorker(ctx context.Context, id int, ch chan types.Recipient, dlq chan types.Recipient, wg *sync.WaitGroup, limiter *rate.Limiter) {

	tracer := otel.Tracer("email-engine")

	for {

		select {

		case <-ctx.Done():
			return

		case recipient := <-ch:

			err := limiter.Wait(ctx)
			if err != nil {
				return
			}

			jobCtx, span := tracer.Start(ctx, "consumer.sendMail")

			span.SetAttributes(
				attribute.Int("worker.id", id),
				attribute.Int("attempt.number", recipient.Attempts),
				attribute.String("recipient.email", recipient.Email),
			)

			// fmt.Printf("[Worker %d] Sending email to: %s\n", id, recipient.Email)

			//send mail
			// err = mail.Send(recipient)

			err = mail.MailSend(jobCtx, recipient)

			if err != nil {
				fmt.Printf("[Worker %d] Failed to send email to: %s, error: %v\n", id, recipient.Email, err)

				failCounter.Add(jobCtx, 1)

				span.RecordError(err)
				span.SetStatus(codes.Error, "Failed to send email")

				slog.ErrorContext(jobCtx, "Failed to send mail", "recipient", recipient.Email, "error", err)

				if recipient.Attempts < 3 {
					recipient.Attempts++
					retryCounter.Add(jobCtx,1)

					go func(r types.Recipient) {
						time.Sleep(time.Duration(r.Attempts) * time.Second)
						select {
						case ch <- r:
							slog.Info("Retry queued", "email", r.Email)
						case <-ctx.Done():
							wg.Done()
						}
					}(recipient)

				} else {
					dlqCounter.Add(jobCtx,1)
					slog.WarnContext(jobCtx,"Max attempts reached, moving to DLQ")

					select {
					case dlq <- recipient:
					case <-ctx.Done():
					}

					wg.Done()
				}
			} else {
				sentCounter.Add(jobCtx,1)
				slog.InfoContext(jobCtx, "Email successfully sent", "worker_id", id)
				fmt.Printf("[Worker %d] Send email to: %s\n", id, recipient.Email)

				wg.Done()
			}

			span.End()

		}
	}
}
