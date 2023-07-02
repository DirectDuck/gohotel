package db

import "fmt"

type ValidationError struct {
	Fields map[string]string
}

func (self ValidationError) Error() string {
	errStr := "ValidationError:"
	for k, v := range self.Fields {
		errStr += fmt.Sprintf("\n%s: %s", k, v)
	}
	return errStr
}
