package mail

import (
	"fmt"
	"net/smtp"

	"github.com/itsanindyak/email-campaign/types"
)

func MailSend(recipent types.Recipient) error {
	url := "localhost:1025"

	formattedMsg := fmt.Sprintf("To: %s\r\nSubject: Test Email\r\n\r\n%s\r\n",recipent.Email,"Just testing our email campaign\r\nname: "+recipent.Name)

	msg := []byte(formattedMsg)

	err := smtp.SendMail(url,nil,"iankoley04@gmail.com",[]string{recipent.Email},msg)

	if err !=nil{
		return  err
	}
	return nil
}