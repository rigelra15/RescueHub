package middlewares

import (
	"crypto/rand"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func GenerateOTP() string {
	n := make([]byte, 3)
	_, err := rand.Read(n)
	if err != nil {
		panic(err)
	}

	otp := fmt.Sprintf("%06d", int(n[0])%10*100000 + int(n[1])%1000 + int(n[2])%100)
	return otp
}

func SendOTPToEmail(toEmail, otp string) error {
	e := email.NewEmail()
	e.From = "RescueHub <your-email@example.com>"
	e.To = []string{toEmail}
	e.Subject = "Kode OTP Anda"
	e.Text = []byte("Kode OTP Anda adalah: " + otp + ". Berlaku selama 5 menit.")

	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "bemysample.id@gmail.com", "uhse oyou lalj konk", "smtp.gmail.com"))
	return err
}
