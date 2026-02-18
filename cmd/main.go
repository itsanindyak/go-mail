package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/itsanindyak/email-campaign/pkg/consumer"
	"github.com/itsanindyak/email-campaign/pkg/dlqueue"
	"github.com/itsanindyak/email-campaign/pkg/producer"
	"github.com/itsanindyak/email-campaign/pkg/telemetry"
	"github.com/itsanindyak/email-campaign/types"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func main() {
	startTime := time.Now()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	oTELctx := context.Background()
	shutdown := telemetry.InitOTel(oTELctx)
	defer shutdown(oTELctx)
	log.Println("Observability initialized. Starting engine...")

	emailsPerSecondStr := os.Getenv("EMAILS_PER_SEC")

	emailsPerSecond, err := strconv.ParseFloat(emailsPerSecondStr, 64)
	if err != nil {
		log.Fatalf("Invalid EMAILS_PER_SEC value: %v", err)
	}

	limiter := rate.NewLimiter(rate.Limit(emailsPerSecond), 1)

	path := "email.csv"
	workerCount := 3

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	var wg sync.WaitGroup
	var dlqWg sync.WaitGroup
	var itemWg sync.WaitGroup
	var producerWg sync.WaitGroup

	recipientChannel := make(chan types.Recipient, 100)
	dlqChannel := make(chan types.Recipient, 100)

	// log.Flags()

	dlqWg.Add(1)
	go dlqueue.DlqWorker(ctx, dlqChannel, &dlqWg)

	for i := range workerCount {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			consumer.EmailWorker(ctx, i, recipientChannel, dlqChannel, &itemWg, limiter)
		}(i)
	}

	log.Println("Producer started...")
	producerWg.Add(1)
	go func() {
		defer producerWg.Done()

		err = producer.LoadFile(path, recipientChannel, &itemWg)
		if err != nil {
			log.Printf("Producer error: %v", err)
		}
	}()

	producerWg.Wait()

	log.Println("Producer finished sending all items.")

	itemWg.Wait()
	log.Println("All items processed.")

	cancel()

	wg.Wait()
	close(recipientChannel)

	close(dlqChannel)
	dlqWg.Wait()

	elapsed := time.Since(startTime)
	log.Printf("Email processing completed in %s", elapsed)
}
