package users

import (
	"fmt"
	"net/smtp"
)

func mapUser(um *userModel, ur userReq) {
	um.Email = ur.Email
	um.Region = ur.Region
	um.Street = ur.Street
	um.Name = ur.Name
	um.Password = ur.Password
}

func sendEmail(to []string, sender, msg string) error {
	// Set up authentication information.
	identity := ""
	username := ""
	password := ""
	hostname := "smtp.gmail.com"
	port := 587
	auth := smtp.PlainAuth(identity, username, password, hostname)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", hostname, port), auth, sender, to, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
