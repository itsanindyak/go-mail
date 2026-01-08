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

	reader.Read() // skip header

	for{
		record, err := reader.Read()

		if err == io.EOF {
        	break
    	}

		if err != nil {
			break
		}

		ch <- types.Recipient{
			Name: record[0],
			Email: record[1],
			Attempts: 0,
		}
	}
	
	return nil

}