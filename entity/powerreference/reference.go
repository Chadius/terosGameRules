package powerreference

// Reference is used to identify a power and is used to quickly identify a power.
type Reference struct {
	Name    string `json:"name" yaml:"name"`
	PowerID string `json:"id" yaml:"id"`
}
