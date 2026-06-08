package service

// ClientError is a service error caused by bad client input → maps to HTTP 4xx.
type ClientError struct {
	Reason string
}

func (e *ClientError) Error() string { return e.Reason }
