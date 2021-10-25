package gauth

import "fmt"

type APIError struct {
	Err string `json:"error"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v", e.Err)
}
