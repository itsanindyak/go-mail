package dlqueue

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/itsanindyak/email-campaign/types"
)

// DlqWorker runs a dedicated worker for the Dead Letter Queue (DLQ).
// It receives failed recipients that have exhausted all retry attempts and logs them
// to a file (dlq.log) for later analysis or manual intervention.
// The worker signals completion via the provided WaitGroup when the channel is closed.
func DlqWorker(ctx context.Context, dlq chan types.Recipient, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.OpenFile("dlq.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[DLQ] failed to open file: %v", err)
		return
	}
	defer f.Close()

	for {
		select {
		case <-ctx.Done():
			return

		case r := <-dlq:
			fmt.Fprintf(f, "Failed: %s,%s\n", r.Name, r.Email)
		}
	}
}
