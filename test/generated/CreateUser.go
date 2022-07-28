package openapi

type CreateUser struct {
	Company  Company
	Id       string
	Merchant Merchant
	Photos   string
	Profile  UserProfile
}

func (instance *CreateUser) Validate() error {
	return nil
}
