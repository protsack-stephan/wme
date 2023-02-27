package schema

// Filter payload for filters in realtime API.
type Filter struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
