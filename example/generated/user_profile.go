package openapi

type UserProfile struct {
	Email string
	Name  string
}

func (instance *UserProfile) Validate() error {
	return nil
}
