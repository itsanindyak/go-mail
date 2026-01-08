package producer

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/itsanindyak/email-campaign/types"
)

// Recipient represents a CSV recipient with name and email.


func LoadFile(path string, ch chan types.Recipient) error {
	defer close(ch)

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}

	for _, record := range records[1:] {
		// fmt.Println(record)
		ch <- types.Recipient{
			Name: record[0],
			Email: record[1],
			Attempts: 0,
		}
	}

	return nil

}