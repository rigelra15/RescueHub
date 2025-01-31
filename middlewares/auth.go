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

	otp := fmt.Sprintf("%06d", int(n[0])%10*100000+int(n[1])%1000+int(n[2])%100)
	return otp
}

func SendOTPToEmail(toEmail, otp string) error {
	e := email.NewEmail()
	e.From = "RescueHub <bemysample.id@gmail.com>"
	e.To = []string{toEmail}
	e.Subject = "Kode OTP Anda"

	e.HTML = []byte(fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Kode OTP Anda</title>
			<style>
				.container {
					font-family: Arial, sans-serif;
					max-width: 500px;
					margin: auto;
					padding: 20px;
					border: 1px solid #ddd;
					border-radius: 5px;
					box-shadow: 0 2px 4px rgba(0,0,0,0.1);
				}
				.otp {
					font-size: 24px;
					font-weight: bold;
					color: #d9534f;
				}
				.footer {
					font-size: 12px;
					color: #777;
					margin-top: 10px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Verifikasi OTP</h2>
				<p>Kode OTP Anda adalah:</p>
				<p class="otp">%s</p>
				<p>Silakan gunakan kode ini dalam waktu <strong>5 menit</strong>.</p>
				<p class="footer">Jika Anda tidak meminta kode ini, abaikan email ini.</p>
			</div>
		</body>
		</html>
	`, otp))

	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "bemysample.id@gmail.com", "uhse oyou lalj konk", "smtp.gmail.com"))
	return err

	// Get the password from Google App Passwords
}
