package openapi

type CreateUser struct {
	Id       string
	Merchant Merchant
	Photos   string
	Profile  UserProfile
	Company  Company
}

func (instance *CreateUser) Validate() error {
	return nil
}
