package openapi

type Foo struct {
	King   FooKing
	Queens []FooQueen
	Bar    string
	Baz    Baz
}

func (instance *Foo) Validate() error {
	return nil
}
