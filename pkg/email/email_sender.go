package email

import "log"

// EmailSender 定義寄送 Email 的介面，未來可替換為真實 SMTP/第三方 Provider 實作
type EmailSender interface {
	SendOTP(toEmail string, otp string) error
}

// LogEmailSender 本輪先以 log 模擬寄送，不接真實 Provider
type LogEmailSender struct{}

func NewLogEmailSender() EmailSender {
	return &LogEmailSender{}
}

func (s *LogEmailSender) SendOTP(toEmail string, otp string) error {
	log.Printf("[MOCK EMAIL] 寄送密碼重設 OTP 至 %s: %s", toEmail, otp)
	return nil
}
