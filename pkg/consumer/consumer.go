package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/itsanindyak/email-campaign/mail"
	"github.com/itsanindyak/email-campaign/types"
	"golang.org/x/time/rate"
)

// EmailWorker processes emails from the recipient channel using the specified worker ID.
// For each recipient, it attempts to send an email and handles success/failure appropriately:
//   - On success: increments the sent counter metric
//   - On failure with attempts < 3: re-queues to retry channel with incremented attempt count
//   - On failure with attempts >= 3: moves recipient to dead-letter queue
//
// The worker signals completion via the provided WaitGroup when the channel is closed.
func EmailWorker(ctx context.Context, id int, ch chan types.Recipient, dlq chan types.Recipient, wg *sync.WaitGroup, limiter *rate.Limiter) {

	for {

		select {

		case <-ctx.Done():
			return

		case recipient := <-ch:

			err := limiter.Wait(ctx)
			if err != nil {
				return
			}

			fmt.Printf("[Worker %d] Sending email to: %s\n", id, recipient.Email)

			//send mail
			//  err := mail.Send(recipient)

			err = mail.MailSend(recipient)

			
			if err != nil {
				fmt.Printf("[Worker %d] Failed to send email to: %s, error: %v\n", id, recipient.Email, err)
				if recipient.Attempts < 3 {
					recipient.Attempts++

					go func(r types.Recipient) {
						time.Sleep(time.Duration(r.Attempts) * time.Second)
						select {
						case ch <- r:

							// Success: Retry queued
							fmt.Printf("Retry queued for %s\n", r.Email)
						case <-ctx.Done():
							wg.Done()
						}
					}(recipient)

				} else {

					select {
					case dlq <- recipient:
					case <-ctx.Done():
					}

					wg.Done()
				}
			} else {
				fmt.Printf("[Worker %d] Send email to: %s\n", id, recipient.Email)

				wg.Done()
			}

		}
	}
}
