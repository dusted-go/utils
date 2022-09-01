package mailman

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"cloud.google.com/go/pubsub"

	"github.com/dusted-go/fault/fault"
)

// ------
// Private email struct:
// ------
// This struct is private on purpose so the user
// is forced to use the NewEmail function and the
// associated mutation functions which allow a
// fluent way to construct a correct email object.
// ------

type email struct {
	TraceID      string
	Domain       string
	Sender       string
	Recipients   []string
	CC           []string
	BCC          []string
	ReplyTo      string
	Subject      string
	Plaintext    string
	HTML         string
	TemplateName string
	TemplateData map[string]string
}

// NewEmail creates a new *email struct.
// nolint
func (c *Client) NewEmail(subject string, recipients ...string) *email {
	return &email{
		Domain:     c.domain,
		Sender:     c.sender,
		Recipients: recipients,
		Subject:    subject,
	}
}

func (e *email) SetTraceID(traceID string) *email {
	e.TraceID = traceID
	return e
}

func (e *email) SetCC(cc ...string) *email {
	e.CC = cc
	return e
}

func (e *email) SetBCC(bcc ...string) *email {
	e.BCC = bcc
	return e
}

func (e *email) SetReplyTo(replyTo string) *email {
	e.ReplyTo = replyTo
	return e
}

func (e *email) SetText(text string) *email {
	e.Plaintext = text
	return e
}

func (e *email) SetHTML(body string) *email {
	e.HTML = body
	return e
}

func (e *email) SetTemplate(templateName string, templateData map[string]string) *email {
	e.TemplateName = templateName
	e.TemplateData = templateData
	return e
}

func (e *email) String() string {
	return fmt.Sprintf(
		"email: { from: \"%v\", to: %v, subject: \"%v\" }",
		e.Sender,
		e.Recipients,
		e.Subject)
}

func (e *email) ToBinary() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(e)
	if err != nil {
		return nil, fault.SystemWrap(err, "Client", "ToBinary",
			"failed to encode email message")
	}
	return b.Bytes(), nil
}

// ------
// Public Client struct:
// ------
// The Client puts an email message into a pub sub
// topic inside Google Cloud which gets subsequently
// picked up by a private Google Cloud Function to
// actually send the message.
// This is a custom implementation for Dusted Codes projects.
// ------

const emptyMessageID = ""

// Client allows sending emails to a recipient.
type Client struct {
	topic           *pubsub.Topic
	domain          string
	sender          string
	environmentName string
}

// New creates a new Client client to send emails.
func New(topic *pubsub.Topic, domain, sender, environmentName string) *Client {
	return &Client{
		topic:           topic,
		domain:          domain,
		sender:          sender,
		environmentName: environmentName,
	}
}

func (c *Client) sendMessage(
	ctx context.Context,
	msg *email) (string, error) {

	if c.topic == nil || len(c.domain) == 0 || len(c.sender) == 0 {
		return emptyMessageID,
			fault.System("Client", "sendMessage",
				"cannot send email because the topic, domain or sender were not set")
	}

	data, err := msg.ToBinary()
	if err != nil {
		return emptyMessageID,
			fault.SystemWrap(err, "Client", "sendMessage",
				"failed to serialize message to byte array.")
	}
	attr := map[string]string{
		"environment": c.environmentName,
	}

	if len(msg.TraceID) > 0 {
		attr["traceID"] = msg.TraceID
	}

	result := c.topic.Publish(ctx,
		// No need to init the rest
		// nolint: exhaustivestruct
		&pubsub.Message{
			Data:       data,
			Attributes: attr,
		})

	// Wait until result is ready or request has been cancelled:
	select {
	case <-result.Ready():
		msgID, err := result.Get(ctx)
		if err != nil {
			return msgID,
				fault.SystemWrap(err, "Client", "sendMessage",
					"failed to publish message to PubSub topic")
		}
		return msgID, nil
	case <-ctx.Done():
		return emptyMessageID,
			fault.System("Client", "sendMessage",
				"context got cancelled before email status could get verified")
	}
}

// Send publishes a new email message to the `emails` topic which is being subscribed
// by the Mailman Google Cloud Function.
func (c *Client) Send(
	ctx context.Context,
	email *email) (msgID string, err error) {
	return c.sendMessage(ctx, email)
}
