package openapi

type CreateUser struct {
	Merchant Merchant
	Photos   string
	Profile  UserProfile
	Company  Company
	Id       string
}

func (instance *CreateUser) Validate() error {
	return nil
}
