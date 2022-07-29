package openapi

type FooKing struct {
	Years int
}

func (instance *FooKing) Validate() error {
	return nil
}
