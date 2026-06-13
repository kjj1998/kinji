package app

// ClientError is an application error caused by bad client input. The HTTP
// adapter maps it to a 4xx response, surfacing Reason to the caller.
type ClientError struct {
	Reason string
}

func (e *ClientError) Error() string { return e.Reason }
