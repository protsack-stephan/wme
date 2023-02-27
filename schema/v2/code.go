package schema

type Code struct {
	Identifier  string `json:"identifier,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
