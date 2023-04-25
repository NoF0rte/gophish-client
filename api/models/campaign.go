package models

import (
	"encoding/json"
	"time"
)

type Campaign struct {
	ID            int64           `json:"id" yaml:"-"`
	Name          string          `json:"name" yaml:"name"`
	CreatedDate   time.Time       `json:"created_date,omitempty" yaml:"-"`
	LaunchDate    time.Time       `json:"launch_date" yaml:"launch-date"`
	SendByDate    time.Time       `json:"send_by_date" yaml:"send-by-date"`
	CompletedDate time.Time       `json:"completed_date,omitempty" yaml:"-"`
	Template      *Template       `json:"template" yaml:"template"`
	Page          *Page           `json:"page" yaml:"page"`
	Status        string          `json:"status,omitempty" yaml:"status"`
	Results       []*Result       `json:"results,omitempty" yaml:"results"`
	Groups        []*Group        `json:"groups" yaml:"groups"`
	Timeline      []*Event        `json:"timeline,omitempty" yaml:"timeline"`
	SMTP          *SendingProfile `json:"smtp" yaml:"smtp"`
	URL           string          `json:"url" yaml:"url"`
}

func (c *Campaign) ToJSON() (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type Result struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Position  string  `json:"position"`
	Status    string  `json:"status"`
	IP        string  `json:"ip"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	SendDate  string  `json:"send_date"`
	Reported  bool    `json:"reported"`
}

type Event struct {
	Email   string `json:"email"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Details string `json:"details"`
}
