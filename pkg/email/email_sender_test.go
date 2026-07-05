package email

import "testing"

func TestLogEmailSenderSendOTPReturnsNoError(t *testing.T) {
	sender := NewLogEmailSender()
	if err := sender.SendOTP("user@example.com", "123456"); err != nil {
		t.Fatalf("預期無錯誤，實際為: %v", err)
	}
}
