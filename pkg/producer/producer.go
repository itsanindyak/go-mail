package producer

import (
	"encoding/csv"
	"io"
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

	_, err = reader.Read() // skip header
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	for{
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		// Ensure the record has at least name and email columns before accessing.
		if len(record) < 2 {
			// Skip malformed or incomplete rows.
			continue
		}
		ch <- types.Recipient{
			Name: record[0],
			Email: record[1],
			Attempts: 0,
		}
	}
	
	return nil

}