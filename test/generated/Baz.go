package openapi

type Baz struct {
	Lol string
}

func (instance *Baz) Validate() error {
	return nil
}
