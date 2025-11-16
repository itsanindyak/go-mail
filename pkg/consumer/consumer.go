package consumer

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/itsanindyak/email-campaign/mail"
	"github.com/itsanindyak/email-campaign/pkg/metrics"
	"github.com/itsanindyak/email-campaign/types"
)

func EmailWorker(id int, ch chan types.Recipient,dlq chan types.Recipient,wg *sync.WaitGroup) {
	
	defer wg.Done()

	for recipient := range ch {

		metrics.WorkerActive.Inc()
		
		success := false

		for attempt := 1; attempt <= 3;attempt++ {

			fmt.Printf("[Worker %d] Sending email to: %s\n", id, recipient.Email)

			start := time.Now()
			err := mail.MailSend(recipient)
			duration := time.Since(start).Seconds()

			metrics.EmailDuration.Observe(duration)

			if err == nil {

				success = true
				time.Sleep(50 * time.Millisecond)
				fmt.Printf("[Worker %d] Send email to: %s\n", id, recipient.Email)
				metrics.EmailsSent.Inc()

				break
			}
			
			log.Printf("[Worker %d] Error sending to %s (attempt %d): %v",
                id, recipient.Email, attempt, err)

            time.Sleep(1 * time.Second)
			
		}

		if !success {
			metrics.EmailsFailed.Inc()
			dlq <- recipient // push to DLQ
		}

		metrics.WorkerActive.Dec()
		
	}
}