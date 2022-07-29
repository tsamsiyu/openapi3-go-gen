package openapi

type Foo struct {
	Bar    string
	Baz    Baz
	King   FooKing
	Queens []FooQueen
}

func (instance *Foo) Validate() error {
	return nil
}
