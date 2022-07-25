package schema

// License schema to represent license object.
type License struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}
