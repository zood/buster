package email

type SendEmailer interface {
	SendEmail(from string, to string, subj string, textMsg string, htmlMsg *string) error
}

type mockSendEmailer struct{}

func (mse mockSendEmailer) SendEmail(from string, to string, subj string, textMsg string, htmlMsg *string) error {
	return nil
}

func NewMock() SendEmailer {
	return mockSendEmailer{}
}
