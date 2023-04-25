package models

import (
	"encoding/json"
	"time"
)

type Group struct {
	ID           int64     `json:"id" yaml:"id"`
	Name         string    `json:"name" yaml:"name"`
	NumTargets   int64     `json:"num_targets,omitempty" yaml:"-"`
	Targets      []Target  `json:"targets,omitempty" yaml:"targets"`
	ModifiedDate time.Time `json:"modified_date" yaml:"-"`
}

func (g *Group) ToJSON() (string, error) {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type Target struct {
	Email     string `json:"email" yaml:"email"`
	FirstName string `json:"first_name" yaml:"first-name"`
	LastName  string `json:"last_name" yaml:"last-name"`
	Position  string `json:"position" yaml:"position"`
}

type GroupsSummary struct {
	Total  int      `json:"total"`
	Groups []*Group `json:"groups"`
}
