package schema

import (
	"time"
)

// Article schema for wikipedia article.
// Tries to compliant with https://schema.org/Article.
type Article struct {
	Name                   string           `json:"name,omitempty"`
	Identifier             int              `json:"identifier,omitempty"`
	DateCreated            *time.Time       `json:"date_created,omitempty"`
	DateModified           *time.Time       `json:"date_modified,omitempty"`
	DatePreviouslyModified *time.Time       `json:"date_previously_modified,omitempty"`
	Protection             []*Protection    `json:"protection,omitempty"`
	Version                *Version         `json:"version,omitempty"`
	PreviousVersion        *PreviousVersion `json:"previous_version,omitempty"`
	VersionIdentifier      string           `json:"-"`
	URL                    string           `json:"url,omitempty"`
	WatchersCount          int              `json:"watchers_count,omitempty"`
	Namespace              *Namespace       `json:"namespace,omitempty"`
	InLanguage             *Language        `json:"in_language,omitempty"`
	MainEntity             *Entity          `json:"main_entity,omitempty"`
	AdditionalEntities     []*Entity        `json:"additional_entities,omitempty"`
	Categories             []*Category      `json:"categories,omitempty"`
	Templates              []*Template      `json:"templates,omitempty"`
	Redirects              []*Redirect      `json:"redirects,omitempty"`
	IsPartOf               *Project         `json:"is_part_of,omitempty"`
	ArticleBody            *ArticleBody     `json:"article_body,omitempty"`
	License                []*License       `json:"license,omitempty"`
	Visibility             *Visibility      `json:"visibility,omitempty"`
	Event                  *Event           `json:"event,omitempty"`
}

// Category article category representation.
type Category struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Redirect article redirect representation.
type Redirect struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Template article template representation.
type Template struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
