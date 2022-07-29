package openapi

type Company struct {
	Name string
}

func (instance *Company) Validate() error {
	return nil
}
