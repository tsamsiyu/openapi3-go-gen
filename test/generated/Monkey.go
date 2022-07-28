package openapi

import (
	"errors"
)

type Monkey struct {
	Age int
}

func (instance *Monkey) Validate() error {
	if instance.Age > 20 {
		return errors.New("Field Age should not be greater than 20")
	}
	if instance.Age <= 3 {
		return errors.New("Field Age should not be less or equal than 3")
	}
	return nil
}
