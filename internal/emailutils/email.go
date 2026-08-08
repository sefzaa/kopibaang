package emailutils

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/smtp"
)

// GenerateOTP membuat 6 digit angka random
func GenerateOTP() string {
	b := make([]byte, 6)
	_, _ = io.ReadAtLeast(rand.Reader, b, 6)
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// SendOTPResetPassword sekarang menerima kredensial SMTP
func SendOTPResetPassword(toEmail string, otp string, smtpEmail string, smtpPassword string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: OTP Reset Password - Kopibang Coffee\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<h2>Reset Password</h2><p>Kode OTP Anda adalah: <b>%s</b></p><p>Kode ini berlaku selama 5 menit. Jangan berikan kode ini kepada siapapun.</p>", otp)
	
	message := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpEmail, []string{toEmail}, message)
	return err
}