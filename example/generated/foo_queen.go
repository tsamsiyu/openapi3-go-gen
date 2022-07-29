package openapi

type FooQueen struct {
	Level int
}

func (instance *FooQueen) Validate() error {
	return nil
}
