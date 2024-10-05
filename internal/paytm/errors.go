package paytm

import "fmt"

// PaytmError represents an error response from Paytm Money API.
type PaytmError struct {
	Code    int
	Message string
}

func (e *PaytmError) Error() string {
	return fmt.Sprintf("PaytmError - Code: %d, Message: %s", e.Code, e.Message)
}
