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

    f, err := os.OpenFile("dlq.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("[DLQ] failed to open file: %v", err)
        return
    }
    defer f.Close()

    for r := range dlq {
        fmt.Fprintf(f, "Failed: %s,%s\n", r.Name, r.Email)
    }
}
