package mail

import (
	"fmt"
)

func FormatMIME(recipient, subject, htmlBody string) string {
	return fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s\r\n",
		recipient,
		subject,
		htmlBody,
	)
}
