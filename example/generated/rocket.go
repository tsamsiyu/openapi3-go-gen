package openapi

type Rocket struct {
	Speed float64
}

func (instance *Rocket) Validate() error {
	return nil
}
