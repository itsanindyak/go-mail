package consumer

import (
	"fmt"
	"sync"
	"time"

	"github.com/itsanindyak/email-campaign/mail"
	"github.com/itsanindyak/email-campaign/pkg/metrics"
	"github.com/itsanindyak/email-campaign/types"
)

// EmailWorker processes emails for a single worker.
func EmailWorker(id int, ch chan types.Recipient, dlq chan types.Recipient, wg *sync.WaitGroup) {

	defer wg.Done()

	for recipient := range ch {

		metrics.WorkerActive.Inc()

		fmt.Printf("[Worker %d] Sending email to: %s\n", id, recipient.Email)

		start := time.Now()
		//send mail
		err := mail.Send(recipient)

		duration := time.Since(start).Seconds()
		metrics.EmailDuration.Observe(duration)

		if err != nil {
			fmt.Printf("[Worker %d] Failed to send email to: %s, error: %v\n", id, recipient.Email, err)
			if recipient.Attempts < 3 {
				recipient.Attempts++
				ch <- recipient
			} else {
				metrics.EmailsFailed.Inc()
				dlq <- recipient
			}
		} else {
			metrics.EmailsSent.Inc()
			fmt.Printf("[Worker %d] Send email to: %s\n", id, recipient.Email)
		}

		metrics.WorkerActive.Dec()

	}
}
