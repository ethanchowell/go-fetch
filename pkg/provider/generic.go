package provider

type Generic struct {
	Store
}

func (p Generic) Fetch(tag string, artifact string) error {
	return nil
}
