package openapi

type Car struct {
	Model string
	Year  int
}

func (instance *Car) Validate() error {
	return nil
}
