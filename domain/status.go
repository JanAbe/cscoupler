package domain

// Status type for conveying the status
// students / project's of companies
// can have.
type Status uint8

const (
	// Available indicates an entity is available
	Available Status = iota

	// Unavailable indicates an entity is unavailable
	Unavailable
)
