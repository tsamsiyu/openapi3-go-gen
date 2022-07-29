package openapi

type Merchant struct {
	Name string
}

func (instance *Merchant) Validate() error {
	return nil
}
