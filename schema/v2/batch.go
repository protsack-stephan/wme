package schema

import "time"

type Batch struct {
	Name         string     `json:"name,omitempty"`
	Identifier   string     `json:"identifier,omitempty"`
	Version      string     `json:"version,omitempty"`
	DateModified *time.Time `json:"date_modified,omitempty"`
	IsPartOf     *Project   `json:"is_part_of,omitempty"`
	InLanguage   *Language  `json:"in_language,omitempty"`
	Namespace    *Namespace `json:"namespace,omitempty"`
	Size         *Size      `json:"size,omitempty"`
}
