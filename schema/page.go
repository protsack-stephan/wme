package schema

import "time"

// ArticleBody content of the page.
type ArticleBody struct {
	HTML     string `json:"html"`
	Wikitext string `json:"wikitext"`
}

// Page schema to represent full page in the stream.
type Page struct {
	Name               string        `json:"name"`
	Identifier         int           `json:"identifier,omitempty"`
	DateModified       *time.Time    `json:"date_modified,omitempty"`
	Protection         []*Protection `json:"protection,omitempty"`
	Version            *Version      `json:"version,omitempty"`
	URL                string        `json:"url,omitempty"`
	Namespace          *Namespace    `json:"namespace,omitempty"`
	InLanguage         *Language     `json:"in_language,omitempty"`
	MainEntity         *Entity       `json:"main_entity,omitempty"`
	AdditionalEntities []*Entity     `json:"additional_entities,omitempty"`
	Categories         []*Page       `json:"categories,omitempty"`
	Templates          []*Page       `json:"templates,omitempty"`
	Redirects          []*Page       `json:"redirects,omitempty"`
	IsPartOf           *Project      `json:"is_part_of,omitempty"`
	ArticleBody        *ArticleBody  `json:"article_body,omitempty"`
	License            []*License    `json:"license,omitempty"`
	Visibility         *Visibility   `json:"visibility,omitempty"`
}
