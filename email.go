package email

//go:generate moq -out email_mock_test.go . Sender SESApi

import (
	"fmt"
	"regexp"
)

// Sender sends an email message
type Sender interface {
	Send(*Message) error
}

// Message collects everything that's needed to send a message
type Message struct {
	sender   Sender
	From     string
	To       []string
	Cc       []string
	Bcc      []string
	Subject  string
	TextBody string
	HtmlBody string
}

// Send sends an email message
func (m *Message) Send() error {
	err := validateMessage(m)
	if err != nil {
		return fmt.Errorf("message validation failed: %s", err.Error())
	}
	return m.sender.Send(m)
}

func validateMessage(m *Message) error {
	if m.sender == nil {
		return fmt.Errorf("sender is not defined")
	}
	if m.From == "" {
		return fmt.Errorf("from email is required")
	}
	if !emailValid(m.From) {
		return fmt.Errorf("invalid from email")
	}
	if len(m.To) < 1 {
		return fmt.Errorf("at least one to email is required")
	}
	for _, t := range m.To {
		if !emailValid(t) {
			return fmt.Errorf("invalid to email")
		}
	}
	return nil
}

// New returns a new message using default (AWS SES) sender
func New(opts ...func(*Message)) *Message {
	ses := newSES()
	sesSender := NewSESSender(ses)
	return NewWithSender(sesSender, opts...)
}

// NewWithSender returns a new message using custom sender
func NewWithSender(s Sender, opts ...func(*Message)) *Message {
	m := &Message{sender: s}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func From(from string) func(*Message) {
	return func(m *Message) {
		m.From = from
	}
}

func To(to ...string) func(*Message) {
	return func(m *Message) {
		m.To = to
	}
}

func Cc(cc ...string) func(*Message) {
	return func(m *Message) {
		m.Cc = cc
	}
}

func Bcc(bcc ...string) func(*Message) {
	return func(m *Message) {
		m.Bcc = bcc
	}
}

func Subject(subject string) func(*Message) {
	return func(m *Message) {
		m.Subject = subject
	}
}

func TextBody(tBody string) func(*Message) {
	return func(m *Message) {
		m.TextBody = tBody
	}
}

func HtmlBody(hBody string) func(*Message) {
	return func(m *Message) {
		m.HtmlBody = hBody
	}
}

func emailValid(e string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(e)
}
