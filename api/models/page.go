package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

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

func (p *Page) ToJSON() (string, error) {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (p *Page) replaceVars(vars map[string]string) error {
	if p.varsReplaced {
		return nil
	}

	name, err := templateReplace(p.Name, vars)
	if err != nil {
		return err
	}
	p.Name = name

	html, err := templateReplace(p.HTML, vars)
	if err != nil {
		return err
	}
	p.HTML = html

	redirectURL, err := templateReplace(p.RedirectURL, vars)
	if err != nil {
		return err
	}
	p.RedirectURL = redirectURL

	p.varsReplaced = true
	return nil
}

func PageFromFile(file string, vars map[string]string) (*Page, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var page Page
	err = yaml.Unmarshal(bytes, &page)
	if err != nil {
		return nil, err
	}

	parentDir := filepath.Dir(file)
	if page.HTML == "" && page.HTMLFile != "" {
		htmlPath := filepath.Join(parentDir, page.HTMLFile)

		html, err := os.ReadFile(htmlPath)
		if err != nil {
			return nil, err
		}

		page.HTML = string(html)
	}

	err = page.replaceVars(vars)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

type ImportSite struct {
	URL              string `json:"url"`
	IncludeResources bool   `json:"include_resources"`
}

type ImportedSite struct {
	HTML string `json:"html"`
}
