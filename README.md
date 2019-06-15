# email

A simple wrapper around cloud email services, to be used in lambda functions and other apps.

## Motivation

I find myself rewriting the same AWS SES code whenever I need to send an email from a lambda function or an app. This package makes it easier.

## Prerequisites

This package assumes that AWS SES service is correctly configured and that the environment variables needed to properly configure an AWS session are set. More about AWS SES requirements here: https://docs.aws.amazon.com/sdk-for-go/api/service/ses/#SES.SendEmail

## How to use

This example shows the complete API for the default sender:

```go
err := email.New(
    email.From("from@email.com"), // required
    email.To("to1@email.com", "to2@email.com"), // at least one address required
    email.Cc("cc1@email.com". "cc2@email.com"),
    email.Bcc("bcc1@email.com", "bcc2@email.com"),
    email.Subject("Example Email"),
    email.TextBody("Text version of the body."),
    email.HtmlBody("<p>Html version of the body.</p>"),
).Send()
```

Alternatively it is also possible to set the fields of the `Message` struct by hand, but the struct itself should still be initialized using the `New` initializer, since it attaches a Sender to the struct.

```go
msg := email.New()
msg.From = "from@email.com"
msg.To = []string{"to1@email.com", "to2@email.com"}
msg.Cc = []string{"cc1@email.com". "cc2@email.com"}
msg.Bcc = []string{"bcc1@email.com". "bcc2@email.com"}
msg.Subject = "Example Email"
msg.TextBody = "Text version of the body."
msg.HtmlBody = "<p>Html version of the body.</p>"
err := msg.Send()
```

## Custom Senders

This pacakge only implements AWS SES (Simple Email Service) as a Sender, since this is what I use, but it can be easily extended to use different Senders (implementations of Sender interface):

```go
type customSender struct {}
func (s *customSender) Send(m *email.Message) error {
    // implementation
}
err := email.NewWithSender(
    customSender,
    email.From("from@email.com"),
    email.To("to@email.com"),
    email.Subject("Example Email"),
    email.TextBody("Sent using a custom sender!"),
).Send()
```
