package dlqueue

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/itsanindyak/email-campaign/types"
)

func DlqWorker(dlq chan types.Recipient, wg *sync.WaitGroup) {
    defer wg.Done()

    for r := range dlq {
        log.Printf("[DLQ] Permanently failed: %s (%s)", r.Email, r.Name)

        // Optional: write to DB or file
        f, _ := os.OpenFile("dlq.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        fmt.Fprintf(f, "Failed: %s,%s\n", r.Name, r.Email)
        f.Close()
    }
}
