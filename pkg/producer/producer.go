package producer

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sync"

	"github.com/itsanindyak/email-campaign/types"
)

// LoadFile reads a CSV file containing recipient data and sends each recipient
// to the provided channel for processing. The CSV file should have at least two
// columns: Name and Email. The first row is treated as a header and skipped.
//
// Note: The caller is responsible for closing the channel after LoadFile returns.
func LoadFile(path string, ch chan types.Recipient,wg *sync.WaitGroup ) error {

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read() // skip header
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Printf("Error reading CSV header: %v", err)
		return err
	}

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error reading CSV record: %v", err)
			break
		}

		// Ensure the record has at least name and email columns before accessing.
		if len(record) < 2 {
			log.Printf("Skipping malformed row: %v", record)
			continue
		}

		wg.Add(1)
		
		ch <- types.Recipient{
			Name:     record[0],
			Email:    record[1],
			Attempts: 0,
		}
	}

	return nil
}
