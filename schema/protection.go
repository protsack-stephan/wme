package schema

// Protection level of protection for the page.
type Protection struct {
	Type   string `json:"type,omitempty"`
	Level  string `json:"level,omitempty"`
	Expiry string `json:"expiry,omitempty"`
}
