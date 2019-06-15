package email

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SESApi interface is a subset of AWS SES Api that this packge utilizes
// AWS SES instance will satisfy this interface
// But it can be also easily mocked
type SESApi interface {
	SendEmail(*ses.SendEmailInput) (*ses.SendEmailOutput, error)
}

// sesSender implements Sender interface for AWS SES
type sesSender struct {
	SESApi
}

// Send sends an email message
func (s *sesSender) Send(m *Message) error {
	i := messagetoSESInput(m)
	_, err := s.SendEmail(i)
	if err != nil {
		return err
	}
	return nil
}

// NewsesSender returns a new instance of sesSender
// which implements Sender interface
func NewSESSender(s SESApi) Sender {
	return &sesSender{s}
}

func newSES() SESApi {
	sess, _ := session.NewSession()
	return ses.New(sess)
}

func messagetoSESInput(m *Message) *ses.SendEmailInput {
	charset := "UTF-8"
	i := &ses.SendEmailInput{
		Source:      &m.From,
		Destination: &ses.Destination{ToAddresses: slicetoPtrSlice(m.To)},
		Message: &ses.Message{
			Body: &ses.Body{},
		},
	}
	if len(m.Cc) > 0 {
		i.Destination.CcAddresses = slicetoPtrSlice(m.Cc)
	}
	if len(m.Bcc) > 0 {
		i.Destination.BccAddresses = slicetoPtrSlice(m.Bcc)
	}
	if m.Subject != "" {
		i.Message.Subject = &ses.Content{
			Charset: &charset,
			Data:    &m.Subject,
		}
	}
	if m.TextBody != "" {
		i.Message.Body.Text = &ses.Content{
			Charset: &charset,
			Data:    &m.TextBody,
		}
	}
	if m.HtmlBody != "" {
		i.Message.Body.Html = &ses.Content{
			Charset: &charset,
			Data:    &m.HtmlBody,
		}
	}
	return i
}

func slicetoPtrSlice(s []string) []*string {
	result := make([]*string, len(s))
	for i := range s {
		result[i] = &s[i]
	}
	return result
}
