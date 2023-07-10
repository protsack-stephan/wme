package schema

import "time"

// Thing schema for Machine Readable entity.
// Tries to compliant with https://schema.org/Thing.
type Thing struct {
	// Name is the name of the thing (article).
	Name string `json:"name,omitempty"`

	// Identifier is a unique identifier for the thing (article, in scope of a single project).
	Identifier int `json:"identifier,omitempty"`

	// Abstract is a summary of the thing (article).
	Abstract string `json:"abstract,omitempty"`

	// Version is the metadata about the version of the thing (article).
	Version *Version `json:"version,omitempty"`

	// URL is the URL of the thing (article).
	URL string `json:"url,omitempty"`

	// DateCreated is the date and time the thing (article) was created.
	DateCreated *time.Time `json:"date_created,omitempty"`

	// DateModified is the date and time the thing (article) was last modified.
	DateModified *time.Time `json:"date_modified,omitempty"`

	// MainEntity is the main (Wikidata) entity of the thing (article).
	MainEntity *Entity `json:"main_entity,omitempty"`

	// AdditionalEntities are the additional (Wikidata) entities used in the thing (article).
	AdditionalEntities []*Entity `json:"additional_entities,omitempty"`

	// IsPartOf is the project that the thing (article) belongs to.
	IsPartOf *Project `json:"is_part_of,omitempty"`

	// InLanguage is the language of the thing (article).
	InLanguage *Language `json:"in_language,omitempty"`

	// Image specifies the image related to the thing (article).
	Image *Image `json:"image,omitempty"`
}
