package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/itsanindyak/email-campaign/pkg/consumer"
	"github.com/itsanindyak/email-campaign/pkg/dlqueue"
	"github.com/itsanindyak/email-campaign/pkg/metrics"
	"github.com/itsanindyak/email-campaign/pkg/producer"
	"github.com/itsanindyak/email-campaign/types"
)



func main() {
	err := godotenv.Load() 
    if err != nil {
    	log.Fatal(err)
    }
	path := "email.csv"
	workerCount := 3

	metrics.Init()

	go metrics.StartMetrics(":2112")

	var wg sync.WaitGroup
	var dlqWg sync.WaitGroup

	recipientChannel := make(chan types.Recipient)
	dlqChannel := make(chan types.Recipient)

	go producer.LoadFile(path, recipientChannel)

	dlqWg.Add(1)
	go dlqueue.DlqWorker(dlqChannel,&dlqWg)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go consumer.EmailWorker(i, recipientChannel, dlqChannel, &wg)
	}

	wg.Wait()

	close(dlqChannel)

	dlqWg.Wait()

	log.Println("Email processing completed. Metrics server running...")
	select {} // stop main from exiting
	
}