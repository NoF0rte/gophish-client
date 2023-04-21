package models

import "time"

type Page struct {
	ID                 int64     `json:"id" yaml:"-"`
	Name               string    `json:"name" yaml:"name"`
	HTML               string    `json:"html,omitempty" yaml:"html,omitempty"`
	HTMLFile           string    `json:"-" yaml:"html-file,omitempty"`
	CaptureCredentials bool      `json:"capture_credentials" yaml:"capture-credentials"`
	CapturePasswords   bool      `json:"capture_passwords" yaml:"capture-passwords"`
	ModifiedDate       time.Time `json:"modified_date" yaml:"-"`
	RedirectURL        string    `json:"redirect_url" yaml:"redirect-url"`
	varsReplaced       bool
}

type ImportSite struct {
	URL              string `json:"url"`
	IncludeResources bool   `json:"include_resources"`
}

type ImportedSite struct {
	HTML string `json:"html"`
}
