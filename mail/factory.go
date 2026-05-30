package mail

import (
	"fmt"
	"os"
)

func NewMailService() (MailService, error) {

	switch vendor := os.Getenv("SMTP_VENDOR"); vendor {

	case "resend":
		apiKey := os.Getenv("RESEND_API_KEY")

		if apiKey == "" {
			return nil, fmt.Errorf("RESEND_API_KEY missing")
		}

		fmt.Printf("SMTP_VENDOR=%q\n", vendor)
    fmt.Printf("RESEND_API_KEY=%q\n", apiKey)

		return NewResendService(apiKey), nil

	case "smtp":
		return NewSMTPService(), nil

	default:
		return nil, fmt.Errorf("unknown mail provider")
	}
}