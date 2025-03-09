package registration

import (
	"crypto/tls"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/smtp"
)

// SendVerificationEmail
//
// Read the email info env file
// Create the message
// Send the email to the specified recipient with a randomized token
func SendVerificationEmail(recipientEmail, token string) error {
	// Read the email info env file
	emailInfoMap, err := godotenv.Read("registration/emailinfo.env")

	if err != nil {
		log.Println(err)
		return err
	}

	var EMAIL_PUBLIC = emailInfoMap["EMAIL_PUBLIC"]
	var EMAIL_PRIVATE = emailInfoMap["EMAIL_PRIVATE"]
	var APP_PASSWORD = emailInfoMap["APP_PASSWORD"]

	// Create the message
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Verify your email\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
		"\r\n"+
		"Please copy and paste the following code into the email verification box on the previous page. It will be valid for 5 minutes.\r\n\r\n"+
		"Please do not share this code with anyone.\r\n\r\n"+
		"Verification Code: %s\r\n\r\n",
		EMAIL_PUBLIC, recipientEmail, token)

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpServer := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	smtpAuth := smtp.PlainAuth("", EMAIL_PRIVATE, APP_PASSWORD, smtpHost)

	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	client, err := smtp.Dial(smtpServer)

	if err != nil {
		log.Println(err)
		return err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(tlsConfig); err != nil {
			log.Println(err)
			return err
		}
	} else {
		log.Println("WARNING: Server does not support STARTTLS")
	}

	// Authenticate
	if err = client.Auth(smtpAuth); err != nil {
		log.Println(err)
		return err
	}

	// Add Sender header
	if err = client.Mail(EMAIL_PUBLIC); err != nil {
		log.Println(err)
		return err
	}

	// Add Recipient header
	if err = client.Rcpt(recipientEmail); err != nil {
		log.Println(err)
		return err
	}

	// Create the io.WriteCloser
	wc, err := client.Data()

	if err != nil {
		log.Println(err)
		return err
	}

	// Write the message
	if _, err = wc.Write([]byte(message)); err != nil {
		log.Println(err)
		return err
	}

	// Close the SMTP Writer
	if err = wc.Close(); err != nil {
		log.Println(err)
		return err
	}

	// Close the SMTP connection
	if err = client.Quit(); err != nil {
		log.Println(err)
		return err
	}

	log.Println("Verification email sent successfully to:", recipientEmail)

	return nil
}
