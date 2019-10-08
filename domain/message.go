package domain

// Message struct conveying messages
// users can send to each other,
// optionally linked to a certain project
type Message struct {
	sender    string // the sender of the message (email / id of user)
	receiver  string // the receiver of the message (email / id of user)
	body      string // the message body
	projectID string // the project, the reason why the sender sent the message (optional)
}

// NewMessage creates a new message based on
// the provided input if all is valid, returning
// an error otherwise
func NewMessage(sender, receiver, body string) (Message, error) {
	// what checks does this need?
	// i'm thinking about a check to see if the provided
	// sender and receiver emails exist in the system
	// but that would result in a dependency on a repo, hmmm

	return Message{
		sender:    sender,
		receiver:  receiver,
		body:      body,
		projectID: "",
	}, nil
}
