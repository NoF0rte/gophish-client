package models

import "time"

type Group struct {
	ID           int64     `json:"id" yaml:"id"`
	Name         string    `json:"name" yaml:"name"`
	Targets      []Target  `json:"targets" yaml:"targets"`
	ModifiedDate time.Time `json:"modified_date" yaml:"-"`
}

type Target struct {
	Email     string `json:"email" yaml:"email"`
	FirstName string `json:"first_name" yaml:"first-name"`
	LastName  string `json:"last_name" yaml:"last-name"`
	Position  string `json:"position" yaml:"position"`
}

type GroupSummary struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	NumTargets   int       `json:"num_targets"`
	ModifiedDate time.Time `json:"modified_date"`
}
