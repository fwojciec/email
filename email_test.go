package email_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/fwojciec/email"
	"github.com/matryer/is"
)

var (
	from    string = "from@email.com"
	to1     string = "to1@email.com"
	to2     string = "to2@email.com"
	cc1     string = "cc1@email.com"
	cc2     string = "cc2@email.com"
	bcc1    string = "bcc1@email.com"
	bcc2    string = "bcc2@email.com"
	subject string = "subject"
	tBody   string = "text body"
	hBody   string = "html body"
	charSet string = "UTF-8"
)

func TestNew(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	m := email.New()
	is.True(m != nil) // should have returned a non nil message
}

func TestSendValidationErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		opts  []func(*email.Message)
		err   error
		calls int
	}{
		{
			name:  "minimum requirements met",
			opts:  []func(*email.Message){email.From(from), email.To(to1)},
			err:   nil,
			calls: 1,
		},
		{
			name:  "no from email",
			opts:  []func(*email.Message){email.To(to1)},
			err:   fmt.Errorf("message validation failed: from email is required"),
			calls: 0,
		},
		{
			name:  "no to emails",
			opts:  []func(*email.Message){email.From(from)},
			err:   fmt.Errorf("message validation failed: at least one to email is required"),
			calls: 0,
		},
		{
			name:  "from email invalid",
			opts:  []func(*email.Message){email.From("wrong"), email.To(to1)},
			err:   fmt.Errorf("message validation failed: invalid from email"),
			calls: 0,
		},
		{
			name:  "to email invalid",
			opts:  []func(*email.Message){email.From(from), email.To("wrong")},
			err:   fmt.Errorf("message validation failed: invalid to email"),
			calls: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			sesMock := &email.SESApiMock{
				SendEmailFunc: func(i *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
					return nil, nil
				},
			}
			s := email.NewSESSender(sesMock)
			m := email.NewWithSender(s, test.opts...)
			err := m.Send()
			is.Equal(err, test.err)                             // excpected a different error value
			is.Equal(len(sesMock.SendEmailCalls()), test.calls) // expected SendEmail to not have been called
		})
	}
}

func TestMessageSend(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		opts []func(*email.Message)
		exp  *ses.SendEmailInput
	}{
		{
			name: "minimal",
			opts: []func(*email.Message){email.From(from), email.To(to1, to2)},
			exp: &ses.SendEmailInput{
				Source: &from,
				Destination: &ses.Destination{
					ToAddresses: []*string{&to1, &to2},
				},
				Message: &ses.Message{
					Body: &ses.Body{},
				},
			},
		},
		{
			name: "everything",
			opts: []func(*email.Message){
				email.From(from),
				email.To(to1, to2),
				email.Cc(cc1, cc2),
				email.Bcc(bcc1, bcc2),
				email.Subject(subject),
				email.TextBody(tBody),
				email.HtmlBody(hBody),
			},
			exp: &ses.SendEmailInput{
				Source: &from,
				Destination: &ses.Destination{
					ToAddresses:  []*string{&to1, &to2},
					CcAddresses:  []*string{&cc1, &cc2},
					BccAddresses: []*string{&bcc1, &bcc2},
				},
				Message: &ses.Message{
					Subject: &ses.Content{
						Data:    &subject,
						Charset: &charSet,
					},
					Body: &ses.Body{
						Text: &ses.Content{
							Data:    &tBody,
							Charset: &charSet,
						},
						Html: &ses.Content{
							Data:    &hBody,
							Charset: &charSet,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			var pInput *ses.SendEmailInput
			s := email.NewSESSender(&email.SESApiMock{
				SendEmailFunc: func(i *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
					pInput = i
					return nil, nil
				},
			})
			m := email.NewWithSender(s, test.opts...)
			err := m.Send()
			is.NoErr(err)                                // // expected no error
			is.True(reflect.DeepEqual(pInput, test.exp)) // expected correctly formed input
		})
	}
}

func TestSenderNotDefined(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	m := email.Message{}
	err := m.Send()
	is.True(err != nil)                                                       // expected an error
	is.Equal(err.Error(), "message validation failed: sender is not defined") // expcted different error message
}
