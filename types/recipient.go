package types

import "github.com/itsanindyak/email-campaign/mail"

// Recipient represents an email recipient loaded from the CSV file.
// It contains the recipient's name, email address, and a counter for retry attempts.
type Recipient struct {
	Name     string // Full name of the recipient
	Email    string // Email address to send to
	Attempts int    // Number of send attempts made (used for retry logic)
}

type EmailJob struct {
    Recipient Recipient
    Template  mail.Template
}