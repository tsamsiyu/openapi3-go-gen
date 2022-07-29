package openapi

import (
	"errors"
	"regexp"
)

type Animal struct {
	Unknown  interface{}
	Unknowns []interface{}
	Meow     string
	Bark     string
	Age      int
}

func (instance *Animal) Validate() error {
	if instance.Unknowns == nil {
		return errors.New("Value for field Unknowns must be present")
	}
	if len(instance.Unknowns) > 100 {
		return errors.New("Number of elements of Unknowns should not exceed 100")
	}
	if len(instance.Unknowns) < 5 {
		return errors.New("Number of elements of Unknowns should not be less than 5")
	}
	if len(instance.Meow) > 255 {
		return errors.New("Field Meow size should not be greater than 255")
	}
	if len(instance.Meow) < 3 {
		return errors.New("Field Meow size should not be less than 3")
	}
	if match, _ := regexp.MatchString(`^\d{3}-\d{2}-\d{4}$`, instance.Meow); !match {
		return errors.New("Field Meow is not formatted correctly")
	}
	containsBark := false
	enum := []string{"rark", "bark", "kararak", "howk"}
	for _, v := range enum {
		if v == instance.Bark {
			containsBark = true
			break
		}
	}

	if !containsBark {
		return errors.New("Value for field Bark is not allowed")
	}
	if instance.Age > 20 {
		return errors.New("Field Age should not be greater than 20")
	}
	if instance.Age <= 3 {
		return errors.New("Field Age should not be less or equal than 3")
	}
	return nil
}
