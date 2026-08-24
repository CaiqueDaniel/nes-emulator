package application

import "fmt"

var _ error = NotFoundError("")

type NotFoundError string

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Not found: %s", string(e))
}

var _ error = InfrastructureError("")

type InfrastructureError string

func (e InfrastructureError) Error() string {
	return fmt.Sprintf("Infrastructure error: %s", string(e))
}
