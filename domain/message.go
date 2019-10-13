package domain

import (
	"errors"
	"strings"
	"time"
)

// Message struct conveying messages
// users can send to each other,
// optionally linked to a certain project
type Message struct {
	Sender    string // the sender of the message (id of user)
	Receiver  string // the receiver of the message (id of user)
	Body      string // the message body
	ProjectID string // the project, the possible subject of the message (optional)
	CreatedAt time.Time
}

// NewMessage creates a new message based on
// the provided input if all is valid, returning
// an error otherwise
func NewMessage(sender, receiver, body string) (Message, error) {
	// what checks does this need?
	// i'm thinking about a check to see if the provided
	// sender and receiver emails exist in the system
	// but that would result in a dependency on a repo, hmmm
	if len(strings.TrimSpace(sender)) == 0 {
		return Message{}, errors.New("provided sender can't be empty")
	}

	if len(strings.TrimSpace(receiver)) == 0 {
		return Message{}, errors.New("provided receiver can't be empty")
	}

	if len(strings.TrimSpace(body)) == 0 {
		return Message{}, errors.New("provided body can't be empty")
	}

	return Message{
		Sender:    sender,
		Receiver:  receiver,
		Body:      body,
		ProjectID: "",
		CreatedAt: time.Now(),
	}, nil
}
